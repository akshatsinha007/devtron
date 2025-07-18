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

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/devtron-labs/devtron/api/bean/gitOps"
	bean3 "github.com/devtron-labs/devtron/api/helm-app/bean"
	bean4 "github.com/devtron-labs/devtron/api/helm-app/gRPC"
	openapi "github.com/devtron-labs/devtron/api/helm-app/openapiClient"
	"github.com/devtron-labs/devtron/api/helm-app/service"
	"github.com/devtron-labs/devtron/api/helm-app/service/bean"
	openapi2 "github.com/devtron-labs/devtron/api/openapi/openapiClient"
	"github.com/devtron-labs/devtron/client/argocdServer"
	"github.com/devtron-labs/devtron/internal/sql/repository/app"
	repository2 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/timelineStatus"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/workflow/cdWorkflow"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/appStore/adapter"
	appStoreBean "github.com/devtron-labs/devtron/pkg/appStore/bean"
	repository3 "github.com/devtron-labs/devtron/pkg/appStore/chartGroup/repository"
	appStoreDiscoverRepository "github.com/devtron-labs/devtron/pkg/appStore/discover/repository"
	installedAppAdapter "github.com/devtron-labs/devtron/pkg/appStore/installedApp/adapter"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/repository"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/EAMode"
	deployment2 "github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/EAMode/deployment"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/FullMode/deployment"
	bean2 "github.com/devtron-labs/devtron/pkg/appStore/installedApp/service/bean"
	"github.com/devtron-labs/devtron/pkg/cluster/environment"
	"github.com/devtron-labs/devtron/pkg/deployment/common"
	bean5 "github.com/devtron-labs/devtron/pkg/deployment/common/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/gitOps/config"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/go-pg/pg"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type AppStoreDeploymentService interface {
	InstallApp(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, ctx context.Context) (*appStoreBean.InstallAppVersionDTO, error)
	UpdateInstalledApp(ctx context.Context, installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*appStoreBean.InstallAppVersionDTO, error)
	DeleteInstalledApp(ctx context.Context, installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*appStoreBean.InstallAppVersionDTO, error)
	LinkHelmApplicationToChartStore(ctx context.Context, request *openapi.UpdateReleaseWithChartLinkingRequest, appIdentifier *bean.AppIdentifier, userId int32) (*openapi.UpdateReleaseResponse, bool, error)
	UpdateProjectHelmApp(updateAppRequest *appStoreBean.UpdateProjectHelmAppDTO) error
	// RollbackApplication Here we are handling the request of rollback using this common function
	RollbackApplication(ctx context.Context, request *openapi2.RollbackReleaseRequest, installedApp *appStoreBean.InstallAppVersionDTO, userId int32) (bool, error)
	GetDeploymentHistory(ctx context.Context, installedApp *appStoreBean.InstallAppVersionDTO) (*bean3.DeploymentHistoryAndInstalledAppInfo, error)
	GetDeploymentHistoryInfo(ctx context.Context, installedApp *appStoreBean.InstallAppVersionDTO, installedAppVersionHistoryId int) (*openapi.HelmAppDeploymentManifestDetail, error)
	InstallAppByHelm(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, ctx context.Context) (*appStoreBean.InstallAppVersionDTO, error)
	UpdatePreviousDeploymentStatusForAppStore(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, triggeredAt time.Time, err error) error
	MarkGitOpsInstalledAppsDeletedIfArgoAppIsDeleted(installedAppId, envId int) error
}

type AppStoreDeploymentServiceImpl struct {
	logger                               *zap.SugaredLogger
	installedAppRepository               repository.InstalledAppRepository
	installedAppService                  EAMode.InstalledAppDBService
	appStoreDeploymentDBService          AppStoreDeploymentDBService
	chartGroupDeploymentRepository       repository3.ChartGroupDeploymentRepository
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository
	appRepository                        app.AppRepository
	eaModeDeploymentService              deployment2.EAModeDeploymentService
	fullModeDeploymentService            deployment.FullModeDeploymentService
	fullModeFluxDeploymentService        deployment.FullModeFluxDeploymentService
	environmentService                   environment.EnvironmentService
	helmAppService                       service.HelmAppService
	installedAppRepositoryHistory        repository.InstalledAppVersionHistoryRepository
	deploymentTypeConfig                 *util2.DeploymentServiceTypeConfig
	aCDConfig                            *argocdServer.ACDConfig
	gitOpsConfigReadService              config.GitOpsConfigReadService
	deletePostProcessor                  DeletePostProcessor
	appStoreValidator                    AppStoreValidator
	deploymentConfigService              common.DeploymentConfigService
	OCIRegistryConfigRepository          repository2.OCIRegistryConfigRepository
}

func NewAppStoreDeploymentServiceImpl(logger *zap.SugaredLogger,
	installedAppRepository repository.InstalledAppRepository,
	installedAppService EAMode.InstalledAppDBService,
	appStoreDeploymentDBService AppStoreDeploymentDBService,
	chartGroupDeploymentRepository repository3.ChartGroupDeploymentRepository,
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository,
	appRepository app.AppRepository,
	eaModeDeploymentService deployment2.EAModeDeploymentService,
	fullModeDeploymentService deployment.FullModeDeploymentService,
	fullModeFluxDeploymentService deployment.FullModeFluxDeploymentService,
	environmentService environment.EnvironmentService,
	helmAppService service.HelmAppService,
	installedAppRepositoryHistory repository.InstalledAppVersionHistoryRepository,
	envVariables *util2.EnvironmentVariables,
	aCDConfig *argocdServer.ACDConfig,
	gitOpsConfigReadService config.GitOpsConfigReadService, deletePostProcessor DeletePostProcessor,
	appStoreValidator AppStoreValidator,
	deploymentConfigService common.DeploymentConfigService,
	OCIRegistryConfigRepository repository2.OCIRegistryConfigRepository) *AppStoreDeploymentServiceImpl {

	return &AppStoreDeploymentServiceImpl{
		logger:                               logger,
		installedAppRepository:               installedAppRepository,
		installedAppService:                  installedAppService,
		appStoreDeploymentDBService:          appStoreDeploymentDBService,
		chartGroupDeploymentRepository:       chartGroupDeploymentRepository,
		appStoreApplicationVersionRepository: appStoreApplicationVersionRepository,
		appRepository:                        appRepository,
		eaModeDeploymentService:              eaModeDeploymentService,
		fullModeDeploymentService:            fullModeDeploymentService,
		fullModeFluxDeploymentService:        fullModeFluxDeploymentService,
		environmentService:                   environmentService,
		helmAppService:                       helmAppService,
		installedAppRepositoryHistory:        installedAppRepositoryHistory,
		deploymentTypeConfig:                 envVariables.DeploymentServiceTypeConfig,
		aCDConfig:                            aCDConfig,
		gitOpsConfigReadService:              gitOpsConfigReadService,
		deletePostProcessor:                  deletePostProcessor,
		appStoreValidator:                    appStoreValidator,
		deploymentConfigService:              deploymentConfigService,
		OCIRegistryConfigRepository:          OCIRegistryConfigRepository,
	}
}

func (impl *AppStoreDeploymentServiceImpl) InstallApp(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, ctx context.Context) (*appStoreBean.InstallAppVersionDTO, error) {

	dbConnection := impl.installedAppRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()
	//step 1 db operation initiated
	installAppVersionRequest, err = impl.appStoreDeploymentDBService.AppStoreDeployOperationDB(installAppVersionRequest, tx, appStoreBean.INSTALL_APP_REQUEST)
	if err != nil {
		impl.logger.Errorw(" error", "err", err)
		return nil, err
	}

	//checking if namespace exists or not
	clusterIdToNsMap := map[int]string{
		installAppVersionRequest.ClusterId: installAppVersionRequest.Namespace,
	}
	err = impl.helmAppService.CheckIfNsExistsForClusterIds(clusterIdToNsMap)
	if err != nil {
		return nil, err
	}
	installedAppDeploymentAction := adapter.NewInstalledAppDeploymentAction(installAppVersionRequest.DeploymentAppType)

	if util.IsAcdApp(installAppVersionRequest.DeploymentAppType) || util.IsManifestDownload(installAppVersionRequest.DeploymentAppType) {
		_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_INITIATED, "Deployment initiated successfully.", time.Now(), tx)
	}

	if util.IsManifestDownload(installAppVersionRequest.DeploymentAppType) {
		_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_MANIFEST_GENERATED, "Manifest generated successfully.", time.Now(), tx)
	}

	var gitOpsResponse *bean2.AppStoreGitOpsResponse
	if installedAppDeploymentAction.PerformGitOps {
		appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(installAppVersionRequest.AppStoreVersion)
		if err != nil {
			impl.logger.Errorw("fetching error", "err", err)
			return nil, err
		}
		manifest, err := impl.fullModeDeploymentService.GenerateManifest(installAppVersionRequest, appStoreAppVersion)
		if err != nil {
			impl.logger.Errorw("error in performing manifest and git operations", "err", err)
			return nil, err
		}
		err = impl.fullModeDeploymentService.CreateArgoRepoSecretIfNeeded(appStoreAppVersion)
		if err != nil {
			impl.logger.Errorw("error in creating argo app repository secret", "appStoreApplicationVersionId", appStoreAppVersion.Id, "err", err)
			return nil, err
		}
		gitOpsResponse, err = impl.fullModeDeploymentService.GitOpsOperations(manifest, installAppVersionRequest)
		if err != nil {
			impl.logger.Errorw("error in doing gitops operation", "err", err)
			if util.IsAcdApp(installAppVersionRequest.DeploymentAppType) {
				_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_GIT_COMMIT_FAILED, fmt.Sprintf("Git commit failed - %v", err), time.Now(), tx)
			}
			return nil, err
		}
		if util.IsAcdApp(installAppVersionRequest.DeploymentAppType) {
			_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_GIT_COMMIT, timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_GIT_COMMIT, time.Now(), tx)
			if impl.aCDConfig.IsManualSyncEnabled() {
				_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_ARGOCD_SYNC_INITIATED, timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_SYNC_INITIATED, time.Now(), tx)
			}
		} else if util.IsFluxApp(installAppVersionRequest.DeploymentAppType) {
			_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(installAppVersionRequest, timelineStatus.TIMELINE_STATUS_GIT_COMMIT, timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_GIT_COMMIT, time.Now(), tx)
		}
		installAppVersionRequest.GitHash = gitOpsResponse.GitHash
		if len(installAppVersionRequest.GitHash) > 0 {
			err = impl.installedAppRepositoryHistory.UpdateGitHash(installAppVersionRequest.InstalledAppVersionHistoryId, gitOpsResponse.GitHash, tx)
			if err != nil {
				impl.logger.Errorw("error in updating git hash ", "err", err)
				return nil, err
			}
		}
	}

	if util2.IsBaseStack() || util2.IsHelmApp(installAppVersionRequest.AppOfferingMode) || util.IsHelmApp(installAppVersionRequest.DeploymentAppType) {
		installAppVersionRequest, err = impl.eaModeDeploymentService.InstallApp(installAppVersionRequest, nil, ctx, tx)
	} else if util.IsAcdApp(installAppVersionRequest.DeploymentAppType) {
		if gitOpsResponse == nil && gitOpsResponse.ChartGitAttribute != nil {
			return nil, errors.New("service err, Error in git operations")
		}
		installAppVersionRequest, err = impl.fullModeDeploymentService.InstallApp(installAppVersionRequest, gitOpsResponse.ChartGitAttribute, ctx, tx)
	} else if util.IsFluxApp(installAppVersionRequest.DeploymentAppType) {
		if gitOpsResponse == nil && gitOpsResponse.ChartGitAttribute != nil {
			return nil, errors.New("service err, Error in git operations")
		}
		installAppVersionRequest, err = impl.fullModeFluxDeploymentService.InstallApp(installAppVersionRequest, gitOpsResponse.ChartGitAttribute, ctx, tx)
	}
	if err != nil {
		return nil, err
	}
	err = tx.Commit()

	err = impl.appStoreDeploymentDBService.InstallAppPostDbOperation(installAppVersionRequest)
	if err != nil {
		return nil, err
	}

	return installAppVersionRequest, nil
}

func (impl *AppStoreDeploymentServiceImpl) DeleteInstalledApp(ctx context.Context, installAppVersionRequest *appStoreBean.InstallAppVersionDTO) (*appStoreBean.InstallAppVersionDTO, error) {
	installAppVersionRequest.InstalledAppDeleteResponse = &appStoreBean.InstalledAppDeleteResponseDTO{
		DeleteInitiated:  false,
		ClusterReachable: true,
	}
	dbConnection := impl.installedAppRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()

	environment, err := impl.environmentService.GetExtendedEnvBeanById(installAppVersionRequest.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return nil, err
	}
	if len(environment.ErrorInConnecting) > 0 {
		installAppVersionRequest.InstalledAppDeleteResponse.ClusterReachable = false
		installAppVersionRequest.InstalledAppDeleteResponse.ClusterName = environment.ClusterName
	}

	app, err := impl.appRepository.FindById(installAppVersionRequest.AppId)
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, fmt.Errorf("App not found in database")
		}
		return nil, err
	}
	model, err := impl.installedAppRepository.GetInstalledApp(installAppVersionRequest.InstalledAppId)
	if err != nil {
		impl.logger.Errorw("error in fetching installed app", "id", installAppVersionRequest.InstalledAppId, "err", err)
		return nil, err
	}

	deploymentConfig, err := impl.deploymentConfigService.GetConfigForHelmApps(model.AppId, model.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in getiting deployment config db object by appId and envId", "appId", model.AppId, "envId", model.EnvironmentId, "err", err)
		return nil, err
	}
	if installAppVersionRequest.AcdPartialDelete == true {
		if !util2.IsBaseStack() && !util2.IsHelmApp(app.AppOfferingMode) && !util.IsHelmApp(deploymentConfig.DeploymentAppType) {
			if !installAppVersionRequest.InstalledAppDeleteResponse.ClusterReachable {
				impl.logger.Errorw("cluster connection error", "err", environment.ErrorInConnecting)
				if !installAppVersionRequest.NonCascadeDelete {
					return installAppVersionRequest, nil
				}
			}
			err = impl.fullModeDeploymentService.DeleteACDAppObject(ctx, app.AppName, environment.Environment, installAppVersionRequest)
		}
		if err != nil {
			impl.logger.Errorw("error on delete installed app", "err", err)
			return nil, err
		}
		model.DeploymentAppDeleteRequest = true
		model.UpdatedBy = installAppVersionRequest.UserId
		model.UpdatedOn = time.Now()
		_, err = impl.installedAppRepository.UpdateInstalledApp(model, tx)
		if err != nil {
			impl.logger.Errorw("error while creating install app", "error", err)
			return nil, err
		}
	} else {
		//soft delete app
		app.Active = false
		app.UpdatedBy = installAppVersionRequest.UserId
		app.UpdatedOn = time.Now()
		err = impl.appRepository.UpdateWithTxn(app, tx)
		if err != nil {
			impl.logger.Errorw("error in update entity ", "entity", app)
			return nil, err
		}

		impl.deletePostProcessor.Process(app, installAppVersionRequest)

		// soft delete install app
		model.Active = false
		model.UpdatedBy = installAppVersionRequest.UserId
		model.UpdatedOn = time.Now()
		_, err = impl.installedAppRepository.UpdateInstalledApp(model, tx)
		if err != nil {
			impl.logger.Errorw("error while creating install app", "error", err)
			return nil, err
		}
		models, err := impl.installedAppRepository.GetInstalledAppVersionByInstalledAppId(installAppVersionRequest.InstalledAppId)
		if err != nil {
			impl.logger.Errorw("error while fetching install app versions", "error", err)
			return nil, err
		}

		// soft delete install app versions
		for _, item := range models {
			item.Active = false
			item.UpdatedBy = installAppVersionRequest.UserId
			item.UpdatedOn = time.Now()
			_, err = impl.installedAppRepository.UpdateInstalledAppVersion(item, tx)
			if err != nil {
				impl.logger.Errorw("error while fetching from db", "error", err)
				return nil, err
			}
		}

		// soft delete chart-group deployment
		chartGroupDeployment, err := impl.chartGroupDeploymentRepository.FindByInstalledAppId(model.Id)
		if err != nil && err != pg.ErrNoRows {
			impl.logger.Errorw("error while fetching chart group deployment", "error", err)
			return nil, err
		}
		if chartGroupDeployment.Id != 0 {
			chartGroupDeployment.Deleted = true
			_, err = impl.chartGroupDeploymentRepository.Update(chartGroupDeployment, tx)
			if err != nil {
				impl.logger.Errorw("error while updating chart group deployment", "error", err)
				return nil, err
			}
		}

		if util2.IsBaseStack() || util2.IsHelmApp(app.AppOfferingMode) || util.IsHelmApp(deploymentConfig.DeploymentAppType) {
			// there might be a case if helm release gets uninstalled from helm cli.
			//in this case on deleting the app from API, it should not give error as it should get deleted from db, otherwise due to delete error, db does not get clean
			// so in helm, we need to check first if the release exists or not, if exists then only delete
			err = impl.eaModeDeploymentService.DeleteInstalledApp(ctx, app.AppName, environment.Environment, installAppVersionRequest, model, tx)
		} else {
			err = impl.fullModeDeploymentService.DeleteInstalledApp(ctx, app.AppName, environment.Environment, installAppVersionRequest, model, tx)
		}
		if err != nil {
			impl.logger.Errorw("error on delete installed app", "err", err)
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		impl.logger.Errorw("error in commit db transaction on delete", "err", err)
		return nil, err
	}
	installAppVersionRequest.InstalledAppDeleteResponse.DeleteInitiated = true
	return installAppVersionRequest, nil
}

func (impl *AppStoreDeploymentServiceImpl) LinkHelmApplicationToChartStore(ctx context.Context, request *openapi.UpdateReleaseWithChartLinkingRequest,
	appIdentifier *bean.AppIdentifier, userId int32) (*openapi.UpdateReleaseResponse, bool, error) {

	impl.logger.Infow("Linking helm application to chart store", "appId", request.GetAppId())

	// check if chart repo is active starts
	isChartRepoActive, err := impl.appStoreDeploymentDBService.IsChartProviderActive(int(request.GetAppStoreApplicationVersionId()))
	if err != nil {
		impl.logger.Errorw("Error in checking if chart repo is active or not", "err", err)
		return nil, isChartRepoActive, err
	}
	if !isChartRepoActive {
		return nil, isChartRepoActive, nil
	}
	// check if chart repo is active ends

	// STEP-1 check if the app is installed or not
	isInstalled, err := impl.helmAppService.IsReleaseInstalled(ctx, appIdentifier)
	if err != nil {
		impl.logger.Errorw("error while checking if the release is installed", "error", err)
		return nil, isChartRepoActive, err
	}
	if !isInstalled {
		return nil, isChartRepoActive, errors.New("release is not installed. so can not be updated")
	}
	// STEP-1 ends

	// Initialise bean
	installAppVersionRequestDto := &appStoreBean.InstallAppVersionDTO{
		AppName:            appIdentifier.GetUniqueAppNameIdentifier(),
		UserId:             userId,
		AppOfferingMode:    util2.SERVER_MODE_HYPERION,
		ClusterId:          appIdentifier.ClusterId,
		Namespace:          appIdentifier.Namespace,
		AppStoreVersion:    int(request.GetAppStoreApplicationVersionId()),
		ValuesOverrideYaml: request.GetValuesYaml(),
		ReferenceValueId:   int(request.GetReferenceValueId()),
		ReferenceValueKind: request.GetReferenceValueKind(),
		DeploymentAppType:  util.PIPELINE_DEPLOYMENT_TYPE_HELM,
		DisplayName:        appIdentifier.ReleaseName,
		IsChartLinkRequest: true,
	}

	// STEP-2 InstallApp with only DB operations
	// STEP-3 update APP with chart info
	res, err := impl.linkHelmApplicationToChartStore(installAppVersionRequestDto, ctx)
	if err != nil {
		impl.logger.Errorw("error while linking helm app with chart store", "error", err)
		return nil, isChartRepoActive, err
	}
	// STEP-2 and STEP-3 ends

	return res, isChartRepoActive, nil
}

func isExternalHelmApp(appId string) bool {
	// for external helm apps, updateAppRequest.AppId is of the form clusterId|namespace|displayAppName
	return len(strings.Split(appId, "|")) > 1
}

func (impl *AppStoreDeploymentServiceImpl) UpdateProjectHelmApp(updateAppRequest *appStoreBean.UpdateProjectHelmAppDTO) error {
	var appName string
	var displayName string
	appName = updateAppRequest.AppName
	if isExternalHelmApp(updateAppRequest.AppId) {
		appIdentifier, err := impl.helmAppService.DecodeAppId(updateAppRequest.AppId)
		if err != nil {
			impl.logger.Errorw("error in decoding app id for external helm apps", "err", err)
			return err
		}
		appName = appIdentifier.GetUniqueAppNameIdentifier()
		displayName = updateAppRequest.AppName
	} else {
		//in case the external app is linked, then it's unique identifier is set in app_name col. hence retrieving appName
		//for this case, although this will also handle the case for non-external apps
		appNameUniqueIdentifier := impl.getAppNameForInstalledApp(updateAppRequest.InstalledAppId)
		if len(appNameUniqueIdentifier) > 0 {
			appName = appNameUniqueIdentifier
		}
	}
	impl.logger.Infow("update helm project request", updateAppRequest)
	err := impl.appStoreDeploymentDBService.UpdateProjectForHelmApp(appName, displayName, updateAppRequest.TeamId, updateAppRequest.UserId)
	if err != nil {
		impl.logger.Errorw("error in linking project to helm app", "appName", updateAppRequest.AppName, "err", err)
		return err
	}
	return nil
}

func (impl *AppStoreDeploymentServiceImpl) RollbackApplication(ctx context.Context, request *openapi2.RollbackReleaseRequest,
	installedApp *appStoreBean.InstallAppVersionDTO, userId int32) (bool, error) {
	var err error
	var success bool
	upgradeRequest := installedApp.NewInstalledAppVersionRequestDTO(userId, installedApp.InstalledAppId)

	//in case of externally cli helm apps, we set its installedAppId 0 in the request body, using which we are deciding the flow
	if installedApp.IsExternalCliApp() {
		installedApp, success, err = impl.eaModeDeploymentService.RollbackRelease(ctx, installedApp, request.GetVersion())
		if err != nil {
			impl.logger.Errorw("error while rollback helm release", "error", err)
			return false, err
		}
	} else {
		//we are fetching the values from installed app version history table (iavh) which is further used as different
		installedAppVersionHistory, err := impl.installedAppRepositoryHistory.GetInstalledAppVersionHistory(int(request.GetVersion()))
		if err != nil {
			impl.logger.Errorw("error while fetching installed app version history from Db", "request", request, "error", err)
			return false, err
		}
		installedAppVersionDTO, err := impl.installedAppService.GetInstalledAppVersionByIdIncludeDeleted(installedAppVersionHistory.InstalledAppVersionId, userId)
		if err != nil {
			impl.logger.Errorw("error while fetching installed app version detail from Db", "installedAppVersion", installedAppVersionHistory.InstalledAppVersionId, "error", err)
			return false, err
		}
		adapter.UpdateRequestDTOForRollback(upgradeRequest, installedApp, installedAppVersionDTO, installedAppVersionHistory)

		upgradeRequest, err = impl.UpdateInstalledApp(ctx, upgradeRequest)
		if err != nil {
			impl.logger.Errorw("error while performing update to the previous version", "upgradeRequest", upgradeRequest, "error", err)
			return false, err
		}
		success = true
	}
	if !success {
		return false, fmt.Errorf("rollback request failed")
	}
	return success, err
}

func (impl *AppStoreDeploymentServiceImpl) GetDeploymentHistory(ctx context.Context, installedApp *appStoreBean.InstallAppVersionDTO) (*bean3.DeploymentHistoryAndInstalledAppInfo, error) {
	newCtx, span := otel.Tracer("orchestrator").Start(ctx, "AppStoreDeploymentServiceImpl.GetDeploymentHistory")
	defer span.End()
	result := &bean3.DeploymentHistoryAndInstalledAppInfo{}
	var err error
	if util2.IsHelmApp(installedApp.AppOfferingMode) {
		deploymentHistory, err := impl.eaModeDeploymentService.GetDeploymentHistory(newCtx, installedApp)
		if err != nil {
			impl.logger.Errorw("error while getting deployment history", "error", err)
			return nil, err
		}
		result.DeploymentHistory = deploymentHistory.GetDeploymentHistory()
	} else {
		deploymentHistory, err := impl.fullModeDeploymentService.GetDeploymentHistory(newCtx, installedApp)
		if err != nil {
			impl.logger.Errorw("error while getting deployment history", "error", err)
			return nil, err
		}
		result.DeploymentHistory = deploymentHistory.GetDeploymentHistory()
	}

	if installedApp.InstalledAppId > 0 {
		result.InstalledAppInfo = &bean3.InstalledAppInfo{
			AppId:                 installedApp.AppId,
			EnvironmentName:       installedApp.EnvironmentName,
			AppOfferingMode:       installedApp.AppOfferingMode,
			InstalledAppId:        installedApp.InstalledAppId,
			InstalledAppVersionId: installedApp.InstalledAppVersionId,
			AppStoreChartId:       installedApp.InstallAppVersionChartDTO.AppStoreChartId,
			ClusterId:             installedApp.ClusterId,
			EnvironmentId:         installedApp.EnvironmentId,
			DeploymentType:        installedApp.DeploymentAppType,
			HelmPackageName: adapter.GetGeneratedHelmPackageName(
				installedApp.AppName,
				installedApp.EnvironmentName,
				installedApp.UpdatedOn),
		}
	}

	return result, err
}

func (impl *AppStoreDeploymentServiceImpl) GetDeploymentHistoryInfo(ctx context.Context, installedApp *appStoreBean.InstallAppVersionDTO, version int) (*openapi.HelmAppDeploymentManifestDetail, error) {
	//var result interface{}
	result := &openapi.HelmAppDeploymentManifestDetail{}
	var err error
	if util2.IsHelmApp(installedApp.AppOfferingMode) {
		_, span := otel.Tracer("orchestrator").Start(ctx, "eaModeDeploymentService.GetDeploymentHistoryInfo")
		result, err = impl.eaModeDeploymentService.GetDeploymentHistoryInfo(ctx, installedApp, int32(version))
		span.End()
		if err != nil {
			impl.logger.Errorw("error while getting deployment history info", "error", err)
			return nil, err
		}
	} else {
		_, span := otel.Tracer("orchestrator").Start(ctx, "fullModeDeploymentService.GetDeploymentHistoryInfo")
		result, err = impl.fullModeDeploymentService.GetDeploymentHistoryInfo(ctx, installedApp, int32(version))
		span.End()
		if err != nil {
			impl.logger.Errorw("error while getting deployment history info", "error", err)
			return nil, err
		}
	}
	return result, err
}

func (impl *AppStoreDeploymentServiceImpl) updateInstalledApp(ctx context.Context, upgradeAppRequest *appStoreBean.InstallAppVersionDTO, tx *pg.Tx) (*appStoreBean.InstallAppVersionDTO, error) {
	installedApp, err := impl.installedAppService.GetInstalledAppById(upgradeAppRequest.InstalledAppId)
	if err != nil {
		impl.logger.Errorw("error in fetching installed app by id", "installedAppId", upgradeAppRequest.InstalledAppId, "err", err)
		return nil, err
	}
	//checking if ns exists or not
	clusterIdToNsMap := map[int]string{
		installedApp.Environment.ClusterId: installedApp.Environment.Namespace,
	}

	deploymentConfig, err := impl.deploymentConfigService.GetAndMigrateConfigIfAbsentForHelmApp(installedApp.AppId, installedApp.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("error in getting deploymentConfig by appId and envId", "appId", installedApp.AppId, "envId", installedApp.EnvironmentId, "err", err)
		return nil, err
	}

	err = impl.helmAppService.CheckIfNsExistsForClusterIds(clusterIdToNsMap)
	if err != nil {
		impl.logger.Errorw("error in checking if namespace exists or not", "clusterId",
			installedApp.Environment.ClusterId, "namespace", installedApp.Environment.Namespace, "err", err)
		return nil, err
	}
	upgradeAppRequest.UpdateDeploymentAppType(deploymentConfig.DeploymentAppType)

	installedAppDeploymentAction := adapter.NewInstalledAppDeploymentAction(deploymentConfig.DeploymentAppType)
	// migrate installedApp.GitOpsRepoName to installedApp.GitOpsRepoUrl
	if (util.IsAcdApp(deploymentConfig.DeploymentAppType) || util.IsFluxApp(deploymentConfig.DeploymentAppType)) &&
		len(deploymentConfig.GetRepoURL()) == 0 {
		gitRepoUrl, err := impl.fullModeDeploymentService.GetAcdAppGitOpsRepoURL(installedApp.App.AppName, installedApp.Environment.Name)
		if err != nil {
			impl.logger.Errorw("error in GitOps repository url migration", "err", err)
			return nil, err
		}
		deploymentConfig.SetRepoURL(gitRepoUrl)
		//installedApp.GitOpsRepoUrl = gitRepoUrl
		installedApp.GitOpsRepoName = impl.gitOpsConfigReadService.GetGitOpsRepoNameFromUrl(gitRepoUrl)
	}
	// migration ends

	var installedAppVersion *repository.InstalledAppVersions

	// mark previous versions of chart as inactive if chart or version is updated
	isChartChanged := false   // flag for keeping track if chart is updated by user or not
	isVersionChanged := false // flag for keeping track if version of chart is upgraded

	//if chart is changed, then installedAppVersion id is sent as 0 from front-end
	if upgradeAppRequest.Id == 0 {
		isChartChanged = true
		err = impl.installedAppService.MarkInstalledAppVersionsInactiveByInstalledAppId(upgradeAppRequest.InstalledAppId, upgradeAppRequest.UserId, tx)
		if err != nil {
			return nil, err
		}
	} else {
		installedAppVersion, err = impl.installedAppRepository.GetInstalledAppVersion(upgradeAppRequest.Id)
		if err != nil {
			impl.logger.Errorw("error in fetching installedAppVersion by upgradeAppRequest id ", "err", err)
			return nil, fmt.Errorf("The values are outdated. Please make your changes to the latest version and try again.")
		}
		// version is upgraded if appStoreApplication version from request payload is not equal to installed app version saved in DB
		if installedAppVersion.AppStoreApplicationVersionId != upgradeAppRequest.AppStoreVersion {
			isVersionChanged = true
			err = impl.installedAppService.MarkInstalledAppVersionModelInActive(installedAppVersion, upgradeAppRequest.UserId, tx)
		}
	}

	appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(upgradeAppRequest.AppStoreVersion)
	if err != nil {
		impl.logger.Errorw("fetching error", "err", err)
		return nil, err
	}
	// create new entry for installed app version if chart or version is changed
	if isChartChanged || isVersionChanged {
		installedAppVersion, err = impl.installedAppService.CreateInstalledAppVersion(upgradeAppRequest, tx)
		if err != nil {
			impl.logger.Errorw("error in creating installed app version", "err", err)
			return nil, err
		}
	} else if installedAppVersion != nil && installedAppVersion.Id != 0 {
		adapter.UpdateInstalledAppVersionModel(installedAppVersion, upgradeAppRequest)
		installedAppVersion, err = impl.installedAppService.UpdateInstalledAppVersion(installedAppVersion, upgradeAppRequest, tx)
		if err != nil {
			impl.logger.Errorw("error in creating installed app version", "err", err)
			return nil, err
		}
	}
	// populate the related model data into repository.InstalledAppVersions
	// related tables: repository.InstalledApps AND appStoreDiscoverRepository.AppStoreApplicationVersion
	installedAppVersion.AppStoreApplicationVersion = *appStoreAppVersion
	installedAppVersion.InstalledApp = *installedApp

	// populate appStoreBean.InstallAppVersionDTO from the DB models
	upgradeAppRequest.Id = installedAppVersion.Id
	upgradeAppRequest.InstalledAppVersionId = installedAppVersion.Id
	adapter.UpdateInstallAppDetails(upgradeAppRequest, installedApp, deploymentConfig)
	adapter.UpdateAppDetails(upgradeAppRequest, &installedApp.App)
	environment, err := impl.environmentService.GetExtendedEnvBeanById(installedApp.EnvironmentId)
	if err != nil {
		impl.logger.Errorw("fetching environment error", "envId", installedApp.EnvironmentId, "err", err)
		return nil, err
	}
	adapter.UpdateAdditionalEnvDetails(upgradeAppRequest, environment)

	helmInstallConfigDTO := appStoreBean.HelmReleaseStatusConfig{
		InstallAppVersionHistoryId: 0,
		Message:                    "Install initiated",
		IsReleaseInstalled:         false,
		ErrorInInstallation:        false,
	}
	installedAppVersionHistory, err := adapter.NewInstallAppVersionHistoryModel(upgradeAppRequest, cdWorkflow.WorkflowInProgress, helmInstallConfigDTO)
	_, err = impl.installedAppRepositoryHistory.CreateInstalledAppVersionHistory(installedAppVersionHistory, tx)
	if err != nil {
		impl.logger.Errorw("error while creating installed app version history for updating installed app", "error", err)
		return nil, err
	}
	upgradeAppRequest.InstalledAppVersionHistoryId = installedAppVersionHistory.Id
	_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(upgradeAppRequest, timelineStatus.TIMELINE_STATUS_DEPLOYMENT_INITIATED, "Deployment initiated successfully.", time.Now(), tx)

	if util.IsManifestDownload(upgradeAppRequest.DeploymentAppType) {
		_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(upgradeAppRequest, timelineStatus.TIMELINE_STATUS_MANIFEST_GENERATED, "Manifest generated successfully.", time.Now(), tx)
	}
	// gitOps operation
	monoRepoMigrationRequired := false
	gitOpsResponse := &bean2.AppStoreGitOpsResponse{}

	if installedAppDeploymentAction.PerformGitOps {
		// manifest contains ChartRepoName where the valuesConfig and requirementConfig files will get committed
		// and that gitOpsRepoUrl is extracted from db inside GenerateManifest func and not from the current
		// orchestrator cm prefix and appName.
		manifest, err := impl.fullModeDeploymentService.GenerateManifest(upgradeAppRequest, appStoreAppVersion)
		if err != nil {
			impl.logger.Errorw("error in generating manifest for helm apps", "installedAppVersionHistoryId", upgradeAppRequest.InstalledAppVersionHistoryId, "err", err)
			_ = impl.appStoreDeploymentDBService.UpdateInstalledAppVersionHistoryStatus(
				upgradeAppRequest.InstalledAppVersionHistoryId,
				installedAppAdapter.FailedStatusUpdateOption(upgradeAppRequest.UserId, err),
			)
			return nil, err
		}
		err = impl.fullModeDeploymentService.CreateArgoRepoSecretIfNeeded(appStoreAppVersion)
		if err != nil {
			impl.logger.Errorw("error in creating argo app repository secret", "appStoreApplicationVersionId", appStoreAppVersion.Id, "err", err)
			return nil, err
		}
		// required if gitOps repo name is changed, gitOps repo name will change if env variable which we use as suffix changes
		monoRepoMigrationRequired = impl.checkIfMonoRepoMigrationRequired(installedApp, deploymentConfig)
		argocdAppName := util2.BuildDeployedAppName(installedApp.App.AppName, installedApp.Environment.Name)
		upgradeAppRequest.ACDAppName = argocdAppName

		var gitOpsErr error
		gitOpsResponse, gitOpsErr = impl.fullModeDeploymentService.UpdateAppGitOpsOperations(manifest, upgradeAppRequest, monoRepoMigrationRequired, isChartChanged || isVersionChanged)
		if gitOpsErr != nil {
			impl.logger.Errorw("error in performing GitOps operation", "err", gitOpsErr)
			_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(upgradeAppRequest, timelineStatus.TIMELINE_STATUS_GIT_COMMIT_FAILED, fmt.Sprintf("Git commit failed - %v", gitOpsErr), time.Now(), tx)
			return nil, gitOpsErr
		}

		upgradeAppRequest.GitHash = gitOpsResponse.GitHash
		_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(upgradeAppRequest, timelineStatus.TIMELINE_STATUS_GIT_COMMIT, timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_GIT_COMMIT, time.Now(), tx)
		if installedAppDeploymentAction.PerformACDDeployment && impl.aCDConfig.IsManualSyncEnabled() { //if acd then only save manifest timeline, filtering flux through this check
			_ = impl.fullModeDeploymentService.SaveTimelineForHelmApps(upgradeAppRequest, timelineStatus.TIMELINE_STATUS_ARGOCD_SYNC_INITIATED, timelineStatus.TIMELINE_DESCRIPTION_ARGOCD_SYNC_INITIATED, time.Now(), tx)
		}
		installedAppVersionHistory.GitHash = gitOpsResponse.GitHash
		_, err = impl.installedAppRepositoryHistory.UpdateInstalledAppVersionHistory(installedAppVersionHistory, tx)
		if err != nil {
			impl.logger.Errorw("error on updating history for chart deployment", "error", err, "installedAppVersion", installedAppVersion)
			return nil, err
		}
	}

	if installedAppDeploymentAction.PerformACDDeployment {
		// refresh update repo details on ArgoCD if repo is changed
		err = impl.fullModeDeploymentService.UpdateAndSyncACDApps(upgradeAppRequest, gitOpsResponse.ChartGitAttribute, monoRepoMigrationRequired, ctx, tx)
		if err != nil {
			impl.logger.Errorw("error in acd patch request", "err", err)
			return nil, err
		}
	} else if installedAppDeploymentAction.PerformHelmDeployment {
		err = impl.eaModeDeploymentService.UpgradeDeployment(upgradeAppRequest, gitOpsResponse.ChartGitAttribute, upgradeAppRequest.InstalledAppVersionHistoryId, ctx)
		if err != nil {
			if err != nil {
				impl.logger.Errorw("error in helm update request", "err", err)
				return nil, err
			}
		}
	} else if installedAppDeploymentAction.PerformFluxDeployment {
		err = impl.fullModeFluxDeploymentService.UpgradeDeployment(upgradeAppRequest, gitOpsResponse.ChartGitAttribute, upgradeAppRequest.InstalledAppVersionHistoryId, ctx)
		if err != nil {
			impl.logger.Errorw("error in flux app patch request", "err", err)
			return nil, err
		}
	}
	installedApp.UpdateStatus(appStoreBean.DEPLOY_SUCCESS)
	installedApp.UpdateAuditLog(upgradeAppRequest.UserId)
	if monoRepoMigrationRequired {
		//if mono repo case is true then repoUrl is changed then also update repo url in database
		installedApp.UpdateGitOpsRepository(gitOpsResponse.ChartGitAttribute.RepoUrl, installedApp.IsCustomRepository)
		deploymentConfig.SetRepoURL(gitOpsResponse.ChartGitAttribute.RepoUrl)
	}
	installedApp, err = impl.installedAppRepository.UpdateInstalledApp(installedApp, tx)
	if err != nil {
		impl.logger.Errorw("error in updating installed app", "err", err)
		return nil, err
	}
	upgradeAppRequest.UpdateLog(installedApp.UpdatedOn)

	deploymentConfig, err = impl.deploymentConfigService.CreateOrUpdateConfig(tx, deploymentConfig, upgradeAppRequest.UserId)
	if err != nil {
		impl.logger.Errorw("error in updating deployment config for helm apps", "appId", deploymentConfig.AppId, "envId", deploymentConfig.EnvironmentId, "err", err)
		return nil, err
	}
	return upgradeAppRequest, nil
}

func (impl *AppStoreDeploymentServiceImpl) UpdateInstalledApp(ctx context.Context, upgradeAppRequest *appStoreBean.InstallAppVersionDTO) (*appStoreBean.InstallAppVersionDTO, error) {
	triggeredAt := time.Now()
	// db operations
	dbConnection := impl.installedAppRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()
	upgradeAppRequest, err = impl.updateInstalledApp(ctx, upgradeAppRequest, tx)
	if err != nil {
		impl.logger.Errorw("error while performing updateInstalledApp", "upgradeRequest", upgradeAppRequest, "err", err)
		return nil, err
	}

	//STEP 8: finish with return response
	err = tx.Commit()
	if err != nil {
		impl.logger.Errorw("error while committing transaction to db", "error", err)
		return nil, err
	}

	if util.IsManifestDownload(upgradeAppRequest.DeploymentAppType) {
		upgradeAppRequest.HelmPackageName = adapter.GetGeneratedHelmPackageName(
			upgradeAppRequest.AppName,
			upgradeAppRequest.EnvironmentName,
			upgradeAppRequest.UpdatedOn)
		err = impl.appStoreDeploymentDBService.UpdateInstalledAppVersionHistoryStatus(
			upgradeAppRequest.InstalledAppVersionHistoryId,
			installedAppAdapter.SuccessStatusUpdateOption(upgradeAppRequest.DeploymentAppType, upgradeAppRequest.UserId),
		)
		if err != nil {
			impl.logger.Errorw("error on creating history for chart deployment", "error", err)
			return nil, err
		}
	} else if util.IsHelmApp(upgradeAppRequest.DeploymentAppType) && !impl.deploymentTypeConfig.HelmInstallASyncMode {
		err = impl.appStoreDeploymentDBService.UpdateInstalledAppVersionHistoryStatus(
			upgradeAppRequest.InstalledAppVersionHistoryId,
			installedAppAdapter.SuccessStatusUpdateOption(upgradeAppRequest.DeploymentAppType, upgradeAppRequest.UserId),
		)
		if err != nil {
			impl.logger.Errorw("error in updating install app version history on sync", "err", err)
			return nil, err
		}
	}
	err1 := impl.UpdatePreviousDeploymentStatusForAppStore(upgradeAppRequest, triggeredAt, err)
	if err1 != nil {
		impl.logger.Errorw("error while update previous installed app version history", "err", err, "installAppVersionRequest", upgradeAppRequest)
		//if installed app is updated and error is in updating previous deployment status, then don't block user, just show error.
	}
	return upgradeAppRequest, err
}

func (impl *AppStoreDeploymentServiceImpl) InstallAppByHelm(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, ctx context.Context) (*appStoreBean.InstallAppVersionDTO, error) {
	installAppVersionRequest, err := impl.eaModeDeploymentService.InstallApp(installAppVersionRequest, nil, ctx, nil)
	if err != nil {
		impl.logger.Errorw("error while installing app via helm", "error", err)
		return installAppVersionRequest, err
	}
	if util.IsHelmApp(installAppVersionRequest.DeploymentAppType) && !impl.deploymentTypeConfig.HelmInstallASyncMode {
		err = impl.appStoreDeploymentDBService.UpdateInstalledAppVersionHistoryStatus(
			installAppVersionRequest.InstalledAppVersionHistoryId,
			installedAppAdapter.SuccessStatusUpdateOption(installAppVersionRequest.DeploymentAppType, installAppVersionRequest.UserId),
		)
		if err != nil {
			impl.logger.Errorw("error in updating installed app version history with sync", "err", err)
			return installAppVersionRequest, err
		}
	}
	return installAppVersionRequest, nil
}

func (impl *AppStoreDeploymentServiceImpl) UpdatePreviousDeploymentStatusForAppStore(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, triggeredAt time.Time, err error) error {
	//creating pipeline status timeline for deployment failed
	if !util.IsAcdApp(installAppVersionRequest.DeploymentAppType) && !util.IsFluxApp(installAppVersionRequest.DeploymentAppType) {
		return nil
	}
	err1 := impl.fullModeDeploymentService.UpdateInstalledAppAndPipelineStatusForFailedDeploymentStatus(installAppVersionRequest, triggeredAt, err)
	if err1 != nil {
		impl.logger.Errorw("error in updating previous deployment status for appStore", "err", err1, "installAppVersionRequestId", installAppVersionRequest.Id)
		return err1
	}
	return nil
}

func (impl *AppStoreDeploymentServiceImpl) MarkGitOpsInstalledAppsDeletedIfArgoAppIsDeleted(installedAppId, envId int) error {
	apiError := &util.ApiError{}
	installedApp, err := impl.installedAppRepository.GetGitOpsInstalledAppsWhereArgoAppDeletedIsTrue(installedAppId, envId)
	if err != nil {
		impl.logger.Errorw("error in fetching partially deleted argoCd apps from installed app repo", "err", err)
		apiError.HttpStatusCode = http.StatusInternalServerError
		apiError.InternalMessage = "error in fetching partially deleted argoCd apps from installed app repo"
		return apiError
	}
	deploymentConfig, err := impl.deploymentConfigService.GetConfigForHelmApps(installedApp.App.Id, envId)
	if err != nil {
		impl.logger.Errorw("error in getting deployment config by appId and envId", "appId", installedAppId, "envId", envId, "err", err)
		apiError.HttpStatusCode = http.StatusInternalServerError
		apiError.InternalMessage = "error in fetching partially deleted argoCd apps from installed app repo"
		return apiError
	}
	if (!util.IsAcdApp(installedApp.DeploymentAppType) && !util.IsAcdApp(deploymentConfig.DeploymentAppType)) || !installedApp.DeploymentAppDeleteRequest {
		return nil
	}
	// Operates for ArgoCd apps only
	acdAppName := util2.BuildDeployedAppName(installedApp.App.AppName, installedApp.Environment.Name)
	isFound, err := impl.fullModeDeploymentService.CheckIfArgoAppExists(acdAppName)
	if err != nil {
		impl.logger.Errorw("error in CheckIfArgoAppExists", "err", err)
		apiError.HttpStatusCode = http.StatusInternalServerError
		apiError.InternalMessage = err.Error()
		return apiError
	}

	if isFound {
		apiError.HttpStatusCode = http.StatusInternalServerError
		apiError.InternalMessage = "App Exist in argo, error in fetching resource tree"
		return apiError
	}

	impl.logger.Warnw("app not found in argo, deleting from db ", "err", err)
	//make call to delete it from pipeline DB
	deleteRequest := &appStoreBean.InstallAppVersionDTO{}
	deleteRequest.ForceDelete = false
	deleteRequest.NonCascadeDelete = false
	deleteRequest.AcdPartialDelete = false
	deleteRequest.InstalledAppId = installedApp.Id
	deleteRequest.AppId = installedApp.AppId
	deleteRequest.AppName = installedApp.App.AppName
	deleteRequest.Namespace = installedApp.Environment.Namespace
	deleteRequest.ClusterId = installedApp.Environment.ClusterId
	deleteRequest.EnvironmentId = installedApp.EnvironmentId
	deleteRequest.AppOfferingMode = installedApp.App.AppOfferingMode
	deleteRequest.UserId = 1
	_, err = impl.DeleteInstalledApp(context.Background(), deleteRequest)
	if err != nil {
		impl.logger.Errorw("error in deleting installed app", "err", err)
		apiError.HttpStatusCode = http.StatusNotFound
		apiError.InternalMessage = "error in deleting installed app"
		return apiError
	}
	apiError.HttpStatusCode = http.StatusNotFound
	return apiError
}

func (impl *AppStoreDeploymentServiceImpl) linkHelmApplicationToChartStore(installAppVersionRequest *appStoreBean.InstallAppVersionDTO, ctx context.Context) (*openapi.UpdateReleaseResponse, error) {
	dbConnection := impl.installedAppRepository.GetConnection()
	tx, err := dbConnection.Begin()
	if err != nil {
		return nil, err
	}
	// Rollback tx on error.
	defer tx.Rollback()

	//step 1 db operation initiated
	appModel, err := impl.appRepository.FindActiveByName(installAppVersionRequest.AppName)
	if err != nil && !util.IsErrNoRows(err) {
		impl.logger.Errorw("error in getting app", "appName", installAppVersionRequest.AppName)
		return nil, err
	}
	if appModel != nil && appModel.Id > 0 {
		impl.logger.Infow("app already exists", "name", installAppVersionRequest.AppName)
		installAppVersionRequest.AppId = appModel.Id
	}
	installAppVersionRequest, err = impl.appStoreDeploymentDBService.AppStoreDeployOperationDB(installAppVersionRequest, tx, appStoreBean.INSTALL_APP_REQUEST)
	if err != nil {
		impl.logger.Errorw("error in linking chart to chart store", "err", err)
		return nil, err
	}

	// fetch app store application version from DB
	appStoreApplicationVersionId := installAppVersionRequest.AppStoreVersion
	appStoreAppVersion, err := impl.appStoreApplicationVersionRepository.FindById(appStoreApplicationVersionId)
	if err != nil {
		impl.logger.Errorw("Error in fetching app store application version", "err", err, "appStoreApplicationVersionId", appStoreApplicationVersionId)
		return nil, err
	}

	// STEP-2 update APP with chart info
	//TODO: below code is duplicated
	var IsOCIRepo bool
	var registryCredential *bean4.RegistryCredential
	var chartRepository *bean4.ChartRepository
	dockerRegistryId := appStoreAppVersion.AppStore.DockerArtifactStoreId
	if dockerRegistryId != "" {
		ociRegistryConfigs, err := impl.OCIRegistryConfigRepository.FindByDockerRegistryId(dockerRegistryId)
		if err != nil {
			impl.logger.Errorw("error in fetching oci registry config", "err", err)
			return nil, err
		}
		var ociRegistryConfig *repository2.OCIRegistryConfig
		for _, config := range ociRegistryConfigs {
			if config.RepositoryAction == repository2.STORAGE_ACTION_TYPE_PULL || config.RepositoryAction == repository2.STORAGE_ACTION_TYPE_PULL_AND_PUSH {
				ociRegistryConfig = config
				break
			}
		}
		IsOCIRepo = true
		registryCredential = &bean4.RegistryCredential{
			RegistryUrl:         appStoreAppVersion.AppStore.DockerArtifactStore.RegistryURL,
			Username:            appStoreAppVersion.AppStore.DockerArtifactStore.Username,
			Password:            appStoreAppVersion.AppStore.DockerArtifactStore.Password,
			AwsRegion:           appStoreAppVersion.AppStore.DockerArtifactStore.AWSRegion,
			AccessKey:           appStoreAppVersion.AppStore.DockerArtifactStore.AWSAccessKeyId,
			SecretKey:           appStoreAppVersion.AppStore.DockerArtifactStore.AWSSecretAccessKey,
			RegistryType:        string(appStoreAppVersion.AppStore.DockerArtifactStore.RegistryType),
			RepoName:            appStoreAppVersion.AppStore.Name,
			IsPublic:            ociRegistryConfig.IsPublic,
			Connection:          appStoreAppVersion.AppStore.DockerArtifactStore.Connection,
			RegistryName:        appStoreAppVersion.AppStore.DockerArtifactStore.Id,
			RegistryCertificate: appStoreAppVersion.AppStore.DockerArtifactStore.Cert,
		}
	} else {
		chartRepository = &bean4.ChartRepository{
			Name:                    appStoreAppVersion.AppStore.ChartRepo.Name,
			Url:                     appStoreAppVersion.AppStore.ChartRepo.Url,
			Username:                appStoreAppVersion.AppStore.ChartRepo.UserName,
			Password:                appStoreAppVersion.AppStore.ChartRepo.Password,
			AllowInsecureConnection: appStoreAppVersion.AppStore.ChartRepo.AllowInsecureConnection,
		}
	}

	updateReleaseRequest := &bean3.UpdateApplicationWithChartInfoRequestDto{
		InstallReleaseRequest: &bean4.InstallReleaseRequest{
			ValuesYaml:   installAppVersionRequest.ValuesOverrideYaml,
			ChartName:    appStoreAppVersion.Name,
			ChartVersion: appStoreAppVersion.Version,
			ReleaseIdentifier: &bean4.ReleaseIdentifier{
				ReleaseNamespace: installAppVersionRequest.Namespace,
				ReleaseName:      installAppVersionRequest.DisplayName,
			},
			RegistryCredential:         registryCredential,
			ChartRepository:            chartRepository,
			IsOCIRepo:                  IsOCIRepo,
			InstallAppVersionHistoryId: 0,
		},
		SourceAppType: bean3.SOURCE_HELM_APP,
	}

	res, err := impl.helmAppService.UpdateApplicationWithChartInfo(ctx, installAppVersionRequest.ClusterId, updateReleaseRequest)
	if err != nil {
		return nil, err
	}
	// STEP-2 ends

	// tx commit here because next operation will be process after this commit.
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	// STEP-3 install app DB post operations
	installAppVersionRequest.UpdateDeploymentAppType(util.PIPELINE_DEPLOYMENT_TYPE_HELM)
	err = impl.appStoreDeploymentDBService.InstallAppPostDbOperation(installAppVersionRequest)
	if err != nil {
		return nil, err
	}
	// STEP-3 ends

	return res, nil
}

// checkIfMonoRepoMigrationRequired checks if gitOps repo name is changed
func (impl *AppStoreDeploymentServiceImpl) checkIfMonoRepoMigrationRequired(installedApp *repository.InstalledApps, deploymentConfig *bean5.DeploymentConfig) bool {
	monoRepoMigrationRequired := false
	if !util.IsAcdApp(deploymentConfig.DeploymentAppType) || gitOps.IsGitOpsRepoNotConfigured(deploymentConfig.GetRepoURL()) || deploymentConfig.ConfigType == bean5.CUSTOM.String() {
		return false
	}
	var err error
	gitOpsRepoName := impl.gitOpsConfigReadService.GetGitOpsRepoNameFromUrl(deploymentConfig.GetRepoURL())
	if len(gitOpsRepoName) == 0 {
		gitOpsRepoName, err = impl.fullModeDeploymentService.GetAcdAppGitOpsRepoName(installedApp.App.AppName, installedApp.Environment.Name)
		if err != nil || gitOpsRepoName == "" {
			return false
		}
	}
	appNameGitOpsRepoPattern := installedApp.App.AppName + "$"
	regex := regexp.MustCompile(appNameGitOpsRepoPattern)

	// if appName is not in the gitOpsRepoName consider it as mono repo
	if !regex.MatchString(gitOpsRepoName) {
		monoRepoMigrationRequired = true
	}
	return monoRepoMigrationRequired
}

// getAppNameForInstalledApp will fetch and returns AppName from app table
func (impl *AppStoreDeploymentServiceImpl) getAppNameForInstalledApp(installedAppId int) string {
	installedApp, err := impl.installedAppRepository.GetInstalledApp(installedAppId)
	if err != nil {
		impl.logger.Errorw("UpdateProjectHelmApp, error in finding app by installedAppId", "installedAppId", installedAppId, "err", err)
		return ""
	}
	if installedApp != nil {
		return installedApp.App.AppName
	}
	return ""
}
