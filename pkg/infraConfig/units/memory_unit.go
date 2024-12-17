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

type MemoryUnitFactory struct {
	logger      *zap.SugaredLogger
	memoryUnits map[MemoryUnitStr]Unit
}

func NewMemoryUnitFactory(logger *zap.SugaredLogger) *MemoryUnitFactory {
	return &MemoryUnitFactory{
		logger:      logger,
		memoryUnits: getMemoryUnit(),
	}
}

func (m *MemoryUnitFactory) GetAllUnits() map[string]Unit {
	memoryUnits := m.memoryUnits
	units := make(map[string]Unit)
	for key, value := range memoryUnits {
		units[string(key)] = value
	}
	return units
}

func (m *MemoryUnitFactory) ParseValAndUnit(val string) (*ParsedValue, error) {
	return ParseCPUorMemoryValue(val)
}

type MemoryUnitStr string

const (
	MILLIBYTE MemoryUnitStr = "m"
	BYTE      MemoryUnitStr = "byte"
	KIBYTE    MemoryUnitStr = "Ki"
	MIBYTE    MemoryUnitStr = "Mi"
	GIBYTE    MemoryUnitStr = "Gi"
	TIBYTE    MemoryUnitStr = "Ti"
	PIBYTE    MemoryUnitStr = "Pi"
	EIBYTE    MemoryUnitStr = "Ei"
	KBYTE     MemoryUnitStr = "k"
	MBYTE     MemoryUnitStr = "M"
	GBYTE     MemoryUnitStr = "G"
	TBYTE     MemoryUnitStr = "T"
	PBYTE     MemoryUnitStr = "P"
	EBYTE     MemoryUnitStr = "E"
)

func (memoryUnitStr MemoryUnitStr) GetUnitSuffix() UnitType {
	switch memoryUnitStr {
	case BYTE:
		return Byte
	case KIBYTE:
		return KiByte
	case MIBYTE:
		return MiByte
	case GIBYTE:
		return GiByte
	case TIBYTE:
		return TiByte
	case PIBYTE:
		return PiByte
	case EIBYTE:
		return EiByte
	case KBYTE:
		return K
	case MBYTE:
		return M
	case GBYTE:
		return G
	case TBYTE:
		return T
	case PBYTE:
		return P
	case EBYTE:
		return E
	default:
		return Byte
	}
}

func (memoryUnitStr MemoryUnitStr) GetUnit() (Unit, bool) {
	memoryUnits := getMemoryUnit()
	memoryUnit, exists := memoryUnits[memoryUnitStr]
	return memoryUnit, exists
}

func (memoryUnitStr MemoryUnitStr) String() string {
	return string(memoryUnitStr)
}

func getMemoryUnit() map[MemoryUnitStr]Unit {
	return map[MemoryUnitStr]Unit{
		MILLIBYTE: {
			Name:             string(MILLIBYTE),
			ConversionFactor: 1e-3,
		},
		BYTE: {
			Name:             string(BYTE),
			ConversionFactor: 1,
		},
		KBYTE: {
			Name:             string(KBYTE),
			ConversionFactor: 1000,
		},
		MBYTE: {
			Name:             string(MBYTE),
			ConversionFactor: 1000000,
		},
		GBYTE: {
			Name:             string(GBYTE),
			ConversionFactor: 1000000000,
		},
		TBYTE: {
			Name:             string(TBYTE),
			ConversionFactor: 1000000000000,
		},
		PBYTE: {
			Name:             string(PBYTE),
			ConversionFactor: 1000000000000000,
		},
		EBYTE: {
			Name:             string(EBYTE),
			ConversionFactor: 1000000000000000000,
		},
		KIBYTE: {
			Name:             string(KIBYTE),
			ConversionFactor: 1024,
		},
		MIBYTE: {
			Name:             string(MIBYTE),
			ConversionFactor: 1024 * 1024,
		},
		GIBYTE: {
			Name:             string(GIBYTE),
			ConversionFactor: 1024 * 1024 * 1024,
		},
		TIBYTE: {
			Name:             string(TIBYTE),
			ConversionFactor: 1024 * 1024 * 1024 * 1024,
		},
		PIBYTE: {
			Name:             string(PIBYTE),
			ConversionFactor: 1024 * 1024 * 1024 * 1024 * 1024,
		},
		EIBYTE: {
			Name:             string(EIBYTE),
			ConversionFactor: 1024 * 1024 * 1024 * 1024 * 1024 * 1024,
		},
	}
}
