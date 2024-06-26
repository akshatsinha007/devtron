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

package user

import (
	"errors"
	"net/http"

	"github.com/devtron-labs/devtron/api/restHandler/common"
	"github.com/devtron-labs/devtron/pkg/auth/authorisation/casbin"
	user2 "github.com/devtron-labs/devtron/pkg/auth/user"
	"github.com/devtron-labs/devtron/util/rbac"
	"go.uber.org/zap"
	"gopkg.in/go-playground/validator.v9"
)

type RbacRoleRestHandler interface {
	GetAllDefaultRoles(w http.ResponseWriter, r *http.Request)
}

type RbacRoleRestHandlerImpl struct {
	logger          *zap.SugaredLogger
	validator       *validator.Validate
	rbacRoleService user2.RbacRoleService
	userService     user2.UserService
	enforcer        casbin.Enforcer
	enforcerUtil    rbac.EnforcerUtil
}

func NewRbacRoleHandlerImpl(logger *zap.SugaredLogger,
	validator *validator.Validate, rbacRoleService user2.RbacRoleService,
	userService user2.UserService, enforcer casbin.Enforcer,
	enforcerUtil rbac.EnforcerUtil) *RbacRoleRestHandlerImpl {
	rbacRoleRestHandlerImpl := &RbacRoleRestHandlerImpl{
		logger:          logger,
		validator:       validator,
		rbacRoleService: rbacRoleService,
		userService:     userService,
		enforcer:        enforcer,
		enforcerUtil:    enforcerUtil,
	}
	return rbacRoleRestHandlerImpl
}

func (handler *RbacRoleRestHandlerImpl) GetAllDefaultRoles(w http.ResponseWriter, r *http.Request) {
	userId, err := handler.userService.GetLoggedInUser(r)
	if userId == 0 || err != nil {
		common.WriteJsonResp(w, err, "Unauthorized User", http.StatusUnauthorized)
		return
	}
	handler.logger.Debugw("request payload, GetAllDefaultRoles")
	// RBAC enforcer applying
	token := r.Header.Get("token")
	teamNames, err := handler.enforcerUtil.GetAllActiveTeamNames()
	if err != nil {
		handler.logger.Errorw("error in finding all active team names", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	if len(teamNames) > 0 {
		rbacResultMap := handler.enforcer.EnforceInBatch(token, casbin.ResourceUser, casbin.ActionGet, teamNames)
		isAuthorized := false
		for _, authorizedOnTeam := range rbacResultMap {
			if authorizedOnTeam {
				isAuthorized = true
				break
			}
		}
		if !isAuthorized {
			common.WriteJsonResp(w, errors.New("unauthorized"), nil, http.StatusForbidden)
			return
		}
	}
	roles, err := handler.rbacRoleService.GetAllDefaultRoles()
	if err != nil {
		handler.logger.Errorw("service error, GetAllDefaultRoles", "err", err)
		common.WriteJsonResp(w, err, nil, http.StatusInternalServerError)
		return
	}
	common.WriteJsonResp(w, nil, roles, http.StatusOK)
}
