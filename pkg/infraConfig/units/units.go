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
	"fmt"
	util2 "github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/infraConfig/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/infraConfig/units/bean"
	"github.com/devtron-labs/devtron/util"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/resource"
	"net/http"
	"strconv"
	"strings"
)

type UnitStrService interface {
	GetUnitSuffix() bean2.UnitType
	GetUnit() (bean.Unit, bool)
	String() string
}

type PropertyType string

const (
	CPU    PropertyType = "CPU"
	MEMORY PropertyType = "MEMORY"
	TIME   PropertyType = "TIME"
)

type UnitService interface {
	GetAllUnits() map[string]bean.Unit
	ParseValAndUnit(val string) (*bean2.ParsedValue, error)
	Validate(profileBean, defaultProfile *bean.ProfileBeanDto) error
}

func NewUnitService(propertyType PropertyType, logger *zap.SugaredLogger) (UnitService, error) {
	switch propertyType {
	case CPU:
		return NewCPUUnitFactory(logger), nil
	case MEMORY:
		return NewMemoryUnitFactory(logger), nil
	case TIME:
		return NewTimeUnitFactory(logger), nil
	default:
		errMsg := fmt.Sprintf("invalid property type '%s'", propertyType)
		return nil, util2.NewApiError(http.StatusBadRequest, errMsg, errMsg)
	}
}

// ParseCPUorMemoryValue parses the quantity that has number values string and returns the value and unitType
// returns error if parsing fails
func ParseCPUorMemoryValue(quantity string) (*bean2.ParsedValue, error) {
	parsedValue := bean2.NewParsedValue()
	positive, _, num, denom, suffix, err := parseQuantityString(quantity)
	if err != nil {
		return parsedValue, err
	}
	if !positive {
		return parsedValue, errors.New("negative value not allowed for cpu limits")
	}
	valStr := num
	if denom != "" {
		valStr = num + "." + denom
	}

	val, err := strconv.ParseFloat(valStr, 64)

	// currently we are not supporting exponential values upto 2 decimals
	val = util.TruncateFloat(val, 2)
	return parsedValue.
		WithValueString(strconv.FormatFloat(val, 'f', -1, 64)).
		WithUnit(suffix), err
}

// parseQuantityString is a fast scanner for quantity values.
// this parsing is only for cpu and mem resources
func parseQuantityString(str string) (positive bool, value, num, denom, suffix string, err error) {
	positive = true
	pos := 0
	end := len(str)

	// handle leading sign
	if pos < end {
		switch str[0] {
		case '-':
			positive = false
			pos++
		case '+':
			pos++
		}
	}

	// strip leading zeros
Zeroes:
	for i := pos; ; i++ {
		if i >= end {
			num = "0"
			value = num
			return
		}
		switch str[i] {
		case '0':
			pos++
		default:
			break Zeroes
		}
	}

	// extract the numerator
Num:
	for i := pos; ; i++ {
		if i >= end {
			num = str[pos:end]
			value = str[0:end]
			return
		}
		switch str[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			num = str[pos:i]
			pos = i
			break Num
		}
	}

	// if we stripped all numerator positions, always return 0
	if len(num) == 0 {
		num = "0"
	}

	// handle a denominator
	if pos < end && str[pos] == '.' {
		pos++
	Denom:
		for i := pos; ; i++ {
			if i >= end {
				denom = str[pos:end]
				value = str[0:end]
				return
			}
			switch str[i] {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				denom = str[pos:i]
				pos = i
				break Denom
			}
		}
		// TODO: we currently allow 1.G, but we may not want to in the future.
		// if len(denom) == 0 {
		// 	err = ErrFormatWrong
		// 	return
		// }
	}
	value = str[0:pos]

	// grab the elements of the suffix
	suffixStart := pos
	for i := pos; ; i++ {
		if i >= end {
			suffix = str[suffixStart:end]
			return
		}
		if !strings.ContainsAny(str[i:i+1], "eEinumkKMGTP") {
			pos = i
			break
		}
	}
	if pos < end {
		switch str[pos] {
		case '-', '+':
			pos++
		}
	}
Suffix:
	for i := pos; ; i++ {
		if i >= end {
			suffix = str[suffixStart:end]
			return
		}
		switch str[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			break Suffix
		}
	}
	// we encountered a non decimal in the Suffix loop, but the last character
	// was not a valid exponent
	err = resource.ErrFormatWrong
	return
}
