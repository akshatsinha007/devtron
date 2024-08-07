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

package util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"math/rand"
	"testing"
)

func TestChartTemplateService(t *testing.T) {

	t.Run("getValues", func(t *testing.T) {
		logger, err := NewSugardLogger()
		assert.Nil(t, err)
		impl := ChartTemplateServiceImpl{
			logger: logger,
		}
		directory := "/scripts/devtron-reference-helm-charts/reference-chart_3-11-0"
		pipelineStrategyPath := "pipeline-values.yaml"
		values, err := impl.getValues(directory, pipelineStrategyPath)
		assert.Nil(t, err)
		assert.NotNil(t, values)
	})

	t.Run("buildChart", func(t *testing.T) {
		logger, err := NewSugardLogger()
		assert.Nil(t, err)
		impl := ChartTemplateServiceImpl{
			logger:     logger,
			randSource: rand.NewSource(0),
		}
		chartMetaData := &chart.Metadata{
			Name:    "sample-app",
			Version: "1.0.0",
		}
		refChartDir := "/scripts/devtron-reference-helm-charts/reference-chart_3-11-0"

		builtChartPath, err := impl.BuildChart(context.Background(), chartMetaData, refChartDir)
		assert.Nil(t, err)
		assert.DirExists(t, builtChartPath)

		isValidChart, err := chartutil.IsChartDir(builtChartPath)
		assert.Nil(t, err)
		assert.Equal(t, isValidChart, true)
	})

	t.Run("LoadChartInBytesWithDeleteFalse", func(t *testing.T) {
		logger, err := NewSugardLogger()
		assert.Nil(t, err)
		impl := ChartTemplateServiceImpl{
			logger:     logger,
			randSource: rand.NewSource(0),
		}
		chartMetaData := &chart.Metadata{
			Name:    "sample-app",
			Version: "1.0.0",
		}
		refChartDir := "/scripts/devtron-reference-helm-charts/reference-chart_3-11-0"

		builtChartPath, err := impl.BuildChart(context.Background(), chartMetaData, refChartDir)

		chartBytes, err := impl.LoadChartInBytes(builtChartPath, false)
		assert.Nil(t, err)

		chartBytesLen := len(chartBytes)
		assert.NotEqual(t, chartBytesLen, 0)

	})

	t.Run("LoadChartInBytesWithDeleteTrue", func(t *testing.T) {
		logger, err := NewSugardLogger()
		assert.Nil(t, err)
		impl := ChartTemplateServiceImpl{
			logger:     logger,
			randSource: rand.NewSource(0),
		}
		chartMetaData := &chart.Metadata{
			Name:    "sample-app",
			Version: "1.0.0",
		}
		refChartDir := "/scripts/devtron-reference-helm-charts/reference-chart_3-11-0"

		builtChartPath, err := impl.BuildChart(context.Background(), chartMetaData, refChartDir)

		chartBytes, err := impl.LoadChartInBytes(builtChartPath, true)
		assert.Nil(t, err)

		assert.NoDirExists(t, builtChartPath)

		chartBytesLen := len(chartBytes)
		assert.NotEqual(t, chartBytesLen, 0)

	})
}
