package devtronResource

import (
	"fmt"
	apiBean "github.com/devtron-labs/devtron/api/devtronResource/bean"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	"github.com/devtron-labs/devtron/pkg/devtronResource/history/deployment/cdPipeline"
	"github.com/gorilla/schema"
	"go.uber.org/zap"
	"net/http"
)

type HistoryRestHandler interface {
	GetDeploymentHistory(w http.ResponseWriter, r *http.Request)
	GetDeploymentHistoryConfigList(w http.ResponseWriter, r *http.Request)
}

type HistoryRestHandlerImpl struct {
	logger                   *zap.SugaredLogger
	enforcer                 casbin.Enforcer
	deploymentHistoryService cdPipeline.DeploymentHistoryService
}

func NewHistoryRestHandlerImpl(logger *zap.SugaredLogger,
	enforcer casbin.Enforcer,
	deploymentHistoryService cdPipeline.DeploymentHistoryService) *HistoryRestHandlerImpl {
	return &HistoryRestHandlerImpl{
		logger:                   logger,
		enforcer:                 enforcer,
		deploymentHistoryService: deploymentHistoryService,
	}
}

func (handler *HistoryRestHandlerImpl) GetDeploymentHistory(w http.ResponseWriter, r *http.Request) {
	//kind, subKind, versionVar, caughtError := getKindSubKindVersion(w, r)
	//if caughtError {
	//	return
	//}
	v := r.URL.Query()
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	queryParams := apiBean.GetHistoryQueryParams{}
	err := decoder.Decode(&queryParams, v)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	isValidated := handler.enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionCreate, "*")
	if !isValidated {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	resp, err := handler.deploymentHistoryService.GetCdPipelineDeploymentHistory(queryParams.OffSet, queryParams.Limit, queryParams.FilterCriteria)
	if err != nil {
		handler.logger.Errorw("service error, GetCdPipelineDeploymentHistory", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, resp, http.StatusOK)
	return
}

func (handler *HistoryRestHandlerImpl) GetDeploymentHistoryConfigList(w http.ResponseWriter, r *http.Request) {
	//kind, subKind, versionVar, caughtError := getKindSubKindVersion(w, r)
	//if caughtError {
	//	return
	//}
	v := r.URL.Query()
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	queryParams := apiBean.GetHistoryConfigQueryParams{}
	err := decoder.Decode(&queryParams, v)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}
	token := r.Header.Get("token")
	isValidated := handler.enforcer.Enforce(token, casbin.ResourceGlobal, casbin.ActionCreate, "*")
	if !isValidated {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), nil, http.StatusForbidden)
		return
	}
	resp, err := handler.deploymentHistoryService.GetCdPipelineDeploymentHistoryConfigList(queryParams.BaseConfigurationId, queryParams.HistoryComponent,
		queryParams.HistoryComponentName, queryParams.FilterCriteria)
	if err != nil {
		handler.logger.Errorw("service error, GetCdPipelineDeploymentHistoryConfigList", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, err, resp, http.StatusOK)
	return
}
