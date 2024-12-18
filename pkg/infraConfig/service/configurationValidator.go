package service

import (
	"github.com/devtron-labs/devtron/pkg/infraConfig/bean"
)

func (impl *InfraConfigServiceImpl) validateInfraConfig(profileBean *bean.ProfileBeanDto, defaultProfile *bean.ProfileBeanDto) error {
	for propertyType, factoryService := range impl.unitFactoryMap {
		validationErr := factoryService.Validate(profileBean, defaultProfile)
		if validationErr != nil {
			impl.logger.Error("Error while validating configuration", "propertyType", propertyType, "error", validationErr)
			return validationErr
		}
	}
	return nil
}
