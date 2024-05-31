/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package appStoreDiscover

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/devtron-labs/devtron/api/restHandler/common"
	appStoreBean "github.com/devtron-labs/devtron/pkg/appStore/bean"
	"github.com/devtron-labs/devtron/pkg/appStore/discover/service"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AppStoreRestHandler interface {
	FindAllApps(w http.ResponseWriter, r *http.Request)
	GetChartDetailsForVersion(w http.ResponseWriter, r *http.Request)
	GetChartVersions(w http.ResponseWriter, r *http.Request)
	GetChartInfo(w http.ResponseWriter, r *http.Request)
	SearchAppStoreChartByName(w http.ResponseWriter, r *http.Request)
}

type AppStoreRestHandlerImpl struct {
	Logger          *zap.SugaredLogger
	appStoreService service.AppStoreService
	userAuthService user.UserService
	enforcer        casbin.Enforcer
}

func NewAppStoreRestHandlerImpl(Logger *zap.SugaredLogger, userAuthService user.UserService, appStoreService service.AppStoreService,
	enforcer casbin.Enforcer) *AppStoreRestHandlerImpl {
	return &AppStoreRestHandlerImpl{
		Logger:          Logger,
		appStoreService: appStoreService,
		userAuthService: userAuthService,
		enforcer:        enforcer,
	}
}

func (handler *AppStoreRestHandlerImpl) FindAllApps(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusUnauthorized)
		return
	}

	v := r.URL.Query()
	deprecated := false
	deprecatedStr := v.Get("includeDeprecated")
	if len(deprecatedStr) > 0 {
		deprecated, err = strconv.ParseBool(deprecatedStr)
		if err != nil {
			deprecated = false
		}
	}

	var chartRepoIds []int
	chartRepoIdsStr := v.Get("chartRepoId")
	if len(chartRepoIdsStr) > 0 {
		chartRepoIdStrArr := strings.Split(chartRepoIdsStr, ",")
		for _, chartRepoIdStr := range chartRepoIdStrArr {
			chartRepoId, err := strconv.Atoi(chartRepoIdStr)
			if err == nil {
				chartRepoIds = append(chartRepoIds, chartRepoId)
			}
		}
	}
	var registryIds []string
	registryIdsStrArr := v.Get("registryId")
	if len(registryIdsStrArr) > 0 {
		registryIdStrArr := strings.Split(registryIdsStrArr, ",")
		for _, registryId := range registryIdStrArr {
			registryIds = append(registryIds, registryId)
		}
	}

	appStoreName := strings.ToLower(v.Get("appStoreName"))

	offset := 0
	offsetStr := v.Get("offset")
	if len(offsetStr) > 0 {
		offset, _ = strconv.Atoi(offsetStr)
	}
	size := 0
	sizeStr := v.Get("size")
	if len(sizeStr) > 0 {
		size, _ = strconv.Atoi(sizeStr)
	}
	filter := &appStoreBean.AppStoreFilter{
		IncludeDeprecated: deprecated,
		ChartRepoId:       chartRepoIds,
		RegistryId:        registryIds,
		AppStoreName:      appStoreName,
	}
	if size > 0 {
		filter.Size = size
		filter.Offset = offset
	}
	handler.Logger.Infow("request payload, FindAllApps, app store", "userId", userId)
	res, err := handler.appStoreService.FindAllApps(filter)
	if err != nil {
		handler.Logger.Errorw("service err, FindAllApps, app store", "err", err, "userId", userId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler *AppStoreRestHandlerImpl) GetChartDetailsForVersion(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		handler.Logger.Errorw("request err, GetChartDetailsForVersion", "err", err, "appStoreApplicationVersionId", id)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	handler.Logger.Infow("request payload, GetChartDetailsForVersion, app store", "appStoreApplicationVersionId", id)
	res, err := handler.appStoreService.FindChartDetailsById(id)
	if err != nil {
		handler.Logger.Errorw("service err, GetChartDetailsForVersion, app store", "err", err, "appStoreApplicationVersionId", id)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler *AppStoreRestHandlerImpl) GetChartVersions(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["appStoreId"])
	if err != nil {
		handler.Logger.Errorw("request err, GetChartVersions", "err", err, "appStoreId", id)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	handler.Logger.Infow("request payload, GetChartVersions, app store", "appStoreId", id)
	res, err := handler.appStoreService.FindChartVersionsByAppStoreId(id)
	if err != nil {
		handler.Logger.Errorw("service err, GetChartVersions, app store", "err", err, "appStoreId", id)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler *AppStoreRestHandlerImpl) GetChartInfo(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["appStoreApplicationVersionId"])
	if err != nil {
		handler.Logger.Errorw("request err, GetChartInfo", "err", err, "appStoreApplicationVersionId", id)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	handler.Logger.Infow("request payload, GetChartInfo, app store", "appStoreApplicationVersionId", id)
	res, err := handler.appStoreService.GetChartInfoByAppStoreApplicationVersionId(id)
	if err != nil {
		handler.Logger.Errorw("service err, GetChartInfo, fetching resource tree", "err", err, "appStoreApplicationVersionId", id)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}

func (handler *AppStoreRestHandlerImpl) SearchAppStoreChartByName(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userAuthService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	chartName := vars["chartName"]
	handler.Logger.Infow("request payload, SearchAppStoreChartByName, app store", "chartName", chartName)
	res, err := handler.appStoreService.SearchAppStoreChartByName(chartName)
	if err != nil {
		handler.Logger.Errorw("service err, SearchAppStoreChartByName, app store", "err", err, "userId", userId)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, res, http.StatusOK)
}
