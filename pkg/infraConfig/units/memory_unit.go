/*
 * Copyright (c) 2024. Devtron Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package units

import (
	"errors"
	"fmt"
	"github.com/devtron-labs/devtron/pkg/infraConfig/adapter"
	"github.com/devtron-labs/devtron/pkg/infraConfig/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/infraConfig/units/bean"
	"go.uber.org/zap"
)

type MemoryUnitFactory struct {
	logger      *zap.SugaredLogger
	memoryUnits map[bean2.MemoryUnitStr]bean.Unit
}

func NewMemoryUnitFactory(logger *zap.SugaredLogger) *MemoryUnitFactory {
	return &MemoryUnitFactory{
		logger:      logger,
		memoryUnits: bean2.GetMemoryUnit(),
	}
}

func (m *MemoryUnitFactory) GetAllUnits() map[string]bean.Unit {
	memoryUnits := m.memoryUnits
	units := make(map[string]bean.Unit)
	for key, value := range memoryUnits {
		units[string(key)] = value
	}
	return units
}

func (m *MemoryUnitFactory) ParseValAndUnit(val string) (*bean2.ParsedValue, error) {
	return ParseCPUorMemoryValue(val)
}

func (m *MemoryUnitFactory) Validate(profileBean, defaultProfile *bean.ProfileBeanDto) error {
	// currently validating cpu and memory limits and reqs only
	var (
		memLimit *bean.ConfigurationBean
		memReq   *bean.ConfigurationBean
	)

	for _, platformConfigurations := range profileBean.Configurations {
		for _, configuration := range platformConfigurations {
			// get cpu limit and req
			switch configuration.Key {
			case bean.MEMORY_LIMIT:
				memLimit = configuration
			case bean.MEMORY_REQUEST:
				memReq = configuration
			}
		}
	}
	// validate mem
	err := validateMEM(memLimit, memReq)
	if err != nil {
		return err
	}
	return nil
}

func validateMEM(memLimit, memReq *bean.ConfigurationBean) error {
	memLimitUnitSuffix := bean2.MemoryUnitStr(memLimit.Unit)
	memReqUnitSuffix := bean2.MemoryUnitStr(memReq.Unit)
	memLimitUnit, ok := memLimitUnitSuffix.GetUnit()
	if !ok {
		return errors.New(fmt.Sprintf(bean.InvalidUnit, memLimit.Unit, memLimit.Key))
	}
	memReqUnit, ok := memReqUnitSuffix.GetUnit()
	if !ok {
		return errors.New(fmt.Sprintf(bean.InvalidUnit, memReq.Unit, memReq.Key))
	}

	// Use getTypedValue to retrieve appropriate types
	memLimitInterfaceVal, err := adapter.GetTypedValue(memLimit.Key, memLimit.Value)
	if err != nil {
		return errors.New(fmt.Sprintf(bean.InvalidTypeValue, memLimit.Key, memLimit.Value))
	}
	memLimitVal, ok := memLimitInterfaceVal.(float64)
	if !ok {
		return errors.New(fmt.Sprintf(bean.InvalidTypeValue, memLimit.Key, memLimit.Value))
	}

	memReqInterfaceVal, err := adapter.GetTypedValue(memReq.Key, memReq.Value)
	if err != nil {
		return errors.New(fmt.Sprintf(bean.InvalidTypeValue, memReq.Key, memReq.Value))
	}

	memReqVal, ok := memReqInterfaceVal.(float64)
	if !ok {
		return errors.New(fmt.Sprintf(bean.InvalidTypeValue, memReq.Key, memReq.Value))
	}

	if !validLimReq(memLimitVal, memLimitUnit.ConversionFactor, memReqVal, memReqUnit.ConversionFactor) {
		return errors.New(bean.MEMLimReqErrorCompErr)
	}
	return nil
}
