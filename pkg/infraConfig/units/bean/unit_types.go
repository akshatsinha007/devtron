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

package bean

// memory units
// Ei, Pi, Ti, Gi, Mi, Ki
// E, P, T, G, M, k, m

type UnitType int

const (
	Byte      UnitType = 1
	KiByte    UnitType = 2 // 1024
	MiByte    UnitType = 3
	GiByte    UnitType = 4
	TiByte    UnitType = 5
	PiByte    UnitType = 6
	EiByte    UnitType = 7
	K         UnitType = 8 // 1000
	M         UnitType = 9
	G         UnitType = 10
	T         UnitType = 11
	P         UnitType = 12
	E         UnitType = 13
	Core      UnitType = 14 // CPU cores
	Milli     UnitType = 15
	Second    UnitType = 16
	Minute    UnitType = 17
	Hour      UnitType = 18
	MilliByte UnitType = 19
)

func (unitType UnitType) GetCPUUnitStr() CPUUnitStr {
	switch unitType {
	case Core:
		return CORE
	case Milli:
		return MILLI
	default:
		return CORE
	}
}

func (unitType UnitType) GetMemoryUnitStr() MemoryUnitStr {
	switch unitType {
	case MilliByte:
		return MILLIBYTE
	case Byte:
		return BYTE
	case KiByte:
		return KIBYTE
	case MiByte:
		return MIBYTE
	case GiByte:
		return GIBYTE
	case TiByte:
		return TIBYTE
	case PiByte:
		return PIBYTE
	case EiByte:
		return EIBYTE
	case K:
		return KBYTE
	case M:
		return MBYTE
	case G:
		return GBYTE
	case T:
		return TBYTE
	case P:
		return PBYTE
	case E:
		return EBYTE
	default:
		return BYTE
	}
}

func (unitType UnitType) GetTimeUnitStr() TimeUnitStr {
	switch unitType {
	case Second:
		return SecondStr
	case Minute:
		return MinuteStr
	case Hour:
		return HourStr
	default:
		return SecondStr
	}
}

type ParsedValue struct {
	valueString string
	unitType    string
}

func NewParsedValue() *ParsedValue {
	return &ParsedValue{}
}

func (p *ParsedValue) WithValueString(value string) *ParsedValue {
	p.valueString = value
	return p
}

func (p *ParsedValue) WithUnit(unit string) *ParsedValue {
	p.unitType = unit
	return p
}

func (p *ParsedValue) GetValueString() string {
	return p.valueString
}

func (p *ParsedValue) GetUnitType() string {
	return p.unitType
}
