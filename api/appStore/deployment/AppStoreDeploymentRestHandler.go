/*
 * Copyright (c) 2020 Devtron Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package appStoreDeployment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	client "github.com/devtron-labs/devtron/api/helm-app"
	openapi "github.com/devtron-labs/devtron/api/helm-app/openapiClient"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/internal/util"
	appStoreBean "github.com/devtron-labs/devtron/pkg/appStore/bean"
	appStoreDeploymentCommon "github.com/devtron-labs/devtron/pkg/appStore/deployment/common"
	"github.com/devtron-labs/devtron/pkg/appStore/deployment/service"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/user"
	"github.com/devtron-labs/devtron/pkg/user/casbin"
	util2 "github.com/devtron-labs/devtron/util"
	"github.com/devtron-labs/devtron/util/argo"
	"github.com/devtron-labs/devtron/util/rbac"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
)

const HELM_APP_UPDATE_COUNTER = "HelmAppUpdateCounter"

type AppStoreDeploymentRestHandler interface {
	InstallApp(w http.ResponseWriter, r *http.Request)
	GetInstalledAppsByAppStoreId(w http.ResponseWriter, r *http.Request)
	DeleteInstalledApp(w http.ResponseWriter, r *http.Request)
	LinkHelmApplicationToChartStore(w http.ResponseWriter, r *http.Request)
	UpdateInstalledApp(w http.ResponseWriter, r *http.Request)
	GetInstalledAppVersion(w http.ResponseWriter, r *http.Request)
}

type AppStoreDeploymentRestHandlerImpl struct {
	Logger                     *zap.SugaredLogger
	userAuthService            user.UserService
	enforcer                   casbin.Enforcer
	enforcerUtil               rbac.EnforcerUtil
	enforcerUtilHelm           rbac.EnforcerUtilHelm
	appStoreDeploymentService  service.AppStoreDeploymentService
	appStoreDeploymentServiceC appStoreDeploymentCommon.AppStoreDeploymentCommonService
	validator                  *validator.Validate
	helmAppService             client.HelmAppService
	helmAppRestHandler         client.HelmAppRestHandler
	argoUserService            argo.ArgoUserService
	attributesService          attributes.AttributesService
}

func NewAppStoreDeploymentRestHandlerImpl(Logger *zap.SugaredLogger, userAuthService user.UserService,
	enforcer casbin.Enforcer, enforcerUtil rbac.EnforcerUtil, enforcerUtilHelm rbac.EnforcerUtilHelm, appStoreDeploymentService service.AppStoreDeploymentService,
	validator *validator.Validate, helmAppService client.HelmAppService, appStoreDeploymentServiceC appStoreDeploymentCommon.AppStoreDeploymentCommonService,
	argoUserService argo.ArgoUserService, attributesService attributes.AttributesService) *AppStoreDeploymentRestHandlerImpl {
	return &AppStoreDeploymentRestHandlerImpl{
		Logger:                     Logger,
		userAuthService:            userAuthService,
		enforcer:                   enforcer,
		enforcerUtil:               enforcerUtil,
		enforcerUtilHelm:           enforcerUtilHelm,
		appStoreDeploymentService:  appStoreDeploymentService,
		validator:                  validator,
		helmAppService:             helmAppService,
		appStoreDeploymentServiceC: appStoreDeploymentServiceC,
		argoUserService:            argoUserService,
		attributesService:          attributesService,
	}
}

func (handler AppStoreDeploymentRestHandlerImpl) InstallApp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var request appStoreBean.InstallAppVersionDTO

	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, CreateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, CreateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")

	//rbac block starts from here
	var rbacObject string
	var rbacObject2 string
	if util2.IsBaseStack() && request.EnvironmentId == 0 {

		rbacObject = handler.enforcerUtilHelm.GetHelmObjectByTeamIdAndClusterId(request.TeamId, request.ClusterId, request.Namespace, request.AppName)
		//rbacObject = handler.enforcerUtilHelm.GetHelmObjectByClusterId(request.ClusterId, request.Namespace, request.AppName)
	} else {
		rbacObject, rbacObject2 = handler.enforcerUtil.GetHelmObjectByProjectIdAndEnvId(request.TeamId, request.EnvironmentId)
	}

	var ok bool

	if rbacObject2 == "" {
		ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionCreate, rbacObject)
	} else {
		ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionCreate, rbacObject) || handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionCreate, rbacObject2)
	}

	if !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	//rbac block ends here

	isChartRepoActive, err := handler.appStoreDeploymentService.IsChartRepoActive(request.AppStoreVersion)
	if err != nil {
		handler.Logger.Errorw("service err, CreateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	if !isChartRepoActive {
		common.WriteJsonResp(w, fmt.Errorf("chart repo is disabled"), nil, http.StatusNotAcceptable)
		return
	}

	request.UserId = userId
	handler.Logger.Infow("request payload, CreateInstalledApp", "payload", request)
	ctx, cancel := context.WithCancel(r.Context())
	if cn, ok := w.(http.CloseNotifier); ok {
		go func(done <-chan struct{}, closed <-chan bool) {
			select {
			case <-done:
			case <-closed:
				cancel()
			}
		}(ctx.Done(), cn.CloseNotify())
	}
	if util2.IsBaseStack() || util2.IsHelmApp(request.AppOfferingMode) {
		ctx = context.WithValue(r.Context(), "token", token)
	} else {
		acdToken, err := handler.argoUserService.GetLatestDevtronArgoCdUserToken()
		if err != nil {
			handler.Logger.Errorw("error in getting acd token", "err", err)
			common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(r.Context(), "token", acdToken)
	}

	defer cancel()
	res, err := handler.appStoreDeploymentService.InstallApp(&request, ctx)
	if err != nil {
		if strings.Contains(err.Error(), "application spec is invalid") {
			err = &util.ApiError{Code: "400", HttpStatusCode: 400, UserMessage: "application spec is invalid, please check provided chart values"}
		}
		handler.Logger.Errorw("service err, CreateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler AppStoreDeploymentRestHandlerImpl) GetInstalledAppsByAppStoreId(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	appStoreId, err := strconv.Atoi(vars["appStoreId"])
	if err != nil {
		handler.Logger.Errorw("request err, GetInstalledAppsByAppStoreId", "err", err, "appStoreId", appStoreId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	handler.Logger.Infow("request payload, GetInstalledAppsByAppStoreId", "appStoreId", appStoreId)
	res, err := handler.appStoreDeploymentService.GetAllInstalledAppsByAppStoreId(w, r, token, appStoreId)
	if err != nil {
		handler.Logger.Errorw("service err, GetInstalledAppsByAppStoreId", "err", err, "appStoreId", appStoreId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	var installedAppsResponse []appStoreBean.InstalledAppsResponse
	for _, app := range res {

		//rbac block starts from here
		var rbacObject string
		var rbacObject2 string
		if util2.IsHelmApp(app.AppOfferingMode) {
			rbacObject = handler.enforcerUtilHelm.GetHelmObjectForEAMode(app.AppName, app.ClusterId, app.Namespace)
		} else {
			rbacObject, rbacObject2 = handler.enforcerUtil.GetHelmObjectByAppNameAndEnvId(app.AppName, app.EnvironmentId)
		}
		var ok bool
		if rbacObject2 == "" {
			ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionGet, rbacObject)
		} else {
			ok = handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionGet, rbacObject) || handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionGet, rbacObject2)
		}

		if !ok {
			continue
		}
		//rback block ends here

		installedAppsResponse = append(installedAppsResponse, app)
	}

	common.WriteJsonResp(w, err, installedAppsResponse, http.StatusOK)
}

func (handler AppStoreDeploymentRestHandlerImpl) DeleteInstalledApp(w http.ResponseWriter, r *http.Request) {

	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	installAppId, err := strconv.Atoi(vars["id"])
	if err != nil {
		handler.Logger.Errorw("request err, DeleteInstalledApp", "err", err, "installAppId", installAppId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	v := r.URL.Query()
	forceDelete := false
	force := v.Get("force")
	if len(force) > 0 {
		forceDelete, err = strconv.ParseBool(force)
		if err != nil {
			handler.Logger.Errorw("request err, DeleteInstalledApp", "err", err, "installAppId", installAppId)
			common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
			return
		}
	}
	handler.Logger.Infow("request payload, DeleteInstalledApp", "installAppId", installAppId)
	token := r.Header.Get("token")
	installedApp, err := handler.appStoreDeploymentService.GetInstalledApp(installAppId)
	if err != nil {
		handler.Logger.Error(err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	//rbac block starts from here
	var rbacObject string
	var rbacObject2 string
	if util2.IsHelmApp(installedApp.AppOfferingMode) {
		rbacObject = handler.enforcerUtilHelm.GetHelmObjectForEAMode(installedApp.AppName, installedApp.ClusterId, installedApp.Namespace)
	} else {
		rbacObject, rbacObject2 = handler.enforcerUtil.GetHelmObjectByAppNameAndEnvId(installedApp.AppName, installedApp.EnvironmentId)
	}

	var ok bool
	if rbacObject2 == "" {
		ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionDelete, rbacObject)
	} else {
		ok = handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionDelete, rbacObject) || handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionDelete, rbacObject2)
	}

	if !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	//rbac block ends here

	request := appStoreBean.InstallAppVersionDTO{}
	request.InstalledAppId = installAppId
	request.AppName = installedApp.AppName
	request.AppId = installedApp.AppId
	request.EnvironmentId = installedApp.EnvironmentId
	request.UserId = userId
	request.ForceDelete = forceDelete
	request.AppOfferingMode = installedApp.AppOfferingMode
	request.ClusterId = installedApp.ClusterId
	request.Namespace = installedApp.Namespace
	ctx, cancel := context.WithCancel(r.Context())
	if cn, ok := w.(http.CloseNotifier); ok {
		go func(done <-chan struct{}, closed <-chan bool) {
			select {
			case <-done:
			case <-closed:
				cancel()
			}
		}(ctx.Done(), cn.CloseNotify())
	}
	if util2.IsBaseStack() || util2.IsHelmApp(request.AppOfferingMode) {
		ctx = context.WithValue(r.Context(), "token", token)
	} else {
		acdToken, err := handler.argoUserService.GetLatestDevtronArgoCdUserToken()
		if err != nil {
			handler.Logger.Errorw("error in getting acd token", "err", err)
			common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(r.Context(), "token", acdToken)
	}
	res, err := handler.appStoreDeploymentService.DeleteInstalledApp(ctx, &request)
	if err != nil {
		handler.Logger.Errorw("service err, DeleteInstalledApp", "err", err, "installAppId", installAppId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler *AppStoreDeploymentRestHandlerImpl) LinkHelmApplicationToChartStore(w http.ResponseWriter, r *http.Request) {
	request := &openapi.UpdateReleaseWithChartLinkingRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	appIdentifier, err := handler.helmAppService.DecodeAppId(*request.AppId)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}

	// RBAC enforcer applying
	rbacObject := handler.enforcerUtilHelm.GetHelmObjectByClusterId(appIdentifier.ClusterId, appIdentifier.Namespace, appIdentifier.ReleaseName)
	token := r.Header.Get("token")
	if ok := handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionUpdate, rbacObject); !ok {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	//RBAC enforcer Ends

	res, isChartRepoActive, err := handler.appStoreDeploymentService.LinkHelmApplicationToChartStore(context.Background(), request, appIdentifier, userId)
	if err != nil {
		handler.Logger.Errorw("Error in UpdateApplicationWithChartStoreLinking", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	} else if !isChartRepoActive {
		common.WriteJsonResp(w, fmt.Errorf("chart repo is disabled"), nil, http.StatusNotAcceptable)
		return
	}

	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler AppStoreDeploymentRestHandlerImpl) UpdateInstalledApp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	var request appStoreBean.InstallAppVersionDTO
	err = decoder.Decode(&request)
	if err != nil {
		handler.Logger.Errorw("request err, UpdateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	err = handler.validator.Struct(request)
	if err != nil {
		handler.Logger.Errorw("validation err, UpdateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	handler.Logger.Debugw("request payload, UpdateInstalledApp", "payload", request)
	installedApp, err := handler.appStoreDeploymentService.GetInstalledApp(request.InstalledAppId)
	if err != nil {
		handler.Logger.Errorw("service err, UpdateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	//rbac block starts from here
	var rbacObject string
	var rbacObject2 string
	if util2.IsHelmApp(installedApp.AppOfferingMode) {
		rbacObject = handler.enforcerUtilHelm.GetHelmObjectForEAMode(installedApp.AppName, installedApp.ClusterId, installedApp.Namespace)
	} else {
		rbacObject, rbacObject2 = handler.enforcerUtil.GetHelmObject(installedApp.AppId, installedApp.EnvironmentId)
	}

	var ok bool

	if rbacObject2 == "" {
		ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionUpdate, rbacObject)
	} else {
		ok = handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionUpdate, rbacObject) || handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionUpdate, rbacObject2)
	}

	if !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	//rbac block ends here

	request.UserId = userId
	ctx, cancel := context.WithCancel(r.Context())
	if cn, ok := w.(http.CloseNotifier); ok {
		go func(done <-chan struct{}, closed <-chan bool) {
			select {
			case <-done:
			case <-closed:
				cancel()
			}
		}(ctx.Done(), cn.CloseNotify())
	}
	if util2.IsBaseStack() || util2.IsHelmApp(request.AppOfferingMode) {
		ctx = context.WithValue(r.Context(), "token", token)
	} else {
		acdToken, err := handler.argoUserService.GetLatestDevtronArgoCdUserToken()
		if err != nil {
			handler.Logger.Errorw("error in getting acd token", "err", err)
			common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
			return
		}
		ctx = context.WithValue(r.Context(), "token", acdToken)
	}
	res, err := handler.appStoreDeploymentService.UpdateInstalledApp(ctx, &request)
	if err != nil {
		if strings.Contains(err.Error(), "application spec is invalid") {
			err = &util.ApiError{Code: "400", HttpStatusCode: 400, UserMessage: "application spec is invalid, please check provided chart values"}
		}
		handler.Logger.Errorw("service err, UpdateInstalledApp", "err", err, "payload", request)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	err = handler.attributesService.UpdateKeyValueByOne(HELM_APP_UPDATE_COUNTER)

	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler AppStoreDeploymentRestHandlerImpl) GetInstalledAppVersion(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	installedAppId, err := strconv.Atoi(vars["installedAppVersionId"])
	if err != nil {
		handler.Logger.Errorw("request err, GetInstalledAppVersion", "err", err, "installedAppVersionId", installedAppId)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	handler.Logger.Infow("request payload, GetInstalledAppVersion", "installedAppVersionId", installedAppId)
	dto, err := handler.appStoreDeploymentService.GetInstalledAppVersion(installedAppId, userId)
	if err != nil {
		handler.Logger.Errorw("service err, GetInstalledAppVersion", "err", err, "installedAppVersionId", installedAppId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	//rbac block starts from here
	var rbacObject string
	var rbacObject2 string
	if util2.IsHelmApp(dto.AppOfferingMode) {
		rbacObject = handler.enforcerUtilHelm.GetHelmObjectForEAMode(dto.AppName, dto.ClusterId, dto.Namespace)
	} else {
		rbacObject, rbacObject2 = handler.enforcerUtil.GetHelmObjectByAppNameAndEnvId(dto.AppName, dto.EnvironmentId)
	}
	var ok bool
	if rbacObject2 == "" {
		ok = handler.enforcer.EnforceByEmail("ayush@devtron.ai", casbin.ResourceHelmApp, casbin.ActionGet, rbacObject)
	} else {
		ok = handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionGet, rbacObject) || handler.enforcer.Enforce(token, casbin.ResourceHelmApp, casbin.ActionGet, rbacObject2)
	}

	if !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	//rbac block ends here

	common.WriteJsonResp(w, err, dto, http.StatusOK)
}
