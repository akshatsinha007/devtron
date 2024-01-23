package service

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/devtron-labs/devtron/pkg/devtronResource"
	"github.com/devtron-labs/devtron/pkg/devtronResource/bean"
	"github.com/devtron-labs/devtron/pkg/infraConfig"
	"github.com/devtron-labs/devtron/pkg/infraConfig/repository"
	"github.com/devtron-labs/devtron/pkg/infraConfig/units"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/devtron-labs/devtron/util"
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

const CannotDeleteDefaultProfile = "cannot delete default profile"
const InvalidProfileName = "profile name is invalid"
const PayloadValidationError = "payload validation failed"

type InfraConfigService interface {

	// todo: @gireesh for all get apis, check if we can get profile and configurations in one db call
	// can use rigth join

	// GetConfigurationUnits fetches all the units for the configurations.
	GetConfigurationUnits() map[infraConfig.ConfigKeyStr]map[string]units.Unit

	// GetDefaultProfile fetches the default profile and its configurations.
	GetDefaultProfile() (*infraConfig.ProfileBean, error)

	// GetProfileByName fetches the profile and its configurations matching the given profileName.
	GetProfileByName(profileName string) (*infraConfig.ProfileBean, error)

	// GetProfileList fetches all the profile and their configurations matching the given profileNameLike string.
	// If profileNameLike is empty, it will fetch all the active profiles.
	GetProfileList(profileNameLike string) (*infraConfig.ProfilesResponse, error)

	CreateProfile(userId int32, profileBean *infraConfig.ProfileBean) error

	// UpdateProfile updates the profile and its configurations matching the given profileName.
	// If profileName is empty, it will return an error.
	UpdateProfile(userId int32, profileName string, profileBean *infraConfig.ProfileBean) error

	// DeleteProfile deletes the profile and its configurations matching the given profileName.
	// If profileName is empty, it will return an error.
	DeleteProfile(profileName string) error
}

type InfraConfigServiceImpl struct {
	logger                              *zap.SugaredLogger
	infraProfileRepo                    repository.InfraConfigRepository
	units                               *units.Units
	infraConfig                         *infraConfig.InfraConfig
	devtronResourceSearchableKeyService devtronResource.DevtronResourceSearchableKeyService
	validator                           *validator.Validate
}

func NewInfraConfigServiceImpl(logger *zap.SugaredLogger,
	infraProfileRepo repository.InfraConfigRepository,
	units *units.Units,
	devtronResourceSearchableKeyService devtronResource.DevtronResourceSearchableKeyService,
	validator *validator.Validate) (*InfraConfigServiceImpl, error) {
	infraConfiguration := &infraConfig.InfraConfig{}
	err := env.Parse(infraConfiguration)
	if err != nil {
		return nil, err
	}
	infraProfileService := &InfraConfigServiceImpl{
		logger:                              logger,
		infraProfileRepo:                    infraProfileRepo,
		units:                               units,
		devtronResourceSearchableKeyService: devtronResourceSearchableKeyService,
		infraConfig:                         infraConfiguration,
		validator:                           validator,
	}
	err = infraProfileService.loadDefaultProfile()
	return infraProfileService, err
}

func (impl *InfraConfigServiceImpl) GetDefaultProfile() (*infraConfig.ProfileBean, error) {
	infraProfile, err := impl.infraProfileRepo.GetProfileByName(repository.DEFAULT_PROFILE_NAME)
	if err != nil {
		impl.logger.Errorw("error in fetching default profile", "error", err)
		return nil, err
	}

	profileBean := infraProfile.ConvertToProfileBean()
	infraConfigurations, err := impl.infraProfileRepo.GetConfigurationsByProfileId([]int{infraProfile.Id})
	if err != nil {
		impl.logger.Errorw("error in fetching default configurations", "error", err)
		return nil, err
	}

	configurationBeans := util.Transform(infraConfigurations, func(config *infraConfig.InfraProfileConfiguration) infraConfig.ConfigurationBean {
		configBean := config.ConvertToConfigurationBean()
		configBean.ProfileName = profileBean.Name
		return configBean
	})
	profileBean.Configurations = configurationBeans
	return &profileBean, nil
}

func (impl *InfraConfigServiceImpl) GetProfileByName(profileName string) (*infraConfig.ProfileBean, error) {
	if profileName == "" {
		return nil, errors.New(InvalidProfileName)
	}

	infraProfile, err := impl.infraProfileRepo.GetProfileByName(profileName)
	if err != nil {
		impl.logger.Errorw("error in fetching default profile", "error", err)
		return nil, err
	}

	profileBean := infraProfile.ConvertToProfileBean()
	infraConfigurations, err := impl.infraProfileRepo.GetConfigurationsByProfileId([]int{infraProfile.Id})
	if err != nil {
		impl.logger.Errorw("error in fetching default configurations", "error", err)
		return nil, err
	}

	configurationBeans := util.Transform(infraConfigurations, func(config *infraConfig.InfraProfileConfiguration) infraConfig.ConfigurationBean {
		configBean := config.ConvertToConfigurationBean()
		configBean.ProfileName = profileBean.Name
		return configBean
	})
	profileBean.Configurations = configurationBeans
	return &profileBean, nil
}

func (impl *InfraConfigServiceImpl) GetProfileList(profileNameLike string) (*infraConfig.ProfilesResponse, error) {
	// fetch all the profiles matching the given profileNameLike filter
	infraProfiles, err := impl.infraProfileRepo.GetProfileList(profileNameLike)
	defaultProfileIndx := 0
	// extract out profileIds from the profiles
	profileIds := make([]int, len(infraProfiles))
	profilesMap := make(map[int]infraConfig.ProfileBean)
	for i, _ := range infraProfiles {
		profileBean := infraProfiles[i].ConvertToProfileBean()
		if profileBean.Name == repository.DEFAULT_PROFILE_NAME {
			defaultProfileIndx = i
		}
		profileIds[i] = profileBean.Id
		profilesMap[profileBean.Id] = profileBean
	}

	// fetch all the configurations matching the given profileIds
	infraConfigurations, err := impl.infraProfileRepo.GetConfigurationsByProfileId(profileIds)
	if err != nil {
		impl.logger.Errorw("error in fetching default configurations", "profileIds", profileIds, "error", err)
		return nil, err
	}

	// map the configurations to their respective profiles
	for _, configuration := range infraConfigurations {
		profileBean := profilesMap[configuration.ProfileId]
		configurationBean := configuration.ConvertToConfigurationBean()
		configurationBean.ProfileName = profileBean.Name
		profileBean.Configurations = append(profileBean.Configurations, configurationBean)
		profilesMap[configuration.ProfileId] = profileBean
	}

	searchableKeyNameIdMap := impl.devtronResourceSearchableKeyService.GetAllSearchableKeyNameIdMap()
	profileAppCount, err := impl.infraProfileRepo.GetIdentifierCountForNonDefaultProfiles(profileIds, searchableKeyNameIdMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_APP_ID])
	if err != nil {
		impl.logger.Errorw("error in fetching app count for non default profiles", "error", err)
		return nil, err
	}

	// set app count for each profile
	for _, profileAppCnt := range profileAppCount {
		profileBean := profilesMap[profileAppCnt.ProfileId]
		profileBean.AppCount = profileAppCnt.IdentifierCount
		profilesMap[profileAppCnt.ProfileId] = profileBean
	}

	// fill the default configurations for each profile if any of the default configuration is missing
	defaultProfile := profilesMap[profileIds[defaultProfileIndx]]
	defaultConfigurations := defaultProfile.Configurations

	profiles := make([]infraConfig.ProfileBean, 0, len(profilesMap))
	for profileId, profile := range profilesMap {
		if profile.Name == repository.DEFAULT_PROFILE_NAME {
			continue
		}

		for _, defaultConfiguration := range defaultConfigurations {
			// if profile doesn't have the default configuration, add it to the profile
			if !util.Contains(profile.Configurations, func(config infraConfig.ConfigurationBean) bool {
				return config.Key == defaultConfiguration.Key
			}) {
				profile.Configurations = append(profile.Configurations, defaultConfiguration)
			}
		}

		profiles = append(profiles, profile)
		// update map with updated profile
		profilesMap[profileId] = profile
	}
	resp := &infraConfig.ProfilesResponse{
		Profiles: profiles,
	}
	resp.DefaultConfigurations = defaultConfigurations
	resp.ConfigurationUnits = impl.GetConfigurationUnits()
	return resp, nil
}

func (impl *InfraConfigServiceImpl) CreateProfile(userId int32, profileBean *infraConfig.ProfileBean) error {
	defaultProfile, err := impl.GetDefaultProfile()
	if err != nil {
		impl.logger.Errorw("error in fetching default profile", "profileCreateRequest", profileBean, "error", err)
		return err
	}
	if err := impl.Validate(profileBean, defaultProfile); err != nil {
		impl.logger.Errorw("error occurred in validation the profile create request", "profileCreateRequest", profileBean, "error", err)
		return err
	}

	infraProfile := profileBean.ConvertToInfraProfile()
	infraProfile.Active = true
	infraProfile.CreatedBy = userId
	infraProfile.CreatedOn = time.Now()

	tx, err := impl.infraProfileRepo.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting transaction to create profile", "profileCreateRequest", profileBean, "error", err)
		return err
	}
	defer impl.infraProfileRepo.RollbackTx(tx)
	err = impl.infraProfileRepo.CreateProfile(tx, infraProfile)
	if err != nil {
		impl.logger.Errorw("error in creating infra config profile", "infraProfile", infraProfile, "err", err)
		return err
	}

	infraConfigurations := make([]*infraConfig.InfraProfileConfiguration, 0, len(profileBean.Configurations))
	for _, configuration := range profileBean.Configurations {
		infraConfiguration := configuration.ConvertToInfraProfileConfiguration()
		infraConfiguration.Active = true
		infraConfiguration.CreatedBy = userId
		infraConfiguration.CreatedOn = time.Now()
		infraConfiguration.ProfileId = infraProfile.Id
		infraConfigurations = append(infraConfigurations, infraConfiguration)
	}

	err = impl.infraProfileRepo.CreateConfigurations(tx, infraConfigurations)
	if err != nil {
		impl.logger.Errorw("error in creating infra configurations", "infraConfigurations", infraConfigurations, "err", err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		impl.logger.Errorw("error in committing transaction while creating infra configuration profile", "profileCreateRequest", profileBean, "err", err)
		return err
	}
	return nil
}

func (impl *InfraConfigServiceImpl) UpdateProfile(userId int32, profileName string, profileBean *infraConfig.ProfileBean) error {
	if profileName == "" {
		return errors.New(InvalidProfileName)
	}

	// validation
	defaultProfile, err := impl.GetDefaultProfile()
	if err != nil {
		impl.logger.Errorw("error in fetching default profile", "profileCreateRequest", profileBean, "error", err)
		return err
	}
	if err := impl.Validate(profileBean, defaultProfile); err != nil {
		impl.logger.Errorw("error occurred in validation the profile create request", "profileCreateRequest", profileBean, "error", err)
		return err
	}
	// validations end

	infraProfile := profileBean.ConvertToInfraProfile()
	// user couldn't delete the profile, always set this to active
	infraProfile.Active = true
	infraConfigurations := util.Transform(profileBean.Configurations, func(config infraConfig.ConfigurationBean) *infraConfig.InfraProfileConfiguration {
		config.ProfileId = infraProfile.Id
		// user couldn't delete the configuration for default profile, always set this to active
		if infraProfile.Name == repository.DEFAULT_PROFILE_NAME {
			config.Active = true
		}
		return config.ConvertToInfraProfileConfiguration()
	})

	tx, err := impl.infraProfileRepo.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting transaction to update profile", "error", err)
		return err
	}
	defer impl.infraProfileRepo.RollbackTx(tx)
	infraProfile.UpdatedOn = time.Now()
	infraProfile.UpdatedBy = userId
	err = impl.infraProfileRepo.UpdateProfile(tx, profileName, infraProfile)
	if err != nil {
		impl.logger.Errorw("error in updating profile", "error", err)
		return err
	}

	err = impl.infraProfileRepo.UpdateConfigurations(tx, infraConfigurations)
	if err != nil {
		impl.logger.Errorw("error in creating configurations", "error", err)
		return err
	}
	err = impl.infraProfileRepo.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing transaction to update profile", "profileCreateRequest", profileBean, "error", err)
	}
	return err
}

func (impl *InfraConfigServiceImpl) DeleteProfile(profileName string) error {
	if profileName == "" {
		return errors.New(InvalidProfileName)
	}

	if profileName == repository.DEFAULT_PROFILE_NAME {
		return errors.New(CannotDeleteDefaultProfile)
	}

	tx, err := impl.infraProfileRepo.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting transaction to delete profile", "profileName", profileName, "error", err)
		return err
	}
	defer impl.infraProfileRepo.RollbackTx(tx)

	// step1: delete profile
	err = impl.infraProfileRepo.DeleteProfile(tx, profileName)
	if err != nil {
		impl.logger.Errorw("error in deleting profile", "profileName", profileName, "error", err)
		return err
	}

	// step2: delete configurations
	err = impl.infraProfileRepo.DeleteConfigurations(tx, profileName)
	if err != nil {
		impl.logger.Errorw("error in deleting configurations", "profileName", profileName, "error", err)
	}

	// step3: delete profile identifier mappings
	// todo: delete from resource_identifier_mapping where resource_id is profileId and resource_type is infraProfile
	err = impl.infraProfileRepo.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing transaction to delete profile", "profileName", profileName, "error", err)
	}
	return err
}

func (impl *InfraConfigServiceImpl) loadDefaultProfile() error {

	profile, err := impl.infraProfileRepo.GetProfileByName(repository.DEFAULT_PROFILE_NAME)
	if err != nil && !errors.Is(err, pg.ErrNoRows) {
		return err
	}
	if profile.Id != 0 {
		// return here because entry already exists and we dont need to create it again
		return nil
	}

	infraConfiguration := impl.infraConfig
	cpuLimit, err := infraConfiguration.GetCiLimitCpu()
	if err != nil {
		return err
	}
	memLimit, err := infraConfiguration.GetCiLimitMem()
	if err != nil {
		return err
	}
	cpuReq, err := infraConfiguration.GetCiReqCpu()
	if err != nil {
		return err
	}
	memReq, err := infraConfiguration.GetCiReqMem()
	if err != nil {
		return err
	}
	timeout, err := infraConfiguration.GetDefaultTimeout()
	if err != nil {
		return err
	}

	defaultConfigurations := []*infraConfig.InfraProfileConfiguration{cpuLimit, memLimit, cpuReq, memReq, timeout}
	defaultProfile := &infraConfig.InfraProfile{
		Name:        repository.DEFAULT_PROFILE_NAME,
		Description: "",
		Active:      true,
		AuditLog: sql.AuditLog{
			CreatedBy: 1, // system user
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			UpdatedBy: 1, // system user
		},
	}
	tx, err := impl.infraProfileRepo.StartTx()
	if err != nil {
		impl.logger.Errorw("error in starting transaction to save default configurations", "error", err)
		return err
	}

	defer impl.infraProfileRepo.RollbackTx(tx)
	err = impl.infraProfileRepo.CreateProfile(tx, defaultProfile)
	if err != nil {
		impl.logger.Errorw("error in saving default profile", "error", err)
		return err
	}

	util.Transform(defaultConfigurations, func(config *infraConfig.InfraProfileConfiguration) *infraConfig.InfraProfileConfiguration {
		config.ProfileId = defaultProfile.Id
		config.Active = true
		config.AuditLog = sql.AuditLog{
			CreatedBy: 1, // system user
			CreatedOn: time.Now(),
			UpdatedOn: time.Now(),
			UpdatedBy: 1, // system user
		}
		return config
	})
	err = impl.infraProfileRepo.CreateConfigurations(tx, defaultConfigurations)
	if err != nil {
		impl.logger.Errorw("error in saving default configurations", "error", err)
		return err
	}
	err = impl.infraProfileRepo.CommitTx(tx)
	if err != nil {
		impl.logger.Errorw("error in committing transaction to save default configurations", "error", err)
	}
	return err
}

func (impl *InfraConfigServiceImpl) GetConfigurationUnits() map[infraConfig.ConfigKeyStr]map[string]units.Unit {
	configurationUnits := make(map[infraConfig.ConfigKeyStr]map[string]units.Unit)
	configurationUnits[infraConfig.CPU_REQUEST] = impl.units.GetCpuUnits()
	configurationUnits[infraConfig.CPU_LIMIT] = impl.units.GetCpuUnits()

	configurationUnits[infraConfig.MEMORY_REQUEST] = impl.units.GetMemoryUnits()
	configurationUnits[infraConfig.MEMORY_LIMIT] = impl.units.GetMemoryUnits()

	configurationUnits[infraConfig.TIME_OUT] = impl.units.GetTimeUnits()

	return configurationUnits
}

func (impl *InfraConfigServiceImpl) Validate(profileBean *infraConfig.ProfileBean, defaultProfile *infraConfig.ProfileBean) error {
	err := impl.validator.Struct(profileBean)
	if err != nil {
		err = errors.Wrap(err, PayloadValidationError)
		return err
	}

	defaultConfigurationsKeyMap := infraConfig.GetDefaultConfigKeysMap()
	// validate configurations only contain default configurations types.(cpu_limit,cpu_request,mem_limit,mem_request,timeout)
	for _, propertyConfig := range profileBean.Configurations {
		if _, ok := defaultConfigurationsKeyMap[propertyConfig.Key]; !ok {
			if err == nil {
				err = errors.New(fmt.Sprintf("invalid configuration property \"%s\"", propertyConfig.Key))
			}
			err = errors.Wrap(err, fmt.Sprintf("invalid configuration property \"%s\"", propertyConfig.Key))
		}

	}

	if err != nil {
		err = errors.Wrap(err, PayloadValidationError)
		return err
	}

	err = impl.validateCpuMem(profileBean, defaultProfile)
	if err != nil {
		err = errors.Wrap(err, PayloadValidationError)
		return err
	}
	return nil
}

func (impl *InfraConfigServiceImpl) validateCpuMem(profileBean *infraConfig.ProfileBean, defaultProfile *infraConfig.ProfileBean) error {

	configurationUnits := impl.units
	// currently validating cpu and memory limits and reqs only
	var (
		cpuLimit *infraConfig.ConfigurationBean
		cpuReq   *infraConfig.ConfigurationBean
		memLimit *infraConfig.ConfigurationBean
		memReq   *infraConfig.ConfigurationBean
	)

	for _, propertyConfig := range profileBean.Configurations {
		// get cpu limit and req
		switch propertyConfig.Key {
		case infraConfig.CPU_LIMIT:
			cpuLimit = &propertyConfig
		case infraConfig.CPU_REQUEST:
			cpuReq = &propertyConfig
		case infraConfig.MEMORY_LIMIT:
			memLimit = &propertyConfig
		case infraConfig.MEMORY_REQUEST:
			memReq = &propertyConfig
		}
	}

	for _, defaultPropertyConfig := range defaultProfile.Configurations {
		// get cpu limit and req
		switch defaultPropertyConfig.Key {
		case infraConfig.CPU_LIMIT:
			if cpuLimit == nil {
				cpuLimit = &defaultPropertyConfig
			}
		case infraConfig.CPU_REQUEST:
			if cpuReq == nil {
				cpuReq = &defaultPropertyConfig
			}
		case infraConfig.MEMORY_LIMIT:
			if memLimit == nil {
				memLimit = &defaultPropertyConfig
			}
		case infraConfig.MEMORY_REQUEST:
			if memReq == nil {
				memReq = &defaultPropertyConfig
			}
		}
	}

	// validate cpu
	cpuLimitUnitSuffix := units.GetCPUUnit(units.CPUUnitStr(cpuLimit.Unit))
	cpuReqUnitSuffix := units.GetCPUUnit(units.CPUUnitStr(cpuReq.Unit))
	var cpuLimitUnit *units.Unit
	var cpuReqUnit *units.Unit
	for cpuUnitSuffix, cpuUnit := range configurationUnits.GetCpuUnits() {
		if string(units.GetCPUUnitStr(cpuLimitUnitSuffix)) == cpuUnitSuffix {
			cpuLimitUnit = &cpuUnit
		}

		if string(units.GetCPUUnitStr(cpuReqUnitSuffix)) == cpuUnitSuffix {
			cpuReqUnit = &cpuUnit
		}

	}

	// this condition should be true for valid case => (lim/req)*(lf/rf) >= 1
	limitToReqRationCPU := cpuLimit.Value / cpuReq.Value
	convFactorCPU := cpuLimitUnit.ConversionFactor / cpuReqUnit.ConversionFactor

	if limitToReqRationCPU*convFactorCPU < 1 {
		return errors.New("cpu limit should not be less than cpu request")
	}

	// validate mem

	memLimitUnitSuffix := units.GetMemoryUnit(units.MemoryUnitStr(memLimit.Unit))
	memReqUnitSuffix := units.GetMemoryUnit(units.MemoryUnitStr(memReq.Unit))
	var memLimitUnit *units.Unit
	var memReqUnit *units.Unit

	for memUnitSuffix, memUnit := range configurationUnits.GetMemoryUnits() {
		if string(units.GetMemoryUnitStr(memLimitUnitSuffix)) == memUnitSuffix {
			memLimitUnit = &memUnit
		}

		if string(units.GetMemoryUnitStr(memReqUnitSuffix)) == memUnitSuffix {
			memReqUnit = &memUnit
		}
	}

	// this condition should be true for valid case => (lim/req)*(lf/rf) >= 1
	limitToReqRationMem := memLimit.Value / memReq.Value
	convFactorMem := memLimitUnit.ConversionFactor / memReqUnit.ConversionFactor

	if limitToReqRationMem*convFactorMem < 1 {
		return errors.New("memory limit should not be less than memory request")
	}

	return nil
}
