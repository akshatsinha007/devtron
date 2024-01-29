package infraConfig

import (
	"github.com/devtron-labs/devtron/pkg/infraConfig/units"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

// repo structs

type InfraProfile struct {
	tableName   struct{} `sql:"infra_profile" pg:",discard_unknown_columns"`
	Id          int      `sql:"id"`
	Name        string   `sql:"name"`
	Description string   `sql:"description"`
	Active      bool     `sql:"active"`
	sql.AuditLog
}

func (infraProfile *InfraProfile) ConvertToProfileBean() ProfileBean {
	profileType := NORMAL
	if infraProfile.Name == DEFAULT_PROFILE_NAME {
		profileType = DEFAULT
	}
	return ProfileBean{
		Id:          infraProfile.Id,
		Name:        infraProfile.Name,
		Type:        profileType,
		Description: infraProfile.Description,
		Active:      infraProfile.Active,
		CreatedBy:   infraProfile.CreatedBy,
		CreatedOn:   infraProfile.CreatedOn,
		UpdatedBy:   infraProfile.UpdatedBy,
		UpdatedOn:   infraProfile.UpdatedOn,
	}
}

type InfraProfileConfiguration struct {
	tableName struct{}         `sql:"infra_profile_configuration" pg:",discard_unknown_columns"`
	Id        int              `sql:"id"`
	Key       ConfigKey        `sql:"key"`
	Value     float64          `sql:"value"`
	Unit      units.UnitSuffix `sql:"unit"`
	ProfileId int              `sql:"profile_id"`
	Active    bool             `sql:"active"`
	sql.AuditLog
}

func (infraProfileConfiguration *InfraProfileConfiguration) ConvertToConfigurationBean() ConfigurationBean {
	return ConfigurationBean{
		Id:        infraProfileConfiguration.Id,
		Key:       GetConfigKeyStr(infraProfileConfiguration.Key),
		Value:     infraProfileConfiguration.Value,
		Unit:      GetUnitSuffixStr(infraProfileConfiguration.Key, infraProfileConfiguration.Unit),
		ProfileId: infraProfileConfiguration.ProfileId,
		Active:    infraProfileConfiguration.Active,
	}
}

// service layer structs

type ProfileBean struct {
	Id             int                 `json:"id"`
	Name           string              `json:"name" validate:"required,min=1,max=50"`
	Description    string              `json:"description" validate:"max=300"`
	Active         bool                `json:"active"`
	Configurations []ConfigurationBean `json:"configurations" validate:"dive"`
	Type           ProfileType         `json:"type"`
	AppCount       int                 `json:"appCount,omitempty"`
	CreatedBy      int32               `json:"createdBy"`
	CreatedOn      time.Time           `json:"createdOn"`
	UpdatedBy      int32               `json:"updatedBy"`
	UpdatedOn      time.Time           `json:"updatedOn"`
}

func (profileBean *ProfileBean) ConvertToInfraProfile() *InfraProfile {
	return &InfraProfile{
		Id:          profileBean.Id,
		Name:        profileBean.Name,
		Description: profileBean.Description,
	}
}

type ConfigurationBean struct {
	Id          int          `json:"id"`
	Key         ConfigKeyStr `json:"key" validate:"required,oneof=cpu_limit cpu_request memory_limit memory_request timeout"`
	Value       float64      `json:"value" validate:"required,min=0"`
	Unit        string       `json:"unit" validate:"required,min=0"`
	ProfileName string       `json:"profileName"`
	ProfileId   int          `json:"profileId"`
	Active      bool         `json:"active"`
}

func (configurationBean *ConfigurationBean) ConvertToInfraProfileConfiguration() *InfraProfileConfiguration {
	return &InfraProfileConfiguration{
		Id:        configurationBean.Id,
		Key:       GetConfigKey(configurationBean.Key),
		Value:     configurationBean.Value,
		Unit:      GetUnitSuffix(configurationBean.Key, configurationBean.Unit),
		ProfileId: configurationBean.ProfileId,
		Active:    configurationBean.Active,
	}
}

type InfraConfigMetaData struct {
	DefaultConfigurations []ConfigurationBean                    `json:"defaultConfigurations"`
	ConfigurationUnits    map[ConfigKeyStr]map[string]units.Unit `json:"configurationUnits"`
}
type ProfileResponse struct {
	Profile ProfileBean `json:"profile"`
	InfraConfigMetaData
}

type ProfilesResponse struct {
	Profiles []ProfileBean `json:"profiles"`
	InfraConfigMetaData
}

type IdentifierListFilter struct {
	IdentifierType     IdentifierType `json:"-"`                              // currently supporting app
	IdentifierNameLike string         `json:"search"`                         // currently app_name  is supported
	ProfileName        string         `json:"profileName"`                    // gets  the list for this profile
	Limit              int            `json:"size" validate:"min=0"`          // limit on the result set , defaults to 20
	Offset             int            `json:"offset" validate:"min=0"`        // offset on the result set, defaults to 0
	SortOrder          string         `json:"sort" validate:"oneof=ASC DESC"` // asc or desc, defaults to asc by appName
}

type Identifier struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Profile *ProfileBean `json:"profile"`

	// for internal use only, do not propagate these values to api response
	ProfileId                 int `json:"-"`
	TotalIdentifierCount      int `json:"-"`
	OverriddenIdentifierCount int `json:"-"`
}

type IdentifierProfileResponse struct {
	Identifiers               []*Identifier `json:"identifiers"`
	TotalIdentifierCount      int           `json:"totalIdentifierCount"`
	OverriddenIdentifierCount int           `json:"overriddenIdentifierCount"`
}

type InfraProfileApplyRequest struct {
	IdentifiersFilter *IdentifierListFilter `json:"identifiersFilter"`
	Identifiers       []string              `json:"identifiers"`
	UpdateToProfile   int                   `json:"updateToProfile"`
	// internal use only
	IdentifierIds []int `json:"-"`
}

type Scope struct {
	AppId int
}

// InfraConfig is used for read only purpose outside this package
type InfraConfig struct {
	// currently only for ci
	CiLimitCpu       string `env:"LIMIT_CI_CPU" envDefault:"0.5"`
	CiLimitMem       string `env:"LIMIT_CI_MEM" envDefault:"3G"`
	CiReqCpu         string `env:"REQ_CI_CPU" envDefault:"0.5"`
	CiReqMem         string `env:"REQ_CI_MEM" envDefault:"3G"`
	CiDefaultTimeout int64  `env:"DEFAULT_TIMEOUT" envDefault:"3600"`
}

func (infraConfig InfraConfig) GetCiLimitCpu() string {
	return infraConfig.CiLimitCpu
}

func (infraConfig InfraConfig) setCiLimitCpu(cpu string) {
	infraConfig.CiLimitCpu = cpu
}

func (infraConfig InfraConfig) GetCiLimitMem() string {
	return infraConfig.CiLimitMem
}

func (infraConfig InfraConfig) setCiLimitMem(mem string) {
	infraConfig.CiLimitMem = mem
}

func (infraConfig InfraConfig) GetCiReqCpu() string {
	return infraConfig.CiReqCpu
}

func (infraConfig InfraConfig) setCiReqCpu(cpu string) {
	infraConfig.CiReqCpu = cpu
}

func (infraConfig InfraConfig) GetCiReqMem() string {
	return infraConfig.CiReqMem
}

func (infraConfig InfraConfig) setCiReqMem(mem string) {
	infraConfig.CiReqMem = mem
}

func (infraConfig InfraConfig) GetCiDefaultTimeout() int64 {
	return infraConfig.CiDefaultTimeout
}

func (infraConfig InfraConfig) setCiDefaultTimeout(timeout int64) {
	infraConfig.CiDefaultTimeout = timeout
}

func (infraConfig InfraConfig) LoadCiLimitCpu() (*InfraProfileConfiguration, error) {
	positive, _, num, denom, suffix, err := units.ParseQuantityString(infraConfig.CiLimitCpu)
	if err != nil {
		return nil, err
	}
	if !positive {
		return nil, errors.New("negative value not allowed for cpu limits")
	}
	valStr := num
	if denom != "" {
		valStr = num + "." + denom
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, err
	}
	return &InfraProfileConfiguration{
		Key:   CPULimit,
		Value: val,
		Unit:  units.GetCPUUnit(units.CPUUnitStr(suffix)),
	}, nil

}

func (infraConfig InfraConfig) LoadCiLimitMem() (*InfraProfileConfiguration, error) {
	positive, _, num, denom, suffix, err := units.ParseQuantityString(infraConfig.CiLimitMem)
	if err != nil {
		return nil, err
	}
	if !positive {
		return nil, errors.New("negative value not allowed for memory limits")
	}
	valStr := num
	if denom != "" {
		valStr = num + "." + denom
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, err
	}
	return &InfraProfileConfiguration{
		Key:   MemoryLimit,
		Value: val,
		Unit:  units.GetMemoryUnit(units.MemoryUnitStr(suffix)),
	}, nil

}

func (infraConfig InfraConfig) LoadCiReqCpu() (*InfraProfileConfiguration, error) {
	positive, _, num, denom, suffix, err := units.ParseQuantityString(infraConfig.CiReqCpu)
	if err != nil {
		return nil, err
	}
	if !positive {
		return nil, errors.New("negative value not allowed for cpu requests")
	}

	valStr := num
	if denom != "" {
		valStr = num + "." + denom
	}

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, err
	}

	return &InfraProfileConfiguration{
		Key:   CPURequest,
		Value: val,
		Unit:  units.GetCPUUnit(units.CPUUnitStr(suffix)),
	}, nil
}

func (infraConfig InfraConfig) LoadCiReqMem() (*InfraProfileConfiguration, error) {
	positive, _, num, denom, suffix, err := units.ParseQuantityString(infraConfig.CiReqMem)
	if err != nil {
		return nil, err
	}
	if !positive {
		return nil, errors.New("negative value not allowed for memory requests")
	}
	valStr := num
	if denom != "" {
		valStr = num + "." + denom
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, err
	}

	return &InfraProfileConfiguration{
		Key:   MemoryRequest,
		Value: val,
		Unit:  units.GetMemoryUnit(units.MemoryUnitStr(suffix)),
	}, nil
}

func (infraConfig InfraConfig) LoadDefaultTimeout() (*InfraProfileConfiguration, error) {
	return &InfraProfileConfiguration{
		Key:   TimeOut,
		Value: float64(infraConfig.CiDefaultTimeout),
		Unit:  units.GetTimeUnit(units.SecondStr),
	}, nil
}
