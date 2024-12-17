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
	"go.uber.org/zap"
)

type CPUUnitFactory struct {
	logger   *zap.SugaredLogger
	cpuUnits map[CPUUnitStr]Unit
}

func NewCPUUnitFactory(logger *zap.SugaredLogger) *CPUUnitFactory {
	return &CPUUnitFactory{
		logger:   logger,
		cpuUnits: getCPUUnit(),
	}
}

func (c *CPUUnitFactory) GetAllUnits() map[string]Unit {
	cpuUnits := c.cpuUnits
	units := make(map[string]Unit)
	for key, value := range cpuUnits {
		units[string(key)] = value
	}
	return units
}

func (c *CPUUnitFactory) ParseValAndUnit(val string) (*ParsedValue, error) {
	return ParseCPUorMemoryValue(val)
}

type CPUUnitStr string

const (
	CORE  CPUUnitStr = "Core"
	MILLI CPUUnitStr = "m"
)

func (cpuUnitStr CPUUnitStr) GetUnitSuffix() UnitType {
	switch cpuUnitStr {
	case CORE:
		return Core
	case MILLI:
		return Milli
	default:
		return Core
	}
}

func (cpuUnitStr CPUUnitStr) GetUnit() (Unit, bool) {
	cpuUnits := getCPUUnit()
	cpuUnit, exists := cpuUnits[cpuUnitStr]
	return cpuUnit, exists
}

func (cpuUnitStr CPUUnitStr) String() string {
	return string(cpuUnitStr)
}

func getCPUUnit() map[CPUUnitStr]Unit {
	return map[CPUUnitStr]Unit{
		MILLI: {
			Name:             string(MILLI),
			ConversionFactor: 1e-3,
		},
		CORE: {
			Name:             string(CORE),
			ConversionFactor: 1,
		},
	}
}
