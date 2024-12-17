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
	"strconv"
)

type TimeUnitFactory struct {
	logger    *zap.SugaredLogger
	timeUnits map[TimeUnitStr]Unit
}

func NewTimeUnitFactory(logger *zap.SugaredLogger) *TimeUnitFactory {
	return &TimeUnitFactory{
		logger:    logger,
		timeUnits: getTimeUnit(),
	}
}

func (t *TimeUnitFactory) GetAllUnits() map[string]Unit {
	timeUnits := t.timeUnits
	units := make(map[string]Unit)
	for key, value := range timeUnits {
		units[string(key)] = value
	}
	return units
}

func (t *TimeUnitFactory) ParseValAndUnit(val string) (*ParsedValue, error) {
	parsedValue := NewParsedValue()
	floatValue, err := strconv.ParseFloat(val, 64)
	if err != nil {
		t.logger.Errorw("Error while parsing value", "value", val, "error", err)
		return nil, err
	}
	parsedValue.WithValueFloat(floatValue).WithUnit(SecondStr.String())
	return parsedValue, nil
}

type TimeUnitStr string

const (
	SecondStr TimeUnitStr = "Seconds"
	MinuteStr TimeUnitStr = "Minutes"
	HourStr   TimeUnitStr = "Hours"
)

func (timeUnitStr TimeUnitStr) GetUnitSuffix() UnitType {
	switch timeUnitStr {
	case SecondStr:
		return Second
	case MinuteStr:
		return Minute
	case HourStr:
		return Hour
	default:
		return Second
	}
}

func getTimeUnit() map[TimeUnitStr]Unit {
	return map[TimeUnitStr]Unit{
		SecondStr: {
			Name:             string(SecondStr),
			ConversionFactor: 1,
		},
		MinuteStr: {
			Name:             string(MinuteStr),
			ConversionFactor: 60,
		},
		HourStr: {
			Name:             string(HourStr),
			ConversionFactor: 3600,
		},
	}
}

func (timeUnitStr TimeUnitStr) GetUnit() (Unit, bool) {
	timeUnits := getTimeUnit()
	timeUnit, exists := timeUnits[timeUnitStr]
	return timeUnit, exists
}

func (timeUnitStr TimeUnitStr) String() string {
	return string(timeUnitStr)
}
