/*
 * Copyright (c) 2020-2024. Devtron Inc.
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

package webhookHelm

import (
	"context"
	"fmt"
	"github.com/devtron-labs/devtron/api/helm-app/bean"
	bean2 "github.com/devtron-labs/devtron/api/helm-app/gRPC"
	client "github.com/devtron-labs/devtron/api/helm-app/service"
	helmBean "github.com/devtron-labs/devtron/api/helm-app/service/bean"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/pkg/attributes"
	bean3 "github.com/devtron-labs/devtron/pkg/attributes/bean"
	"github.com/devtron-labs/devtron/pkg/chartRepo"
	"github.com/devtron-labs/devtron/pkg/cluster"
	clientErrors "github.com/devtron-labs/devtron/pkg/errors"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"net/http"
)

const (
	DEFAULT_NAMESPACE   = "default"
	HELM_APP_DETAIL_URL = "%s/orchestrator/application/app?appId=%s"
)

type WebhookHelmService interface {
	CreateOrUpdateHelmApplication(ctx context.Context, request *HelmAppCreateUpdateRequest) (result interface{}, errorCode string, errorMessage string, statusCode int)
}

type WebhookHelmServiceImpl struct {
	logger                 *zap.SugaredLogger
	helmAppService         client.HelmAppService
	clusterService         cluster.ClusterService
	chartRepositoryService chartRepo.ChartRepositoryService
	attributesService      attributes.AttributesService
}

func NewWebhookHelmServiceImpl(logger *zap.SugaredLogger, helmAppService client.HelmAppService, clusterService cluster.ClusterService,
	chartRepositoryService chartRepo.ChartRepositoryService, attributesService attributes.AttributesService) *WebhookHelmServiceImpl {
	return &WebhookHelmServiceImpl{
		logger:                 logger,
		helmAppService:         helmAppService,
		clusterService:         clusterService,
		chartRepositoryService: chartRepositoryService,
		attributesService:      attributesService,
	}
}

func (impl WebhookHelmServiceImpl) CreateOrUpdateHelmApplication(ctx context.Context, request *HelmAppCreateUpdateRequest) (result interface{}, errorCode string, errorMessage string, statusCode int) {
	impl.logger.Infow("Request for create/update helm application from webhook", "request", request)

	// initialise clusterId
	var clusterId int

	// STEP-1 - get cluster info
	clusterName := request.ClusterName
	if len(clusterName) > 0 {
		cluster, err := impl.clusterService.FindOneActive(clusterName)
		if err != nil {
			impl.logger.Errorw("Error in getting cluster", "clusterName", clusterName, "err", err)
			if err == pg.ErrNoRows {
				return nil, common.ResourceNotFound, "cluster not found for given cluster name", http.StatusOK
			} else {
				return nil, common.InternalServerError, err.Error(), http.StatusInternalServerError
			}
		}
		clusterId = cluster.Id
	}

	// STEP-2 - set namespace as default if not supplied
	if len(request.Namespace) == 0 {
		request.Namespace = DEFAULT_NAMESPACE
	}

	// STEP-3 - get chart repository info
	if request.Chart.Repo.Identifier == nil {
		chartRepoName := request.Chart.Repo.Name
		chartRepo, err := impl.chartRepositoryService.GetChartRepoByName(chartRepoName)
		if err != nil {
			impl.logger.Errorw("Error in getting chart repo", "chartRepoName", chartRepoName, "err", err)
			return nil, common.InternalServerError, err.Error(), http.StatusInternalServerError
		}
		if chartRepo.Id == 0 {
			return nil, common.ResourceNotFound, "chart repository not found for given chart repo name", http.StatusOK
		}
		request.Chart.Repo.Identifier = &ChartRepoIdentifierSpec{
			Url:      chartRepo.Url,
			Username: chartRepo.UserName,
			Password: chartRepo.Password,
		}
	}

	// STEP-4 - build app identifier
	appIdentifier := &helmBean.AppIdentifier{
		ClusterId:   clusterId,
		Namespace:   request.Namespace,
		ReleaseName: request.ReleaseName,
	}

	// STEP-5 - check if the release is installed or not
	isInstalled, err := impl.helmAppService.IsReleaseInstalled(ctx, appIdentifier)
	if err != nil {
		impl.logger.Errorw("Error in checking if release is installed or not", "appIdentifier", appIdentifier, "err", err)
		return nil, common.InternalServerError, err.Error(), http.StatusInternalServerError
	}

	// STEP-6 install/update release
	chart := request.Chart
	chartRepo := request.Chart.Repo
	installReleaseRequest := &bean2.InstallReleaseRequest{
		ReleaseIdentifier: &bean2.ReleaseIdentifier{
			ReleaseName:      appIdentifier.ReleaseName,
			ReleaseNamespace: appIdentifier.Namespace,
		},
		ChartName:    chart.ChartName,
		ChartVersion: chart.ChartVersion,
		ValuesYaml:   request.ValuesOverrideYaml,
		ChartRepository: &bean2.ChartRepository{
			Name:     chartRepo.Name,
			Url:      chartRepo.Identifier.Url,
			Username: chartRepo.Identifier.Username,
			Password: chartRepo.Identifier.Password,
		},
	}
	if isInstalled {
		updateReleaseRequest := &bean.UpdateApplicationWithChartInfoRequestDto{
			InstallReleaseRequest: installReleaseRequest,
			SourceAppType:         bean.SOURCE_HELM_APP,
		}
		res, err := impl.helmAppService.UpdateApplicationWithChartInfo(ctx, clusterId, updateReleaseRequest)
		if err != nil {
			impl.logger.Errorw("Error in updating helm release", "appIdentifier", appIdentifier, "err", err)
			return nil, common.InternalServerError, err.Error(), http.StatusInternalServerError
		}
		if !res.GetSuccess() {
			return nil, common.UnknownError, "helm application update un-successful", http.StatusOK
		}
	} else {
		res, err := impl.helmAppService.InstallRelease(ctx, clusterId, installReleaseRequest)
		if err != nil {
			impl.logger.Errorw("Error in installing helm release", "appIdentifier", appIdentifier, "err", err)
			apiError := clientErrors.ConvertToApiError(err)
			if apiError != nil {
				err = apiError
			}
			return nil, common.InternalServerError, err.Error(), http.StatusInternalServerError
		}
		if !res.GetSuccess() {
			return nil, common.UnknownError, "helm application install un-successful", http.StatusOK
		}
	}

	// STEP-7 build app detail url (if error, then return success as operations has been completed already, just result is sent to be nil)
	hostUrlAttribute, err := impl.attributesService.GetByKey(bean3.HostUrlKey)
	if err != nil || hostUrlAttribute == nil {
		impl.logger.Errorw("error while getting host url attribute from DB", "error", err)
		return nil, "", "", http.StatusOK
	}
	appDetailUrl := fmt.Sprintf(HELM_APP_DETAIL_URL, hostUrlAttribute.Value, impl.helmAppService.EncodeAppId(appIdentifier))
	return appDetailUrl, "", "", http.StatusOK
}
