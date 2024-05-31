/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package repository

import (
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type DeploymentGroupAppRepository interface {
	Create(model *DeploymentGroupApp) (*DeploymentGroupApp, error)
	GetById(id int) (*DeploymentGroupApp, error)
	GetAll() ([]*DeploymentGroupApp, error)
	Update(model *DeploymentGroupApp) (*DeploymentGroupApp, error)
	Delete(model *DeploymentGroupApp) error
	GetByDeploymentGroup(deploymentGroupId int) ([]*DeploymentGroupApp, error)
}

type DeploymentGroupAppRepositoryImpl struct {
	dbConnection *pg.DB
	Logger       *zap.SugaredLogger
}

func NewDeploymentGroupAppRepositoryImpl(Logger *zap.SugaredLogger, dbConnection *pg.DB) *DeploymentGroupAppRepositoryImpl {
	return &DeploymentGroupAppRepositoryImpl{dbConnection: dbConnection, Logger: Logger}
}

type DeploymentGroupApp struct {
	TableName         struct{} `sql:"deployment_group_app" pg:",discard_unknown_columns"`
	Id                int      `sql:"id,pk"`
	DeploymentGroupId int      `sql:"deployment_group_id"`
	AppId             int      `sql:"app_id"`
	Active            bool     `sql:"active,notnull"`
	sql.AuditLog
}

func (impl DeploymentGroupAppRepositoryImpl) Create(model *DeploymentGroupApp) (*DeploymentGroupApp, error) {
	err := impl.dbConnection.Insert(model)
	if err != nil {
		impl.Logger.Error(err)
		return model, err
	}
	return model, nil
}

func (impl DeploymentGroupAppRepositoryImpl) GetById(id int) (*DeploymentGroupApp, error) {
	var model DeploymentGroupApp
	err := impl.dbConnection.Model(&model).Where("id = ?", id).Select()
	return &model, err
}

func (impl DeploymentGroupAppRepositoryImpl) GetAll() ([]*DeploymentGroupApp, error) {
	var models []*DeploymentGroupApp
	err := impl.dbConnection.Model(&models).Select()
	return models, err
}

func (impl DeploymentGroupAppRepositoryImpl) Update(model *DeploymentGroupApp) (*DeploymentGroupApp, error) {
	err := impl.dbConnection.Update(model)
	if err != nil {
		impl.Logger.Error(err)
		return model, err
	}
	return model, nil
}

func (impl DeploymentGroupAppRepositoryImpl) Delete(model *DeploymentGroupApp) error {
	err := impl.dbConnection.Delete(model)
	if err != nil {
		impl.Logger.Error(err)
		return err
	}
	return nil
}

func (impl DeploymentGroupAppRepositoryImpl) GetByDeploymentGroup(deploymentGroupId int) ([]*DeploymentGroupApp, error) {
	var models []*DeploymentGroupApp
	err := impl.dbConnection.Model(&models).
		Where("deployment_group_id = ?", deploymentGroupId).
		Select()
	return models, err
}
