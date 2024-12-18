package adapter

import (
	"errors"
	"fmt"
	"github.com/devtron-labs/devtron/pkg/infraConfig/bean"
	"github.com/devtron-labs/devtron/pkg/infraConfig/repository"
	unitsBean "github.com/devtron-labs/devtron/pkg/infraConfig/units/bean"
	"github.com/devtron-labs/devtron/pkg/infraConfig/util"
	"github.com/devtron-labs/devtron/pkg/sql"
	globalUtil "github.com/devtron-labs/devtron/util"
	"math"
	"reflect"
	"strconv"
)

func ConvertToPlatformMap(infraProfileConfigurationEntities []*repository.InfraProfileConfigurationEntity, profileName string) (map[string][]*bean.ConfigurationBean, error) {
	// Validate input parameters
	if infraProfileConfigurationEntities == nil {
		return nil, fmt.Errorf("input infraProfileConfigurationEntities is empty")
	}
	if profileName == "" {
		return nil, fmt.Errorf("profileName cannot be empty")
	}
	platformMap := make(map[string][]*bean.ConfigurationBean)
	for _, infraProfileConfiguration := range infraProfileConfigurationEntities {
		configurationBean, err := getConfigurationBean(infraProfileConfiguration, profileName)
		if err != nil {
			return nil, fmt.Errorf("failed to get configuration bean for profile from infraConfiguration '%s': %w", profileName, err)
		}
		platform := infraProfileConfiguration.Platform
		if len(platform) == 0 {
			platform = bean.DEFAULT_PLATFORM
		}

		// Add the ConfigurationBean to the corresponding platform entry in the map
		platformMap[platform] = append(platformMap[platform], configurationBean)
	}
	return platformMap, nil
}

// ConvertFromPlatformMap converts map[platform][]*ConfigurationBean back to []InfraProfileConfigurationEntity
func ConvertFromPlatformMap(platformMap map[string][]*bean.ConfigurationBean, profileBean *bean.ProfileBeanDto, userId int32) []*repository.InfraProfileConfigurationEntity {
	var entities []*repository.InfraProfileConfigurationEntity
	for platform, beans := range platformMap {
		for _, configBean := range beans {
			entity := getInfraProfileEntity(configBean, profileBean, platform, userId)
			entities = append(entities, entity)
		}
	}
	return entities
}

// Function to convert valueString to interface{} based on key
func convertValueStringToInterface(configKey bean.ConfigKeyStr, valueString string) (interface{}, error) {
	switch configKey {
	case bean.CPU_LIMIT, bean.CPU_REQUEST, bean.MEMORY_LIMIT, bean.MEMORY_REQUEST:
		// Convert string to float64 and truncate to 2 decimal places
		valueFloat, err := strconv.ParseFloat(valueString, 64)
		truncateValue := globalUtil.TruncateFloat(valueFloat, 2)
		return truncateValue, err // Returning float64 for resource values
	case bean.TIME_OUT:
		// Convert string to float64 and ensure it's within integer range
		valueFloat, err := strconv.ParseFloat(valueString, 64)
		modifiedValue := math.Min(math.Floor(valueFloat), math.MaxInt64)
		return modifiedValue, err // Returning float64 for timeout

	// Add more cases as needed for different config keys
	default:
		// Default case, return the string as is
		err := errors.New(fmt.Sprintf("unsupported key found %s", configKey))
		return nil, err
	}
}

func getConfigurationBean(infraProfileConfiguration *repository.InfraProfileConfigurationEntity, profileName string) (*bean.ConfigurationBean, error) {
	valueString := infraProfileConfiguration.ValueString
	// handle old values
	if len(valueString) == 0 && infraProfileConfiguration.Unit > 0 {
		valueString = strconv.FormatFloat(infraProfileConfiguration.Value, 'f', -1, 64)
	}
	valueInterface, err := convertValueStringToInterface(util.GetConfigKeyStr(infraProfileConfiguration.Key), valueString)
	if err != nil {
		return &bean.ConfigurationBean{}, err
	}
	return &bean.ConfigurationBean{
		ConfigurationBeanAbstract: bean.ConfigurationBeanAbstract{
			Id:  infraProfileConfiguration.Id,
			Key: util.GetConfigKeyStr(infraProfileConfiguration.Key),

			Unit:        util.GetUnitSuffixStr(infraProfileConfiguration.Key, infraProfileConfiguration.Unit),
			ProfileId:   infraProfileConfiguration.ProfileId,
			Active:      infraProfileConfiguration.Active,
			ProfileName: profileName,
		},
		Value: valueInterface,
	}, nil
}

func getInfraProfileEntity(configurationBean *bean.ConfigurationBean, profileBean *bean.ProfileBeanDto, platform string, userId int32) *repository.InfraProfileConfigurationEntity {

	infraProfile := &repository.InfraProfileConfigurationEntity{
		Id:          configurationBean.Id,
		Key:         util.GetConfigKey(configurationBean.Key),
		ValueString: FormatTypedValueAsString(configurationBean.Key, configurationBean.Value),
		Unit:        util.GetUnitSuffix(configurationBean.Key, configurationBean.Unit),
		ProfileId:   profileBean.Id,
		Platform:    platform,
		Active:      configurationBean.Active,
		AuditLog:    sql.NewDefaultAuditLog(userId),
	}
	if profileBean.Name == bean.GLOBAL_PROFILE_NAME {
		infraProfile.Active = true
	}
	return infraProfile
}

func FormatTypedValueAsString(configKey bean.ConfigKeyStr, configValue interface{}) string {
	if configKey == bean.CPU_LIMIT ||
		configKey == bean.CPU_REQUEST ||
		configKey == bean.MEMORY_LIMIT ||
		configKey == bean.MEMORY_REQUEST {
		var valueFloat float64
		// Handle string input or directly as float64
		switch v := configValue.(type) {
		case string:
			valueFloat, _ = strconv.ParseFloat(v, 64)
		case float64:
			valueFloat = v
		}
		// Truncate and format the float value
		truncateValue := globalUtil.TruncateFloat(valueFloat, 2)
		return strconv.FormatFloat(truncateValue, 'f', -1, 64)
		//valueFloat, _ := strconv.ParseFloat(configValue.(float64), 64)
	}

	if configKey == bean.TIME_OUT {
		var valueFloat float64
		switch v := configValue.(type) {
		case string:
			valueFloat, _ = strconv.ParseFloat(v, 64)
		case float64:
			valueFloat = v
		}
		//valueFloat, _ := strconv.ParseFloat(configValue, 64)
		modifiedValue := math.Min(math.Floor(valueFloat), math.MaxInt64)
		return strconv.FormatFloat(modifiedValue, 'f', -1, 64)
	}

	return configValue.(string)
}

func GetV0ProfileBean(profileBean *bean.ProfileBeanDto) *bean.ProfileBeanV0 {
	if profileBean == nil {
		return &bean.ProfileBeanV0{}
	}
	profileName := profileBean.Name
	if profileName == bean.GLOBAL_PROFILE_NAME {
		profileName = bean.DEFAULT_PROFILE_NAME
	}

	profileType := profileBean.Type
	if profileType == bean.GLOBAL {
		profileType = bean.DEFAULT
	}

	ciRunnerConfig := profileBean.Configurations[bean.DEFAULT_PLATFORM]
	return &bean.ProfileBeanV0{
		ProfileBeanAbstract: bean.ProfileBeanAbstract{
			Id:          profileBean.Id,
			Name:        profileName,
			Description: profileBean.Description,
			Active:      profileBean.Active,
			Type:        profileType,
			AppCount:    profileBean.AppCount,
			CreatedBy:   profileBean.CreatedBy,
			CreatedOn:   profileBean.CreatedOn,
			UpdatedBy:   profileBean.UpdatedBy,
			UpdatedOn:   profileBean.UpdatedOn,
		},
		Configurations: GetV0ConfigurationBeans(ciRunnerConfig, profileName),
	}

}

func GetV1ProfileBean(profileBean *bean.ProfileBeanV0) *bean.ProfileBeanDto {
	if profileBean == nil {
		return nil
	}
	profileName := profileBean.Name
	if profileName == bean.DEFAULT_PROFILE_NAME {
		profileName = bean.GLOBAL_PROFILE_NAME
	}
	profileType := profileBean.Type
	if profileType == bean.GLOBAL {
		profileType = bean.DEFAULT
	}
	return &bean.ProfileBeanDto{
		ProfileBeanAbstract: bean.ProfileBeanAbstract{
			Id:          profileBean.Id,
			Name:        profileName,
			Description: profileBean.Description,
			Active:      profileBean.Active,
			Type:        profileType,
			AppCount:    profileBean.AppCount,
			CreatedBy:   profileBean.CreatedBy,
			CreatedOn:   profileBean.CreatedOn,
			UpdatedBy:   profileBean.UpdatedBy,
			UpdatedOn:   profileBean.UpdatedOn,
		},
		Configurations: map[string][]*bean.ConfigurationBean{bean.DEFAULT_PLATFORM: GetV1ConfigurationBeans(profileBean.Configurations, profileName)},
	}

}

func GetV1ConfigurationBeans(configBeans []bean.ConfigurationBeanV0, profileName string) []*bean.ConfigurationBean {
	if len(configBeans) == 0 {
		return nil
	}
	resp := make([]*bean.ConfigurationBean, 0)
	for _, configBean := range configBeans {
		valueString := strconv.FormatFloat(configBean.Value, 'f', -1, 64)

		configBeanV1 := &bean.ConfigurationBean{
			ConfigurationBeanAbstract: bean.ConfigurationBeanAbstract{
				Id:          configBean.Id,
				Key:         configBean.Key,
				Unit:        configBean.Unit,
				ProfileName: profileName,
				ProfileId:   configBean.ProfileId,
				Active:      configBean.Active,
			},
			Value: valueString,
		}
		resp = append(resp, configBeanV1)
	}
	return resp
}

func GetV0ConfigurationBeans(configBeans []*bean.ConfigurationBean, profileName string) []bean.ConfigurationBeanV0 {
	if len(configBeans) == 0 {
		return []bean.ConfigurationBeanV0{}
	}

	resp := make([]bean.ConfigurationBeanV0, 0)
	for _, configBean := range configBeans {
		valueFloat, _ := configBean.Value.(float64)
		//valueFloat, _ := strconv.ParseFloat(configBean.Value, 64)

		beanv0 := bean.ConfigurationBeanV0{
			ConfigurationBeanAbstract: bean.ConfigurationBeanAbstract{
				Id:          configBean.Id,
				Key:         configBean.Key,
				Unit:        configBean.Unit,
				ProfileName: profileName,
				ProfileId:   configBean.ProfileId,
				Active:      configBean.Active,
			},
			Value: valueFloat,
		}
		resp = append(resp, beanv0)
	}
	return resp
}

func ConvertToProfileBean(infraProfile *repository.InfraProfileEntity) bean.ProfileBeanDto {
	profileType := bean.GLOBAL
	if infraProfile.Name != bean.GLOBAL_PROFILE_NAME {
		profileType = bean.NORMAL
	}
	return bean.ProfileBeanDto{
		ProfileBeanAbstract: bean.ProfileBeanAbstract{
			Id:          infraProfile.Id,
			Name:        infraProfile.Name,
			Type:        profileType,
			Description: infraProfile.Description,
			Active:      infraProfile.Active,
			CreatedBy:   infraProfile.CreatedBy,
			CreatedOn:   infraProfile.CreatedOn,
			UpdatedBy:   infraProfile.UpdatedBy,
			UpdatedOn:   infraProfile.UpdatedOn,
		},
	}
}

func ConvertToInfraProfileEntity(profileBean *bean.ProfileBeanDto) *repository.InfraProfileEntity {
	return &repository.InfraProfileEntity{
		Id:          profileBean.Id,
		Name:        profileBean.Name,
		Description: profileBean.Description,
	}
}

func LoadCiLimitCpu(parsedValue *unitsBean.ParsedValue) (*repository.InfraProfileConfigurationEntity, error) {
	return &repository.InfraProfileConfigurationEntity{
		Key:         bean.CPULimitKey,
		ValueString: parsedValue.GetValueString(),
		Unit:        unitsBean.CPUUnitStr(parsedValue.GetUnitType()).GetUnitSuffix(),
		Platform:    bean.DEFAULT_PLATFORM,
	}, nil

}

func LoadCiReqCpu(parsedValue *unitsBean.ParsedValue) (*repository.InfraProfileConfigurationEntity, error) {
	return &repository.InfraProfileConfigurationEntity{
		Key:         bean.CPURequestKey,
		ValueString: parsedValue.GetValueString(),
		Unit:        unitsBean.CPUUnitStr(parsedValue.GetUnitType()).GetUnitSuffix(),
		Platform:    bean.DEFAULT_PLATFORM,
	}, nil
}

func LoadCiReqMem(parsedValue *unitsBean.ParsedValue) (*repository.InfraProfileConfigurationEntity, error) {
	return &repository.InfraProfileConfigurationEntity{
		Key:         bean.MemoryRequestKey,
		ValueString: parsedValue.GetValueString(),
		Unit:        unitsBean.MemoryUnitStr(parsedValue.GetUnitType()).GetUnitSuffix(),
		Platform:    bean.DEFAULT_PLATFORM,
	}, nil
}

func LoadCiLimitMem(parsedValue *unitsBean.ParsedValue) (*repository.InfraProfileConfigurationEntity, error) {
	return &repository.InfraProfileConfigurationEntity{
		Key:         bean.MemoryLimitKey,
		ValueString: parsedValue.GetValueString(),
		Unit:        unitsBean.MemoryUnitStr(parsedValue.GetUnitType()).GetUnitSuffix(),
		Platform:    bean.DEFAULT_PLATFORM,
	}, nil

}

func LoadDefaultTimeout(parsedValue *unitsBean.ParsedValue) (*repository.InfraProfileConfigurationEntity, error) {
	return &repository.InfraProfileConfigurationEntity{
		Key:         bean.TimeOutKey,
		ValueString: parsedValue.GetValueString(),
		Unit:        unitsBean.TimeUnitStr(parsedValue.GetUnitType()).GetUnitSuffix(),
		Platform:    bean.DEFAULT_PLATFORM,
	}, nil
}

func GetTypedValue(configKey bean.ConfigKeyStr, value interface{}) (interface{}, error) {
	switch configKey {
	case bean.CPU_LIMIT, bean.CPU_REQUEST, bean.MEMORY_LIMIT, bean.MEMORY_REQUEST:
		//value is float64 or convertible to it
		switch v := value.(type) {
		case string:
			valueFloat, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse string to float for %s: %w", configKey, err)
			}
			return globalUtil.TruncateFloat(valueFloat, 2), nil
		case float64:
			return globalUtil.TruncateFloat(v, 2), nil
		default:
			return nil, fmt.Errorf("unsupported type for %s: %v", configKey, reflect.TypeOf(value))
		}
	case bean.TIME_OUT:
		switch v := value.(type) {
		case string:
			valueFloat, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse string to float for %s: %w", configKey, err)
			}
			return math.Min(math.Floor(valueFloat), math.MaxInt64), nil
		case float64:
			return math.Min(math.Floor(v), math.MaxInt64), nil
		default:
			return nil, fmt.Errorf("unsupported type for %s: %v", configKey, reflect.TypeOf(value))
		}
	// Default case
	default:
		return nil, fmt.Errorf("unsupported config key: %s", configKey)
	}
}
