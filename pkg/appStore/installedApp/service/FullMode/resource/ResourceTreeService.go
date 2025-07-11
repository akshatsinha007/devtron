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

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	k8s2 "github.com/devtron-labs/common-lib/utils/k8s"
	"github.com/devtron-labs/common-lib/utils/k8s/commonBean"
	"github.com/devtron-labs/common-lib/utils/k8sObjectsUtil"
	"github.com/devtron-labs/devtron/api/bean/AppView"
	"github.com/devtron-labs/devtron/api/helm-app/gRPC"
	"github.com/devtron-labs/devtron/api/helm-app/service"
	"github.com/devtron-labs/devtron/api/helm-app/service/bean"
	"github.com/devtron-labs/devtron/api/helm-app/service/read"
	"github.com/devtron-labs/devtron/client/argocdServer"
	"github.com/devtron-labs/devtron/internal/constants"
	repository2 "github.com/devtron-labs/devtron/internal/sql/repository/dockerRegistry"
	cdWorkflow2 "github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig/bean/workflow/cdWorkflow"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/appStatus"
	"github.com/devtron-labs/devtron/pkg/appStore/bean"
	appStoreDiscoverRepository "github.com/devtron-labs/devtron/pkg/appStore/discover/repository"
	"github.com/devtron-labs/devtron/pkg/appStore/installedApp/repository"
	"github.com/devtron-labs/devtron/pkg/argoApplication"
	bean3 "github.com/devtron-labs/devtron/pkg/argoApplication/bean"
	bean4 "github.com/devtron-labs/devtron/pkg/auth/user/bean"
	"github.com/devtron-labs/devtron/pkg/deployment/common"
	bean2 "github.com/devtron-labs/devtron/pkg/deployment/common/bean"
	"github.com/devtron-labs/devtron/pkg/fluxApplication"
	bean5 "github.com/devtron-labs/devtron/pkg/fluxApplication/bean"
	"github.com/devtron-labs/devtron/pkg/k8s"
	application3 "github.com/devtron-labs/devtron/pkg/k8s/application"
	util3 "github.com/devtron-labs/devtron/pkg/util"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type InstalledAppResourceService interface {
	FetchResourceTreeWithHibernateForACD(rctx context.Context, cn http.CloseNotifier, appDetail *AppView.AppDetailContainer) AppView.AppDetailContainer
	FetchResourceTree(rctx context.Context, cn http.CloseNotifier, appDetailsContainer *AppView.AppDetailsContainer, installedApp repository.InstalledApps, deploymentConfig *bean2.DeploymentConfig, helmReleaseInstallStatus string, status string) error
	FetchChartNotes(installedAppId int, envId int, token string, checkNotesAuth func(token string, appName string, envId int) bool) (string, error)
}

type InstalledAppResourceServiceImpl struct {
	logger                               *zap.SugaredLogger
	installedAppRepository               repository.InstalledAppRepository
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository
	acdClientWrapper                     argocdServer.ArgoClientWrapperService
	aCDAuthConfig                        *util3.ACDAuthConfig
	installedAppRepositoryHistory        repository.InstalledAppVersionHistoryRepository
	helmAppService                       service.HelmAppService
	helmAppReadService                   read.HelmAppReadService
	appStatusService                     appStatus.AppStatusService
	k8sCommonService                     k8s.K8sCommonService
	k8sApplicationService                application3.K8sApplicationService
	K8sUtil                              k8s2.K8sService
	deploymentConfigurationService       common.DeploymentConfigService
	OCIRegistryConfigRepository          repository2.OCIRegistryConfigRepository
	argoApplicationService               argoApplication.ArgoApplicationService
	fluxApplicationService               fluxApplication.FluxApplicationService
}

func NewInstalledAppResourceServiceImpl(logger *zap.SugaredLogger,
	installedAppRepository repository.InstalledAppRepository,
	appStoreApplicationVersionRepository appStoreDiscoverRepository.AppStoreApplicationVersionRepository,
	acdClientWrapper argocdServer.ArgoClientWrapperService,
	aCDAuthConfig *util3.ACDAuthConfig,
	installedAppRepositoryHistory repository.InstalledAppVersionHistoryRepository,
	helmAppService service.HelmAppService,
	helmAppReadService read.HelmAppReadService,
	appStatusService appStatus.AppStatusService,
	k8sCommonService k8s.K8sCommonService, k8sApplicationService application3.K8sApplicationService, K8sUtil k8s2.K8sService,
	deploymentConfigurationService common.DeploymentConfigService,
	OCIRegistryConfigRepository repository2.OCIRegistryConfigRepository,
	argoApplicationService argoApplication.ArgoApplicationService,
	fluxApplicationService fluxApplication.FluxApplicationService) *InstalledAppResourceServiceImpl {
	return &InstalledAppResourceServiceImpl{
		logger:                               logger,
		installedAppRepository:               installedAppRepository,
		appStoreApplicationVersionRepository: appStoreApplicationVersionRepository,
		acdClientWrapper:                     acdClientWrapper,
		aCDAuthConfig:                        aCDAuthConfig,
		installedAppRepositoryHistory:        installedAppRepositoryHistory,
		helmAppService:                       helmAppService,
		helmAppReadService:                   helmAppReadService,
		appStatusService:                     appStatusService,
		k8sCommonService:                     k8sCommonService,
		k8sApplicationService:                k8sApplicationService,
		K8sUtil:                              K8sUtil,
		deploymentConfigurationService:       deploymentConfigurationService,
		OCIRegistryConfigRepository:          OCIRegistryConfigRepository,
		argoApplicationService:               argoApplicationService,
		fluxApplicationService:               fluxApplicationService,
	}
}

func (impl *InstalledAppResourceServiceImpl) FetchResourceTree(rctx context.Context, cn http.CloseNotifier, appDetailsContainer *AppView.AppDetailsContainer, installedApp repository.InstalledApps, deploymentConfig *bean2.DeploymentConfig, helmReleaseInstallStatus string, status string) error {
	var err error
	var resourceTree map[string]interface{}
	deploymentAppName := util2.BuildDeployedAppName(installedApp.App.AppName, installedApp.Environment.Name)
	if util.IsAcdApp(deploymentConfig.DeploymentAppType) {
		resourceTree, err = impl.fetchResourceTreeForACD(rctx, cn, installedApp.App.Id, installedApp.EnvironmentId, installedApp.Environment.ClusterId, deploymentAppName, installedApp.Environment.Namespace)
	} else if util.IsHelmApp(deploymentConfig.DeploymentAppType) {
		req := &bean.AppIdentifier{
			ClusterId:   installedApp.Environment.ClusterId,
			Namespace:   installedApp.Environment.Namespace,
			ReleaseName: installedApp.App.AppName,
		}
		detail, err := impl.helmAppService.GetApplicationDetail(rctx, req)
		if err != nil {
			impl.logger.Errorw("error in fetching app detail", "err", err)
		}

		/* helmReleaseInstallStatus is nats message sent from kubelink to orchestrator and has the following details about installation :-
		1) isReleaseInstalled -> whether release object is created or not in this installation
		2) ErrorInInstallation -> if there is error in installation
		3) Message -> error message/ success message
		4) InstallAppVersionHistoryId
		5) Status -> Progressing, Failed, Succeeded
		*/

		if detail != nil && detail.ReleaseExist {

			resourceTree = util2.InterfaceToMapAdapter(detail.ResourceTreeResponse)
			resourceTree["status"] = detail.ApplicationStatus
			appDetailsContainer.Notes = detail.ChartMetadata.Notes

			helmInstallStatus := &appStoreBean.HelmReleaseStatusConfig{}
			releaseStatus := detail.ReleaseStatus

			if len(helmReleaseInstallStatus) > 0 {
				err := json.Unmarshal([]byte(helmReleaseInstallStatus), helmInstallStatus)
				if err != nil {
					impl.logger.Errorw("error in unmarshalling helm release install status")
					return err
				}
				// ReleaseExist=true in app detail container but helm install status says that isReleaseInstalled=false which means this release was created externally
				if helmInstallStatus.IsReleaseInstalled == false && status != "Progressing" {
					/*
						Handling case when :-
						1) An external release with name "foo" exist
						2) User creates an app with same name i.e "foo"
						3) In this case we use helmReleaseInstallStatus which will have status of our release and not external release
					*/
					resourceTree = make(map[string]interface{})
					releaseStatus = impl.getReleaseStatusFromHelmReleaseInstallStatus(helmReleaseInstallStatus, status)
				}
			}
			releaseStatusMap := util2.InterfaceToMapAdapter(releaseStatus)
			appDetailsContainer.ReleaseStatus = releaseStatusMap
		} else {
			// case when helm release is not created
			releaseStatus := impl.getReleaseStatusFromHelmReleaseInstallStatus(helmReleaseInstallStatus, status)
			releaseStatusMap := util2.InterfaceToMapAdapter(releaseStatus)
			appDetailsContainer.ReleaseStatus = releaseStatusMap
		}
	} else if util.IsFluxApp(deploymentConfig.DeploymentAppType) {
		req := &bean5.FluxAppIdentifier{
			ClusterId:      installedApp.Environment.ClusterId,
			Namespace:      installedApp.Environment.Namespace,
			Name:           deploymentAppName,
			IsKustomizeApp: false,
		}
		detail, err := impl.fluxApplicationService.GetFluxAppDetail(rctx, req)
		if err != nil {
			impl.logger.Errorw("error in fetching app detail", "payload", req, "err", err)
		}
		if detail != nil {
			iaVersion, err := impl.installedAppRepositoryHistory.GetLatestInstalledAppVersionHistoryByInstalledAppId(installedApp.Id)
			if err != nil {
				impl.logger.Errorw("error in fetching latest installed app version history", "err", err)
			} else if iaVersion != nil {
				switch detail.FluxAppStatusDetail.Reason {
				case bean5.InstallSucceededReason, bean5.UpgradeSucceededReason, bean5.TestSucceededReason, bean5.RollbackSucceededReason:
					if detail.AppHealthStatus == commonBean.HealthStatusHealthy {
						iaVersion.Status = cdWorkflow2.WorkflowSucceeded
					}
				case bean5.UpgradeFailedReason,
					bean5.TestFailedReason,
					bean5.RollbackFailedReason,
					bean5.UninstallFailedReason,
					bean5.ArtifactFailedReason,
					bean5.InstallFailedReason:
					iaVersion.Status = cdWorkflow2.WorkflowFailed
				}
				iaVersion.UpdateAuditLog(bean4.SYSTEM_USER_ID)
				_, err = impl.installedAppRepositoryHistory.UpdateInstalledAppVersionHistory(iaVersion, nil)
				if err != nil {
					impl.logger.Errorw("error in updating app version history", "err", err)
				}
			}
			resourceTree = util2.InterfaceToMapAdapter(detail.ResourceTreeResponse)
			applicationStatus := detail.AppHealthStatus
			if detail.ResourceTreeResponse != nil && len(detail.ResourceTreeResponse.Nodes) == 0 {
				resourceTree["status"] = commonBean.HealthStatusUnknown
			} else {
				resourceTree["status"] = applicationStatus
			}
		}
	}
	if resourceTree != nil {
		version, err := impl.k8sCommonService.GetK8sServerVersion(installedApp.Environment.ClusterId)
		if err != nil {
			impl.logger.Errorw("error in fetching k8s version in resource tree call fetching", "clusterId", installedApp.Environment.ClusterId, "err", err)
		} else {
			resourceTree["serverVersion"] = version.String()
		}
		appDetailsContainer.ResourceTree = resourceTree
	}
	return err
}

func (impl *InstalledAppResourceServiceImpl) FetchResourceTreeWithHibernateForACD(rctx context.Context, cn http.CloseNotifier, appDetail *AppView.AppDetailContainer) AppView.AppDetailContainer {
	ctx, cancel := context.WithCancel(rctx)
	if cn != nil {
		go func(done <-chan struct{}, closed <-chan bool) {
			select {
			case <-done:
			case <-closed:
				cancel()
			}
		}(ctx.Done(), cn.CloseNotify())
	}
	defer cancel()
	deploymentAppName := util2.BuildDeployedAppName(appDetail.AppName, appDetail.EnvironmentName)
	resourceTree, err := impl.fetchResourceTreeForACD(rctx, cn, appDetail.AppId, appDetail.EnvironmentId, appDetail.ClusterId, deploymentAppName, appDetail.Namespace)
	appDetail.ResourceTree = resourceTree
	if err != nil {
		return *appDetail
	}
	if appDetail.ResourceTree["nodes"] == nil {
		return *appDetail
	}
	appDetail.ResourceTree, _ = impl.checkHibernate(appDetail.ResourceTree, deploymentAppName, ctx)
	return *appDetail
}

func (impl *InstalledAppResourceServiceImpl) fetchResourceTreeForACD(rctx context.Context, cn http.CloseNotifier, appId int, envId, clusterId int, deploymentAppName, namespace string) (map[string]interface{}, error) {
	var resourceTree map[string]interface{}
	query := &application.ResourcesQuery{
		ApplicationName: &deploymentAppName,
	}
	ctx, cancel := context.WithCancel(rctx)
	if cn != nil {
		go func(done <-chan struct{}, closed <-chan bool) {
			select {
			case <-done:
			case <-closed:
				cancel()
			}
		}(ctx.Done(), cn.CloseNotify())
	}
	defer cancel()
	start := time.Now()
	resp, err := impl.argoApplicationService.GetResourceTree(ctx, bean3.NewImperativeQueryRequest(query))
	elapsed := time.Since(start)
	impl.logger.Debugf("Time elapsed %s in fetching app-store installed application %s for environment %s", elapsed, deploymentAppName, envId)
	if err != nil {
		impl.logger.Errorw("service err, FetchAppDetailsForInstalledAppV2, fetching resource tree", "err", err, "installedAppId", appId, "envId", envId)
		err = &util.ApiError{
			Code:            constants.AppDetailResourceTreeNotFound,
			InternalMessage: "app detail fetched, failed to get resource tree from acd",
			UserMessage:     "app detail fetched, failed to get resource tree from acd",
		}
		return resourceTree, err
	}
	label := fmt.Sprintf("app.kubernetes.io/instance=%s", deploymentAppName)
	pods, err := impl.k8sApplicationService.GetPodListByLabel(clusterId, namespace, label)
	if err != nil {
		impl.logger.Errorw("error in getting pods by label", "err", err, "clusterId", clusterId, "namespace", namespace, "label", label)
		return resourceTree, err
	}
	ephemeralContainersMap := k8sObjectsUtil.ExtractEphemeralContainers(pods)
	for _, metaData := range resp.PodMetadata {
		metaData.EphemeralContainers = ephemeralContainersMap[metaData.Name]
	}
	resourceTree = util2.InterfaceToMapAdapter(resp)
	resourceTree, hibernationStatus := impl.checkHibernate(resourceTree, deploymentAppName, ctx)
	appStatus := resp.Status
	if resourceTree != nil {
		if hibernationStatus != "" {
			resourceTree["status"] = hibernationStatus
			appStatus = hibernationStatus
		}
	}
	// using this resp.Status to update in app_status table
	//FIXME: remove this dangling thread
	go func() {
		err = impl.appStatusService.UpdateStatusWithAppIdEnvId(appId, envId, appStatus)
		if err != nil {
			impl.logger.Warnw("error in updating app status", "err", err, appId, "envId", envId)
		}
	}()
	impl.logger.Debugf("application %s in environment %s had status %+v\n", appId, envId, resp)
	k8sAppDetail := AppView.AppDetailContainer{
		DeploymentDetailContainer: AppView.DeploymentDetailContainer{
			ClusterId: clusterId,
			Namespace: namespace,
		},
	}
	clusterIdString := strconv.Itoa(clusterId)
	validRequest := impl.k8sCommonService.FilterK8sResources(rctx, resourceTree, k8sAppDetail, clusterIdString, []string{commonBean.ServiceKind, commonBean.EndpointsKind, commonBean.IngressKind}, "")
	response, err := impl.k8sCommonService.GetManifestsByBatch(rctx, validRequest)
	if err != nil {
		impl.logger.Errorw("error in getting manifest by batch", "clusterId", clusterIdString, "err", err)
		return nil, err
	}
	newResourceTree := impl.k8sCommonService.PortNumberExtraction(response, resourceTree)
	return newResourceTree, err
}

func (impl *InstalledAppResourceServiceImpl) checkHibernate(resp map[string]interface{}, deploymentAppName string, ctx context.Context) (map[string]interface{}, string) {

	if resp == nil {
		return resp, ""
	}
	responseTree := resp
	var canBeHibernated uint64 = 0
	var hibernated uint64 = 0
	responseTreeNodes, ok := responseTree["nodes"]
	if !ok {
		return resp, ""
	}
	replicaNodes := impl.filterOutReplicaNodes(responseTreeNodes)
	batchSize := impl.aCDAuthConfig.ResourceListForReplicasBatchSize
	requestsLength := len(replicaNodes)
	for i := 0; i < requestsLength; {
		//requests left to process
		remainingBatch := requestsLength - i
		if remainingBatch < batchSize {
			batchSize = remainingBatch
		}
		var wg sync.WaitGroup
		for j := 0; j < batchSize; j++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				canBeHibernatedFlag, hibernatedFlag := impl.processReplicaNodeForHibernation(replicaNodes[i+j], deploymentAppName, ctx)
				if canBeHibernatedFlag {
					atomic.AddUint64(&canBeHibernated, 1)
				}
				if hibernatedFlag {
					atomic.AddUint64(&hibernated, 1)
				}
			}(j)
		}
		wg.Wait()
		i += batchSize
	}

	status := ""
	if hibernated > 0 && canBeHibernated > 0 {
		if hibernated == canBeHibernated {
			status = appStatus.HealthStatusHibernating
		} else if hibernated < canBeHibernated {
			status = appStatus.HealthStatusPartiallyHibernated
		}
	}

	return responseTree, status
}

func (impl *InstalledAppResourceServiceImpl) processReplicaNodeForHibernation(node interface{}, deploymentAppName string, ctx context.Context) (bool, bool) {
	currNode := node.(interface{}).(map[string]interface{})
	resName := util2.InterfaceToString(currNode["name"])
	resKind := util2.InterfaceToString(currNode["kind"])
	resGroup := util2.InterfaceToString(currNode["group"])
	resVersion := util2.InterfaceToString(currNode["version"])
	resNamespace := util2.InterfaceToString(currNode["namespace"])
	rQuery := &application.ApplicationResourceRequest{
		Name:         &deploymentAppName,
		ResourceName: &resName,
		Kind:         &resKind,
		Group:        &resGroup,
		Version:      &resVersion,
		Namespace:    &resNamespace,
	}
	canBeHibernatedFlag := false
	alreadyHibernated := false

	if currNode["parentRefs"] == nil {
		canBeHibernatedFlag, alreadyHibernated = impl.checkForHibernation(ctx, rQuery, currNode)
	}
	return canBeHibernatedFlag, alreadyHibernated
}

func (impl *InstalledAppResourceServiceImpl) checkForHibernation(ctx context.Context, rQuery *application.ApplicationResourceRequest, currNode map[string]interface{}) (bool, bool) {
	t0 := time.Now()
	canBeHibernated := false
	alreadyHibernated := false
	ctx, _ = context.WithTimeout(ctx, 60*time.Second)
	res, err := impl.acdClientWrapper.GetApplicationResource(ctx, rQuery)
	if err != nil {
		impl.logger.Errorw("error getting response from acdClient", "request", rQuery, "data", res, "timeTaken", time.Since(t0), "err", err)
		return canBeHibernated, alreadyHibernated
	}
	if res.Manifest != nil {
		manifest, _ := gjson.Parse(*res.Manifest).Value().(map[string]interface{})
		replicas := util2.InterfaceToMapAdapter(manifest["spec"])["replicas"]
		if replicas != nil {
			currNode["canBeHibernated"] = true
			canBeHibernated = true
		}
		annotations := util2.InterfaceToMapAdapter(manifest["metadata"])["annotations"]
		if annotations != nil {
			val := util2.InterfaceToMapAdapter(annotations)["hibernator.devtron.ai/replicas"]
			if val != nil {
				if util2.InterfaceToString(val) != "0" && util2.InterfaceToFloat(replicas) == 0 {
					currNode["isHibernated"] = true
					alreadyHibernated = true
				}
			}
		}
	}
	return canBeHibernated, alreadyHibernated
}

func (impl *InstalledAppResourceServiceImpl) filterOutReplicaNodes(responseTreeNodes interface{}) []interface{} {
	resourceListForReplicas := impl.aCDAuthConfig.ResourceListForReplicas
	entries := strings.Split(resourceListForReplicas, ",")
	resourceListMap := util2.ConvertStringSliceToMap(entries)
	var replicaNodes []interface{}
	for _, node := range responseTreeNodes.(interface{}).([]interface{}) {
		currNode := node.(interface{}).(map[string]interface{})
		resKind := util2.InterfaceToString(currNode["kind"])
		if _, ok := resourceListMap[resKind]; ok {
			replicaNodes = append(replicaNodes, node)
		}
	}
	return replicaNodes
}

func (impl *InstalledAppResourceServiceImpl) getReleaseStatusFromHelmReleaseInstallStatus(helmReleaseInstallStatus string, status string) *gRPC.ReleaseStatus {
	//release status is sent in resource tree call and is shown on UI as helm config apply status
	releaseStatus := &gRPC.ReleaseStatus{}
	if len(helmReleaseInstallStatus) > 0 {
		helmInstallStatus := &appStoreBean.HelmReleaseStatusConfig{}
		err := json.Unmarshal([]byte(helmReleaseInstallStatus), helmInstallStatus)
		if err != nil {
			impl.logger.Errorw("error in unmarshalling helm release install status")
			return releaseStatus
		}
		if status == appStoreBean.HELM_RELEASE_STATUS_FAILED {
			releaseStatus.Status = status
			releaseStatus.Description = helmInstallStatus.Message
			releaseStatus.Message = "Release install/upgrade failed"
		} else if status == appStoreBean.HELM_RELEASE_STATUS_PROGRESSING {
			releaseStatus.Status = status
			releaseStatus.Description = helmInstallStatus.Message
			releaseStatus.Message = helmInstallStatus.Message
		} else {
			// there can be a case when helm release is created but we are not able to fetch it
			releaseStatus.Status = appStoreBean.HELM_RELEASE_STATUS_UNKNOWN
			releaseStatus.Description = "Unable to fetch release for app"
			releaseStatus.Message = "Unable to fetch release for app"
		}
	} else {
		releaseStatus.Status = appStoreBean.HELM_RELEASE_STATUS_UNKNOWN
		releaseStatus.Description = "Release not found"
		releaseStatus.Message = "Release not found "
	}
	return releaseStatus
}
