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

/*
@description: user authentication and authorization
*/
package repository

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/pkg/sql"
	bean2 "github.com/devtron-labs/devtron/pkg/user/bean"
	casbin2 "github.com/devtron-labs/devtron/pkg/user/casbin"
	"github.com/devtron-labs/devtron/util"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"strings"
	"time"
)

type UserAuthRepository interface {
	CreateRole(role *RoleModel) (*RoleModel, error)
	CreateRoleWithTxn(userModel *RoleModel, tx *pg.Tx) (*RoleModel, error)
	GetRoleById(id int) (*RoleModel, error)
	GetRolesByIds(ids []int) ([]RoleModel, error)
	GetRoleByRoles(roles []string) ([]RoleModel, error)
	GetRolesByUserId(userId int32) ([]RoleModel, error)
	GetRolesByGroupId(userId int32) ([]*RoleModel, error)
	GetAllRole() ([]RoleModel, error)
	GetRolesByActionAndAccessType(action string, accessType string) ([]RoleModel, error)
	GetRoleByFilterForAllTypes(entity, team, app, env, act string, approver bool, accessType, cluster, namespace, group, kind, resource, action string, oldValues bool, workflow string) (RoleModel, error)
	CreateUserRoleMapping(userRoleModel *UserRoleModel, tx *pg.Tx) (*UserRoleModel, error)
	GetUserRoleMappingByUserId(userId int32) ([]*UserRoleModel, error)
	DeleteUserRoleMapping(userRoleModel *UserRoleModel, tx *pg.Tx) (bool, error)
	DeleteUserRoleByRoleId(roleId int, tx *pg.Tx) error
	DeleteUserRoleByRoleIds(roleIds []int, tx *pg.Tx) error
	CreateDefaultPoliciesForAllTypes(team, entityName, env, entity, cluster, namespace, group, kind, resource, actionType, accessType string, approver bool, UserId int32) (bool, error, []casbin2.Policy)
	CreateRoleForSuperAdminIfNotExists(tx *pg.Tx, UserId int32) (bool, error)
	SyncOrchestratorToCasbin(team string, entityName string, env string, tx *pg.Tx) (bool, error)
	GetRolesForEnvironment(envName, envIdentifier string) ([]*RoleModel, error)
	GetRolesForProject(teamName string) ([]*RoleModel, error)
	GetRolesForApp(appName string) ([]*RoleModel, error)
	GetRolesForChartGroup(chartGroupName string) ([]*RoleModel, error)
	DeleteRole(role *RoleModel, tx *pg.Tx) error
	DeleteRolesByIds(roleIds []int, tx *pg.Tx) error
	//GetRoleByFilterForClusterEntity(cluster, namespace, group, kind, resource, action string) (RoleModel, error)
	GetRolesByUserIdAndEntityType(userId int32, entityType string) ([]*RoleModel, error)
	CreateRolesWithAccessTypeAndEntity(team, entityName, env, entity, cluster, namespace, group, kind, resource, actionType, accessType string, UserId int32, role string) (bool, error)
	GetRolesByEntityAccessTypeAndAction(entity, accessType, action string) ([]*RoleModel, error)
	GetApprovalUsersByEnv(appName, envName string) ([]string, []string, error)
	GetConfigApprovalUsersByEnv(appName, envName string) ([]string, []string, error)
	GetRolesForWorkflow(workflow, entityName string) ([]*RoleModel, error)
	GetRoleForClusterEntity(cluster, namespace, group, kind, resource, action string) (RoleModel, error)
	GetRoleForJobsEntity(entity, team, app, env, act string, workflow string) (RoleModel, error)
	GetRoleForOtherEntity(team, app, env, act, accessType string, oldValues, approver bool) (RoleModel, error)
	GetRoleForChartGroupEntity(entity, app, act, accessType string) (RoleModel, error)
}

type UserAuthRepositoryImpl struct {
	dbConnection                *pg.DB
	Logger                      *zap.SugaredLogger
	defaultAuthPolicyRepository DefaultAuthPolicyRepository
	defaultAuthRoleRepository   DefaultAuthRoleRepository
}

func NewUserAuthRepositoryImpl(dbConnection *pg.DB, Logger *zap.SugaredLogger,
	defaultAuthPolicyRepository DefaultAuthPolicyRepository,
	defaultAuthRoleRepository DefaultAuthRoleRepository) *UserAuthRepositoryImpl {
	return &UserAuthRepositoryImpl{
		dbConnection:                dbConnection,
		Logger:                      Logger,
		defaultAuthPolicyRepository: defaultAuthPolicyRepository,
		defaultAuthRoleRepository:   defaultAuthRoleRepository,
	}
}

type RoleModel struct {
	TableName   struct{} `sql:"roles" pg:",discard_unknown_columns"`
	Id          int      `sql:"id,pk"`
	Role        string   `sql:"role,notnull"`
	Entity      string   `sql:"entity"`
	Team        string   `sql:"team"`
	EntityName  string   `sql:"entity_name"`
	Environment string   `sql:"environment"`
	Action      string   `sql:"action"`
	AccessType  string   `sql:"access_type"`
	Approver    bool     `sql:"approver"`
	Cluster     string   `sql:"cluster"`
	Namespace   string   `sql:"namespace"`
	Group       string   `sql:"group"`
	Kind        string   `sql:"kind"`
	Resource    string   `sql:"resource"`
	Workflow    string   `sql:"workflow"`
	sql.AuditLog
}

type RolePolicyDetails struct {
	Team       string
	Env        string
	App        string
	TeamObj    string
	EnvObj     string
	AppObj     string
	Entity     string
	EntityName string

	Cluster      string
	Namespace    string
	Group        string
	Kind         string
	Resource     string
	ClusterObj   string
	NamespaceObj string
	GroupObj     string
	KindObj      string
	ResourceObj  string
	Approver     bool
}

type ClusterRolePolicyDetails struct {
	Entity       string
	Cluster      string
	Namespace    string
	Group        string
	Kind         string
	Resource     string
	ClusterObj   string
	NamespaceObj string
	GroupObj     string
	KindObj      string
	ResourceObj  string
}

func (impl UserAuthRepositoryImpl) CreateRole(role *RoleModel) (*RoleModel, error) {
	err := impl.dbConnection.Insert(role)
	if err != nil {
		impl.Logger.Error("error in creating role", "err", err, "role", role)
		return role, err
	}
	return role, nil
}

func (impl UserAuthRepositoryImpl) CreateRoleWithTxn(userModel *RoleModel, tx *pg.Tx) (*RoleModel, error) {
	err := tx.Insert(userModel)
	if err != nil {
		impl.Logger.Error(err)
		return userModel, err
	}
	return userModel, nil
}

func (impl UserAuthRepositoryImpl) GetRoleById(id int) (*RoleModel, error) {
	var model RoleModel
	err := impl.dbConnection.Model(&model).Where("id = ?", id).Select()
	if err != nil {
		impl.Logger.Error(err)
		return &model, err
	}
	return &model, nil
}
func (impl UserAuthRepositoryImpl) GetRolesByIds(ids []int) ([]RoleModel, error) {
	var model []RoleModel
	err := impl.dbConnection.Model(&model).Where("id IN (?)", pg.In(ids)).Select()
	if err != nil {
		impl.Logger.Error(err)
		return model, err
	}
	return model, nil
}
func (impl UserAuthRepositoryImpl) GetRoleByRoles(roles []string) ([]RoleModel, error) {
	var model []RoleModel
	err := impl.dbConnection.Model(&model).Where("role IN (?)", pg.In(roles)).Select()
	if err != nil {
		impl.Logger.Error(err)
		return model, err
	}
	return model, nil
}

func (impl UserAuthRepositoryImpl) GetRolesByUserId(userId int32) ([]RoleModel, error) {
	var models []RoleModel
	err := impl.dbConnection.Model(&models).
		Column("role_model.*").
		Join("INNER JOIN user_roles ur on ur.role_id=role_model.id").
		Where("ur.user_id = ?", userId).Select()
	if err != nil {
		impl.Logger.Error(err)
		return models, err
	}
	return models, nil
}
func (impl UserAuthRepositoryImpl) GetRolesByGroupId(roleGroupId int32) ([]*RoleModel, error) {
	var models []*RoleModel
	err := impl.dbConnection.Model(&models).
		Column("role_model.*").
		Join("INNER JOIN role_group_role_mapping rgrm on rgrm.role_id=role_model.id").
		Join("INNER JOIN role_group rg on rg.id=rgrm.role_group_id").
		Where("rg.id = ?", roleGroupId).Select()
	if err != nil {
		impl.Logger.Error(err)
		return models, err
	}
	return models, nil
}
func (impl UserAuthRepositoryImpl) GetRole(role string) (*RoleModel, error) {
	var model RoleModel
	err := impl.dbConnection.Model(&model).Where("role = ?", role).Select()
	if err != nil {
		impl.Logger.Error(err)
		return &model, err
	}
	return &model, nil
}
func (impl UserAuthRepositoryImpl) GetAllRole() ([]RoleModel, error) {
	var models []RoleModel
	err := impl.dbConnection.Model(&models).Select()
	if err != nil {
		impl.Logger.Error(err)
		return models, err
	}
	return models, nil
}

func (impl UserAuthRepositoryImpl) GetRolesByActionAndAccessType(action string, accessType string) ([]RoleModel, error) {
	var models []RoleModel
	var err error
	if accessType == "" {
		err = impl.dbConnection.Model(&models).Where("action = ?", action).
			Where("access_type is NULL").
			Select()
	} else {
		err = impl.dbConnection.Model(&models).Where("action = ?", action).
			Where("access_type = ?", accessType).
			Select()
	}
	if err != nil {
		impl.Logger.Error("err in getting role by action", "err", err, "action", action, "accessType", accessType)
		return models, err
	}
	return models, nil
}

func (impl UserAuthRepositoryImpl) GetRoleByFilterForAllTypes(entity, team, app, env, act string, approver bool, accessType, cluster, namespace, group, kind, resource, action string, oldValues bool, workflow string) (RoleModel, error) {
	switch entity {
	case bean2.CLUSTER:
		{
			return impl.GetRoleForClusterEntity(cluster, namespace, group, kind, resource, action)
		}
	case bean.CHART_GROUP_ENTITY:
		{
			return impl.GetRoleForChartGroupEntity(entity, app, act, accessType)
		}
	case bean2.EntityJobs:
		{
			return impl.GetRoleForJobsEntity(entity, team, app, env, act, workflow)
		}
	default:
		{
			return impl.GetRoleForOtherEntity(team, app, env, act, accessType, oldValues, approver)
		}
	}
	return RoleModel{}, nil
}

func (impl UserAuthRepositoryImpl) CreateUserRoleMapping(userRoleModel *UserRoleModel, tx *pg.Tx) (*UserRoleModel, error) {
	err := tx.Insert(userRoleModel)
	if err != nil {
		impl.Logger.Error(err)
		return userRoleModel, err
	}

	return userRoleModel, nil
}
func (impl UserAuthRepositoryImpl) GetUserRoleMappingByUserId(userId int32) ([]*UserRoleModel, error) {
	var userRoleModels []*UserRoleModel
	err := impl.dbConnection.Model(&userRoleModels).Where("user_id = ?", userId).Select()
	if err != nil {
		impl.Logger.Error(err)
		return userRoleModels, err
	}
	return userRoleModels, nil
}
func (impl UserAuthRepositoryImpl) DeleteUserRoleMapping(userRoleModel *UserRoleModel, tx *pg.Tx) (bool, error) {
	err := tx.Delete(userRoleModel)
	if err != nil {
		impl.Logger.Error(err)
		return false, err
	}
	return true, nil
}

func (impl UserAuthRepositoryImpl) DeleteUserRoleByRoleId(roleId int, tx *pg.Tx) error {
	var userRoleModel *UserRoleModel
	_, err := tx.Model(userRoleModel).
		Where("role_id = ?", roleId).Delete()
	if err != nil {
		impl.Logger.Error("err in deleting user role by role id", "err", err, "roleId", roleId)
		return err
	}
	return nil
}
func (impl UserAuthRepositoryImpl) DeleteUserRoleByRoleIds(roleIds []int, tx *pg.Tx) error {
	var userRoleModel *UserRoleModel
	_, err := tx.Model(userRoleModel).
		Where("role_id in (?)", pg.In(roleIds)).Delete()
	if err != nil {
		impl.Logger.Error("err in deleting user role by role id", "err", err, "roleIds", roleIds)
		return err
	}
	return nil
}

func (impl UserAuthRepositoryImpl) CreateDefaultPoliciesForAllTypes(team, entityName, env, entity, cluster, namespace, group, kind, resource, actionType, accessType string, approver bool, UserId int32) (bool, error, []casbin2.Policy) {
	//not using txn from parent caller because of conflicts in fetching of transactional save
	dbConnection := impl.dbConnection
	tx, err := dbConnection.Begin()
	var policiesToBeAdded []casbin2.Policy
	if err != nil {
		return false, err, policiesToBeAdded
	}
	// Rollback tx on error.
	defer tx.Rollback()

	//for START in Casbin Object
	teamObj := team
	envObj := env
	appObj := entityName
	if teamObj == "" {
		teamObj = "*"
	}
	if envObj == "" {
		envObj = "*"
	}
	if appObj == "" {
		appObj = "*"
	}

	clusterObj := cluster
	namespaceObj := namespace
	groupObj := group
	kindObj := kind
	resourceObj := resource

	if cluster == "" {
		clusterObj = "*"
	}
	if namespace == "" {
		namespaceObj = "*"
	}
	if group == "" {
		groupObj = "*"
	}
	if kind == "" {
		kindObj = "*"
	}
	if resource == "" {
		resourceObj = "*"
	}
	rolePolicyDetails := RolePolicyDetails{
		Team:         team,
		App:          entityName,
		Env:          env,
		TeamObj:      teamObj,
		EnvObj:       envObj,
		AppObj:       appObj,
		Entity:       entity,
		EntityName:   entityName,
		Approver:     approver,
		Cluster:      cluster,
		Namespace:    namespace,
		Group:        group,
		Kind:         kind,
		Resource:     resource,
		ClusterObj:   clusterObj,
		NamespaceObj: namespaceObj,
		GroupObj:     groupObj,
		KindObj:      kindObj,
		ResourceObj:  resourceObj,
	}

	//getting policies from db
	PoliciesDb, err := impl.defaultAuthPolicyRepository.GetPolicyByRoleTypeAndEntity(bean2.RoleType(actionType), accessType, entity)
	if err != nil {
		return false, err, policiesToBeAdded
	}
	//getting updated policies
	Policies, err := util.Tprintf(PoliciesDb, rolePolicyDetails)
	if err != nil {
		impl.Logger.Errorw("error in getting updated policies", "err", err, "roleType", bean2.RoleType(actionType), accessType)
		return false, err, policiesToBeAdded
	}
	//for START in Casbin Object Ends Here
	var policies bean.PolicyRequest
	err = json.Unmarshal([]byte(Policies), &policies)
	if err != nil {
		impl.Logger.Errorw("decode err", "err", err)
		return false, err, policiesToBeAdded
	}
	impl.Logger.Debugw("add policy request", "policies", policies)
	policiesToBeAdded = append(policiesToBeAdded, policies.Data...)
	//Creating ROLES
	//getting roles from db
	roleDb, err := impl.defaultAuthRoleRepository.GetRoleByRoleTypeAndEntityType(bean2.RoleType(actionType), accessType, entity)
	if err != nil {
		return false, err, nil
	}
	role, err := util.Tprintf(roleDb, rolePolicyDetails)
	if err != nil {
		impl.Logger.Errorw("error in getting updated role", "err", err, "roleType", bean2.RoleType(actionType))
		return false, err, nil
	}
	//getting updated role
	var roleData bean.RoleData
	err = json.Unmarshal([]byte(role), &roleData)
	if err != nil {
		impl.Logger.Errorw("decode err", "err", err)
		return false, err, nil
	}
	_, err = impl.createRole(&roleData, UserId)
	if err != nil && strings.Contains("duplicate key value violates unique constraint", err.Error()) {
		return false, err, nil
	}
	err = tx.Commit()
	if err != nil {
		return false, err, nil
	}
	return true, nil, policiesToBeAdded
}

func (impl UserAuthRepositoryImpl) GetApprovalUsersByEnv(appName, envName string) ([]string, []string, error) {
	var emailIds []string
	var roleGroups []string

	query := "select distinct(email_id) from users us inner join user_roles ur on us.id=ur.user_id inner join roles on ur.role_id = roles.id " +
		"where ((roles.approver = true and (roles.environment=? OR roles.environment is null) and (entity_name=? OR entity_name is null)) OR roles.role = ?) " +
		"and us.id not in (1);"
	_, err := impl.dbConnection.Query(&emailIds, query, envName, appName, "role:super-admin___")
	if err != nil && err != pg.ErrNoRows {
		return emailIds, roleGroups, err
	}

	roleGroupQuery := "select rg.casbin_name from role_group rg inner join role_group_role_mapping rgrm on rg.id = rgrm.role_group_id " +
		"inner join roles r on rgrm.role_id = r.id where r.approver = true  and r.environment=? and r.entity_name=?;"
	_, err = impl.dbConnection.Query(&roleGroups, roleGroupQuery, envName, appName)
	if err != nil && err != pg.ErrNoRows {
		return emailIds, roleGroups, err
	}

	return emailIds, roleGroups, nil
}

func (impl UserAuthRepositoryImpl) GetConfigApprovalUsersByEnv(appName, envName string) ([]string, []string, error) {
	var emailIds []string
	var roleGroups []string

	query := "select distinct(email_id) from users us inner join user_roles ur on us.id=ur.user_id inner join roles on ur.role_id = roles.id " +
		"where ((roles.action = ? and (roles.environment=? OR roles.environment is null) and (entity_name=? OR entity_name is null)) OR roles.role = ?) " +
		"and us.id not in (1);"
	_, err := impl.dbConnection.Query(&emailIds, query, "configApprover", envName, appName, "role:super-admin___")
	if err != nil && err != pg.ErrNoRows {
		return emailIds, roleGroups, err
	}

	roleGroupQuery := "select rg.casbin_name from role_group rg inner join role_group_role_mapping rgrm on rg.id = rgrm.role_group_id " +
		"inner join roles r on rgrm.role_id = r.id where r.action = ?  and r.environment=? and r.entity_name=?;"
	_, err = impl.dbConnection.Query(&roleGroups, roleGroupQuery, "configApprover", envName, appName)
	if err != nil && err != pg.ErrNoRows {
		return emailIds, roleGroups, err
	}

	return emailIds, roleGroups, nil
}

func (impl UserAuthRepositoryImpl) CreateRolesWithAccessTypeAndEntity(team, entityName, env, entity, cluster, namespace, group, kind, resource, actionType, accessType string, UserId int32, role string) (bool, error) {
	roleData := bean.RoleData{
		Role:        role,
		Entity:      entity,
		Team:        team,
		EntityName:  entityName,
		Environment: env,
		Action:      actionType,
		AccessType:  accessType,
		Cluster:     cluster,
		Namespace:   namespace,
		Group:       group,
		Kind:        kind,
		Resource:    resource,
	}
	_, err := impl.createRole(&roleData, UserId)
	if err != nil && strings.Contains("duplicate key value violates unique constraint", err.Error()) {
		return false, err
	}
	return true, nil
}

func (impl UserAuthRepositoryImpl) CreateRoleForSuperAdminIfNotExists(tx *pg.Tx, UserId int32) (bool, error) {
	transaction, err := impl.dbConnection.Begin()
	if err != nil {
		return false, err
	}

	//Creating ROLES
	roleModel, err := impl.GetRoleByFilterForAllTypes("", "", "", "", bean2.SUPER_ADMIN, false, "", "", "", "", "", "", "", false, "")
	if err != nil && err != pg.ErrNoRows {
		return false, err
	}
	if roleModel.Id == 0 || err == pg.ErrNoRows {
		roleManager := "{\r\n    \"role\": \"role:super-admin___\",\r\n    \"casbinSubjects\": [\r\n        \"role:super-admin___\"\r\n    ],\r\n    \"team\": \"\",\r\n    \"entityName\": \"\",\r\n    \"environment\": \"\",\r\n    \"action\": \"super-admin\"\r\n}"

		var roleManagerData bean.RoleData
		err = json.Unmarshal([]byte(roleManager), &roleManagerData)
		if err != nil {
			impl.Logger.Errorw("decode err", "err", err)
			return false, err
		}
		_, err = impl.createRole(&roleManagerData, UserId)
		if err != nil && strings.Contains("duplicate key value violates unique constraint", err.Error()) {
			return false, err
		}
	}
	err = transaction.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (impl UserAuthRepositoryImpl) createRole(roleData *bean.RoleData, UserId int32) (bool, error) {
	roleModel := &RoleModel{
		Role:        roleData.Role,
		Entity:      roleData.Entity,
		Team:        roleData.Team,
		EntityName:  roleData.EntityName,
		Environment: roleData.Environment,
		Action:      roleData.Action,
		AccessType:  roleData.AccessType,
		Approver:    roleData.Approver,
		Cluster:     roleData.Cluster,
		Namespace:   roleData.Namespace,
		Group:       roleData.Group,
		Kind:        roleData.Kind,
		Resource:    roleData.Resource,
		AuditLog: sql.AuditLog{
			CreatedBy: UserId,
			CreatedOn: time.Now(),
			UpdatedBy: UserId,
			UpdatedOn: time.Now(),
		},
	}
	roleModel, err := impl.CreateRole(roleModel)
	if err != nil || roleModel == nil {
		return false, err
	}
	return true, nil
}

func (impl UserAuthRepositoryImpl) SyncOrchestratorToCasbin(team string, entityName string, env string, tx *pg.Tx) (bool, error) {

	//getting policies from db
	triggerPoliciesDb, err := impl.defaultAuthPolicyRepository.GetPolicyByRoleTypeAndEntity(bean2.TRIGGER_TYPE, bean2.DEVTRON_APP, bean2.ENTITY_APPS)
	if err != nil {
		return false, err
	}
	viewPoliciesDb, err := impl.defaultAuthPolicyRepository.GetPolicyByRoleTypeAndEntity(bean2.VIEW_TYPE, bean2.DEVTRON_APP, bean2.ENTITY_APPS)
	if err != nil {
		return false, err
	}

	//for START in Casbin Object
	teamObj := team
	envObj := env
	appObj := entityName
	if teamObj == "" {
		teamObj = "*"
	}
	if envObj == "" {
		envObj = "*"
	}
	if appObj == "" {
		appObj = "*"
	}

	policyDetails := RolePolicyDetails{
		Team:    team,
		App:     entityName,
		Env:     env,
		TeamObj: teamObj,
		EnvObj:  envObj,
		AppObj:  appObj,
	}

	//getting updated trigger policies
	triggerPolicies, err := util.Tprintf(triggerPoliciesDb, policyDetails)
	if err != nil {
		impl.Logger.Errorw("error in getting updated policies", "err", err, "roleType", bean2.TRIGGER_TYPE)
		return false, err
	}

	//getting updated view policies
	viewPolicies, err := util.Tprintf(viewPoliciesDb, policyDetails)
	if err != nil {
		impl.Logger.Errorw("error in getting updated policies", "err", err, "roleType", bean2.VIEW_TYPE)
		return false, err
	}

	//for START in Casbin Object Ends Here
	var policies []casbin2.Policy
	var policiesTrigger bean.PolicyRequest
	err = json.Unmarshal([]byte(triggerPolicies), &policiesTrigger)
	if err != nil {
		impl.Logger.Errorw("decode err", "err", err)
		return false, err
	}
	impl.Logger.Debugw("add policy request", "policies", policiesTrigger)
	policies = append(policies, policiesTrigger.Data...)
	var policiesView bean.PolicyRequest
	err = json.Unmarshal([]byte(viewPolicies), &policiesView)
	if err != nil {
		impl.Logger.Errorw("decode err", "err", err)
		return false, err
	}
	impl.Logger.Debugw("add policy request", "policies", policiesView)
	policies = append(policies, policiesView.Data...)
	err = casbin2.AddPolicy(policies)
	if err != nil {
		impl.Logger.Errorw("casbin policy addition failed", "err", err)
		return false, err
	}
	return true, nil
}

func (impl UserAuthRepositoryImpl) GetDefaultPolicyByRoleType(roleType bean2.RoleType) (policy string, err error) {
	policy, err = impl.defaultAuthPolicyRepository.GetPolicyByRoleTypeAndEntity(roleType, bean2.DEVTRON_APP, bean2.ENTITY_APPS)
	if err != nil {
		return "", err
	}
	return policy, nil
}

func (impl UserAuthRepositoryImpl) GetRolesForEnvironment(envName, envIdentifier string) ([]*RoleModel, error) {
	var roles []*RoleModel
	err := impl.dbConnection.Model(&roles).WhereOr("environment = ?", envName).
		WhereOr("environment = ?", envIdentifier).Select()
	if err != nil {
		impl.Logger.Errorw("error in getting roles for environment", "err", err, "envName", envName, "envIdentifier", envIdentifier)
		return nil, err
	}
	return roles, nil
}

func (impl UserAuthRepositoryImpl) GetRolesForProject(teamName string) ([]*RoleModel, error) {
	var roles []*RoleModel
	err := impl.dbConnection.Model(&roles).Where("team = ?", teamName).Select()
	if err != nil {
		impl.Logger.Errorw("error in getting roles for team", "err", err, "teamName", teamName)
		return nil, err
	}
	return roles, nil
}

func (impl UserAuthRepositoryImpl) GetRolesForApp(appName string) ([]*RoleModel, error) {
	var roles []*RoleModel
	err := impl.dbConnection.Model(&roles).Where("(entity ='apps' and access_type='devtron-app') OR (entity ='jobs' and access_type='')").
		Where("entity_name = ?", appName).Select()
	if err != nil {
		impl.Logger.Errorw("error in getting roles for app", "err", err, "appName", appName)
		return nil, err
	}
	return roles, nil
}

func (impl UserAuthRepositoryImpl) GetRolesForChartGroup(chartGroupName string) ([]*RoleModel, error) {
	var roles []*RoleModel
	err := impl.dbConnection.Model(&roles).Where("entity = ?", bean2.CHART_GROUP_TYPE).
		Where("entity_name = ?", chartGroupName).Select()
	if err != nil {
		impl.Logger.Errorw("error in getting roles for chart group", "err", err, "chartGroupName", chartGroupName)
		return nil, err
	}
	return roles, nil
}

func (impl UserAuthRepositoryImpl) DeleteRole(role *RoleModel, tx *pg.Tx) error {
	err := tx.Delete(role)
	if err != nil {
		impl.Logger.Errorw("error in deleting role", "err", err, "role", role)
		return err
	}
	return nil
}

func (impl UserAuthRepositoryImpl) DeleteRolesByIds(roleIds []int, tx *pg.Tx) error {
	var models []RoleModel
	_, err := tx.Model(&models).Where("id in (?)", pg.In(roleIds)).Delete()
	if err != nil {
		impl.Logger.Errorw("error in deleting roles by roleIds", "err", err, "roles", roleIds)
		return err
	}
	return nil
}

func (impl UserAuthRepositoryImpl) GetRolesByUserIdAndEntityType(userId int32, entityType string) ([]*RoleModel, error) {
	var models []*RoleModel
	err := impl.dbConnection.Model(&models).
		Column("role_model.*").
		Join("INNER JOIN user_roles ur on ur.role_id=role_model.id").
		Where("role_model.entity = ?", entityType).
		Where("ur.user_id = ?", userId).Select()
	if err != nil {
		impl.Logger.Error(err)
		return models, err
	}
	return models, nil
}

func (impl UserAuthRepositoryImpl) GetRolesByEntityAccessTypeAndAction(entity, accessType, action string) ([]*RoleModel, error) {
	var models []*RoleModel
	var err error
	if accessType == "" {
		err = impl.dbConnection.Model(&models).Where("action = ?", action).
			Where("entity = ?", entity).Where("access_type is NULL").
			Select()
	} else {
		err = impl.dbConnection.Model(&models).Where("action = ?", action).
			Where("entity = ?", entity).Where("access_type = ?", accessType).
			Select()
	}
	if err != nil {
		impl.Logger.Error("err, GetRolesByEntityAccessTypeAndAction", "err", err, "entity", entity, "accessType", accessType, "action", action)
		return models, err
	}
	return models, nil
}

func (impl UserAuthRepositoryImpl) GetRolesForWorkflow(workflow, entityName string) ([]*RoleModel, error) {
	var roles []*RoleModel
	err := impl.dbConnection.Model(&roles).Where("workflow = ?", workflow).
		Where("entity_name = ?", entityName).
		Select()
	if err != nil {
		impl.Logger.Errorw("error in getting roles for team", "err", err, "workflow", workflow)
		return nil, err
	}
	return roles, nil
}

func (impl UserAuthRepositoryImpl) GetRoleForClusterEntity(cluster, namespace, group, kind, resource, action string) (RoleModel, error) {
	var model RoleModel
	query := "SELECT * FROM roles  WHERE entity = ? "
	var err error

	if len(cluster) > 0 {
		query += " and cluster='" + cluster + "' "
	} else {
		query += " and cluster IS NULL "
	}
	if len(namespace) > 0 {
		query += " and namespace='" + namespace + "' "
	} else {
		query += " and namespace IS NULL "
	}
	if len(group) > 0 {
		query += " and \"group\"='" + group + "' "
	} else {
		query += " and \"group\" IS NULL "
	}
	if len(kind) > 0 {
		query += " and kind='" + kind + "' "
	} else {
		query += " and kind IS NULL "
	}
	if len(resource) > 0 {
		query += " and resource='" + resource + "' "
	} else {
		query += " and resource IS NULL "
	}
	if len(action) > 0 {
		query += " and action='" + action + "' ;"
	} else {
		query += " and action IS NULL ;"
	}
	_, err = impl.dbConnection.Query(&model, query, bean.CLUSTER_ENTITIY)
	if err != nil {
		impl.Logger.Errorw("error in getting roles for clusterEntity", "err", err,
			bean2.CLUSTER, cluster, "namespace", namespace, "kind", kind, "group", group, "resource", resource)
		return model, err
	}
	return model, nil

}
func (impl UserAuthRepositoryImpl) GetRoleForJobsEntity(entity, team, app, env, act string, workflow string) (RoleModel, error) {
	var model RoleModel
	var err error
	if len(team) > 0 && len(act) > 0 {
		query := "SELECT role.* FROM roles role WHERE role.team = ? AND role.action=? AND role.entity=? "
		if len(env) == 0 {
			query = query + " AND role.environment is NULL"
		} else {
			query += "AND role.environment='" + env + "'"
		}
		if len(app) == 0 {
			query = query + " AND role.entity_name is NULL"
		} else {
			query += " AND role.entity_name='" + app + "'"
		}
		if len(workflow) == 0 {
			query = query + " AND role.workflow is NULL;"
		} else {
			query += " AND role.workflow='" + workflow + "';"
		}
		_, err = impl.dbConnection.Query(&model, query, team, act, entity)
	} else {
		return model, nil
	}
	return model, err
}

func (impl UserAuthRepositoryImpl) GetRoleForChartGroupEntity(entity, app, act, accessType string) (RoleModel, error) {
	var model RoleModel
	var err error
	if len(app) > 0 && act == "update" {
		query := "SELECT role.* FROM roles role WHERE role.entity = ? AND role.entity_name=? AND role.action=?"
		if len(accessType) == 0 {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		_, err = impl.dbConnection.Query(&model, query, entity, app, act)
	} else if app == "" {
		query := "SELECT role.* FROM roles role WHERE role.entity = ? AND role.action=?"
		if len(accessType) == 0 {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		_, err = impl.dbConnection.Query(&model, query, entity, act)
	}
	if err != nil {
		impl.Logger.Errorw("error in getting role for chart group entity", "err", err, "entity", entity, "app", app, "act", act, "accessType", accessType)
	}
	return model, err
}

func (impl UserAuthRepositoryImpl) GetRoleForOtherEntity(team, app, env, act, accessType string, oldValues, approver bool) (RoleModel, error) {
	var model RoleModel
	var err error
	if len(team) > 0 && len(app) > 0 && len(env) > 0 && len(act) > 0 {
		query := "SELECT role.* FROM roles role WHERE role.team = ? AND role.entity_name=? AND role.environment=? AND role.action=?"
		if oldValues {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		if approver {
			query += " and role.approver = true"
		} else {
			query += " and ( role.approver = false OR role.approver is null)"
		}

		_, err = impl.dbConnection.Query(&model, query, team, app, env, act)
	} else if len(team) > 0 && app == "" && len(env) > 0 && len(act) > 0 {

		query := "SELECT role.* FROM roles role WHERE role.team=? AND coalesce(role.entity_name,'')=? AND role.environment=? AND role.action=?"
		if oldValues {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		if approver {
			query += " and role.approver = true"
		} else {
			query += " and ( role.approver = false OR role.approver is null)"
		}
		_, err = impl.dbConnection.Query(&model, query, team, EMPTY, env, act)
	} else if len(team) > 0 && len(app) > 0 && env == "" && len(act) > 0 {
		//this is applicable for all environment of a team
		query := "SELECT role.* FROM roles role WHERE role.team = ? AND role.entity_name=? AND coalesce(role.environment,'')=? AND role.action=?"
		if oldValues {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		if approver {
			query += " and role.approver = true"
		} else {
			query += " and ( role.approver = false OR role.approver is null)"
		}

		_, err = impl.dbConnection.Query(&model, query, team, app, EMPTY, act)
	} else if len(team) > 0 && app == "" && env == "" && len(act) > 0 {
		//this is applicable for all environment of a team
		query := "SELECT role.* FROM roles role WHERE role.team = ? AND coalesce(role.entity_name,'')=? AND coalesce(role.environment,'')=? AND role.action=?"
		if oldValues {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		if approver {
			query += " and role.approver = true"
		} else {
			query += " and ( role.approver = false OR role.approver is null)"
		}

		_, err = impl.dbConnection.Query(&model, query, team, EMPTY, EMPTY, act)
	} else if team == "" && app == "" && env == "" && len(act) > 0 {
		//this is applicable for super admin, all env, all team, all app
		query := "SELECT role.* FROM roles role WHERE coalesce(role.team,'') = ? AND coalesce(role.entity_name,'')=? AND coalesce(role.environment,'')=? AND role.action=?"
		if len(accessType) == 0 {
			query = query + " and role.access_type is NULL"
		} else {
			query += " and role.access_type='" + accessType + "'"
		}
		if approver {
			query += " and role.approver = true"
		} else {
			query += " and ( role.approver = false OR role.approver is null)"
		}
		_, err = impl.dbConnection.Query(&model, query, EMPTY, EMPTY, EMPTY, act)
	} else if team == "" && app == "" && env == "" && act == "" {
		return model, nil
	} else {
		return model, nil
	}
	return model, err
}
