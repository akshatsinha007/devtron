package artifactPromotionApprovalRequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/devtron-labs/devtron/api/restHandler/common"
	repository3 "github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	"github.com/devtron-labs/devtron/pkg/auth/user"
	bean2 "github.com/devtron-labs/devtron/pkg/bean"
	"github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	artifactPromotion2 "github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion"
	"github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/bean"
	repository2 "github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/repository"
	"github.com/devtron-labs/devtron/util/rbac"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type PromotionApprovalRequestRestHandler interface {
	HandleArtifactPromotionRequest(w http.ResponseWriter, r *http.Request)
	GetByPromotionRequestId(w http.ResponseWriter, r *http.Request)
}

type PromotionApprovalMaterialRestHandler interface {
	GetArtifactsForPromotion(w http.ResponseWriter, r *http.Request)
}

type PromotionApprovalRestHandlerImpl struct {
	promotionApprovalRequestService            artifactPromotion2.ArtifactPromotionApprovalService
	logger                                     *zap.SugaredLogger
	userService                                user.UserService
	enforcer                                   casbin.Enforcer
	validator                                  *validator.Validate
	userCommonService                          user.UserCommonService
	enforcerUtil                               rbac.EnforcerUtil
	environmentRepository                      repository.EnvironmentRepository
	artifactPromotionApprovalRequestRepository repository2.ArtifactPromotionApprovalRequestRepository
	appArtifactManager                         pipeline.AppArtifactManager
	CiArtifactRepository                       repository3.CiArtifactRepository
}

func NewArtifactPromotionApprovalRestHandlerImpl(
	promotionApprovalRequestService artifactPromotion2.ArtifactPromotionApprovalService,
	logger *zap.SugaredLogger,
	userService user.UserService,
	validator *validator.Validate,
	userCommonService user.UserCommonService,
	enforcerUtil rbac.EnforcerUtil,
	environmentRepository repository.EnvironmentRepository,
	artifactPromotionApprovalRequestRepository repository2.ArtifactPromotionApprovalRequestRepository,
	appArtifactManager pipeline.AppArtifactManager,
) *PromotionApprovalRestHandlerImpl {
	return &PromotionApprovalRestHandlerImpl{
		promotionApprovalRequestService: promotionApprovalRequestService,
		logger:                          logger,
		userService:                     userService,
		validator:                       validator,
		userCommonService:               userCommonService,
		enforcerUtil:                    enforcerUtil,
		environmentRepository:           environmentRepository,
		artifactPromotionApprovalRequestRepository: artifactPromotionApprovalRequestRepository,
		appArtifactManager:                         appArtifactManager,
	}
}

func (handler PromotionApprovalRestHandlerImpl) HandleArtifactPromotionRequest(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userService.GetLoggedInUser(r)
	if err != nil || userId == 0 {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("token")
	isAuthorised, err := handler.userService.IsUserAdminOrManagerForAnyApp(userId, token)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	if !isAuthorised {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}
	var promotionRequest bean.ArtifactPromotionRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&promotionRequest)
	if err != nil {
		handler.logger.Errorw("err in decoding request in promotionRequest", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	authorizedEnvironments := make(map[string]bool)

	switch promotionRequest.Action {
	case bean.ACTION_PROMOTE:

		appName := promotionRequest.AppName
		appRbacObject := handler.enforcerUtil.GetAppRBACName(appName)
		ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, appRbacObject)
		if !ok {
			common.WriteJsonResp(w, err, nil, http.StatusForbidden)
			return
		}

		environmentNames := promotionRequest.EnvironmentNames
		envRbacObjectMap := handler.enforcerUtil.GetEnvRBACByAppNameAndEnvNames(appName, environmentNames)
		envObjectArr := make([]string, 0)
		for _, obj := range envObjectArr {
			envObjectArr = append(envObjectArr, obj)
		}
		results := handler.enforcer.EnforceInBatch(token, casbin.ResourceEnvironment, casbin.ActionTrigger, envObjectArr)
		for _, env := range environmentNames {
			rbacObject := envRbacObjectMap[env]
			isAuthorised = results[rbacObject]
			authorizedEnvironments[env] = isAuthorised
		}

	case bean.ACTION_APPROVE:
		appName := promotionRequest.AppName
		environmentNames := promotionRequest.EnvironmentNames
		teamEnvRbacObjectMap := handler.enforcerUtil.GetTeamEnvRbacObjByAppAndEnvNames(appName, environmentNames)
		teamEnvObjectArr := make([]string, 0)
		for _, obj := range teamEnvObjectArr {
			teamEnvObjectArr = append(teamEnvObjectArr, obj)
		}
		results := handler.enforcer.EnforceInBatch(token, casbin.ResourceApprovalPolicy, casbin.ActionArtifactPromote, teamEnvObjectArr)
		for _, env := range promotionRequest.EnvironmentNames {
			rbacObject := teamEnvRbacObjectMap[env]
			isAuthorised = results[rbacObject]
			authorizedEnvironments[env] = isAuthorised
		}

	case bean.ACTION_CANCEL:
		artifactPromotionDao, err := handler.artifactPromotionApprovalRequestRepository.FindById(promotionRequest.PromotionRequestId)
		if err == pg.ErrNoRows {
			handler.logger.Errorw("promotion request for given id does not exist", "promotionRequestId", promotionRequest.PromotionRequestId, "err", err)
			common.WriteJsonResp(w, errors.New("promotion request for given id does not exist"), nil, http.StatusNotFound)
			return
		}
		if err != nil {
			handler.logger.Errorw("error in fetching artifact promotion request by id", "artifactPromotionRequestId", promotionRequest.PromotionRequestId, "err", err)
			return
		}
		appRbacObj, envRbacObj := handler.getAppAndEnvObjectByCdPipelineId(artifactPromotionDao.DestinationPipelineId)
		if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionTrigger, appRbacObj); !ok {
			common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
			return
		}
		if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionTrigger, envRbacObj); !ok {
			common.WriteJsonResp(w, err, "Unauthorized User", http.StatusForbidden)
			return
		}
	}

	_, err = handler.promotionApprovalRequestService.HandleArtifactPromotionRequest(&promotionRequest, authorizedEnvironments)
	if err != nil {
		handler.logger.Errorw("error in handling promotion artifact request", "promotionRequest", promotionRequest, "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, nil, http.StatusOK)
}

func (handler PromotionApprovalRestHandlerImpl) GetByPromotionRequestId(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userService.GetLoggedInUser(r)
	if err != nil || userId == 0 {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("token")
	isAuthorised, err := handler.userService.IsUserAdminOrManagerForAnyApp(userId, token)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	if !isAuthorised {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	promotionRequestId, err := strconv.Atoi(vars["promotionRequestId"])
	if err != nil {
		handler.logger.Errorw("error in parsing promotionRequestId from string to int", "promotionRequestId", vars["promotionRequestId"])
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	artifactPromotionDao, err := handler.artifactPromotionApprovalRequestRepository.FindById(promotionRequestId)
	if err == pg.ErrNoRows {
		handler.logger.Errorw("promotion request for given id does not exist", "promotionRequestId", promotionRequestId, "err", err)
		common.WriteJsonResp(w, errors.New("promotion request for given id does not exist"), nil, http.StatusNotFound)
		return
	}
	if err != nil {
		handler.logger.Errorw("error in fetching artifact promotion request by id", "artifactPromotionRequestId", promotionRequestId, "err", err)
		return
	}

	// rbac block starts from here
	appRbacObj, envRbacObj := handler.getAppAndEnvObjectByCdPipelineId(artifactPromotionDao.DestinationPipelineId)
	if ok := handler.enforcer.Enforce(token, casbin.ResourceApplications, casbin.ActionGet, appRbacObj); !ok {
		common.WriteJsonResp(w, fmt.Errorf("unauthorized user"), "Unauthorized User", http.StatusForbidden)
		return
	}
	if ok := handler.enforcer.Enforce(token, casbin.ResourceEnvironment, casbin.ActionGet, envRbacObj); !ok {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusForbidden)
		return
	}

	_, err = handler.promotionApprovalRequestService.GetByPromotionRequestId(artifactPromotionDao)
	if err != nil {
		handler.logger.Errorw("error in getting data for promotion request id", "promotionRequestId", promotionRequestId, "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}

	common.WriteJsonResp(w, nil, nil, http.StatusOK)
}

func (handler PromotionApprovalRestHandlerImpl) getAppAndEnvObjectByCdPipelineId(cdPipelineId int) (string, string) {
	object := handler.enforcerUtil.GetAppAndEnvObjectByPipelineIds([]int{cdPipelineId})
	rbacObjects := object[cdPipelineId]
	return rbacObjects[0], rbacObjects[1]
}

func (handler PromotionApprovalRestHandlerImpl) GetArtifactsForPromotion(w http.ResponseWriter, r *http.Request) {

	userId, err := handler.userService.GetLoggedInUser(r)
	if err != nil || userId == 0 {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("token")
	isAuthorised, err := handler.userService.IsUserAdminOrManagerForAnyApp(userId, token)
	if err != nil {
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	if !isAuthorised {
		common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
		return
	}

	queryParams := r.URL.Query()

	resource := queryParams.Get("resource")
	resourceName := queryParams.Get("resourceName")

	appId, err := strconv.Atoi(queryParams.Get("appId"))
	if err != nil {
		handler.logger.Errorw("error in parsing appId from string to int", "appId", queryParams.Get("appId"))
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	pendingForCurrentUser, err := strconv.ParseBool(queryParams.Get("pendingForCurrentUser"))
	if err != nil {
		handler.logger.Errorw("error in parsing pendingForCurrentUser from string to bool", "pendingForCurrentUser", queryParams.Get("pendingForCurrentUser"))
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	workflowId, err := strconv.Atoi(queryParams.Get("workflowId"))
	if err != nil {
		handler.logger.Errorw("error in parsing workflowId from string to int", "workflowId", queryParams.Get("workflowId"))
		common.WriteJsonResp(w, err, nil, http.StatusBadRequest)
		return
	}

	offset := 0
	limit := 10
	offsetQueryParam := queryParams.Get("offset")
	if offsetQueryParam != "" {
		offset, err = strconv.Atoi(offsetQueryParam)
		if err != nil || offset < 0 {
			handler.logger.Errorw("error in parsing ", "err", err, "offsetQueryParam", offsetQueryParam)
			common.WriteJsonResp(w, err, "invalid offset", http.StatusBadRequest)
			return
		}
	}

	sizeQueryParam := r.URL.Query().Get("size")
	if sizeQueryParam != "" {
		limit, err = strconv.Atoi(sizeQueryParam)
		if err != nil {
			handler.logger.Errorw("request err, GetArtifactsForRollback", "err", err, "sizeQueryParam", sizeQueryParam)
			common.WriteJsonResp(w, err, "invalid size", http.StatusBadRequest)
			return
		}
	}

	if resourceName == "CI" {

	}

	artifactPromotionMaterialRequest := bean2.ArtifactPromotionMaterialRequest{}
	artifactPromotionMaterialRequest.Resource = resource
	artifactPromotionMaterialRequest.ResourceName = resourceName
	artifactPromotionMaterialRequest.AppId = appId
	artifactPromotionMaterialRequest.WorkflowId = workflowId
	artifactPromotionMaterialRequest.Offset = offset
	artifactPromotionMaterialRequest.Limit = limit

	artifactPromotionMaterialRequest.PendingForCurrentUser = pendingForCurrentUser
	artifactPromotionMaterialRequest.UserId = userId

	artifactPromotionMaterialResponse, err := handler.appArtifactManager.FetchMaterialForArtifactPromotion(artifactPromotionMaterialRequest)
	if err != nil {
		handler.logger.Errorw("error in fetching artifacts for given promotion request parameters", "artifactPromotionRequest", artifactPromotionMaterialRequest, "err", err)
		common.WriteJsonResp(w, errors.New("error in fetching artifacts response for given request parameters"), nil, http.StatusInternalServerError)
		return
	}

	common.WriteJsonResp(w, nil, artifactPromotionMaterialResponse, http.StatusOK)
}
