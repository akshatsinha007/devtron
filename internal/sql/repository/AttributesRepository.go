/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package repository

import (
	"github.com/devtron-labs/devtron/pkg/attributes/bean"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/go-pg/pg"
)

type Attributes struct {
	tableName struct{} `sql:"attributes" pg:",discard_unknown_columns"`
	Id        int      `sql:"id,pk"`
	Key       string   `sql:"key,notnull"`
	Value     string   `sql:"value,notnull"`
	Active    bool     `sql:"active, notnull"`
	sql.AuditLog
}

type AttributesRepository interface {
	Save(model *Attributes, tx *pg.Tx) (*Attributes, error)
	Update(model *Attributes, tx *pg.Tx) error
	FindByKey(key string) (*Attributes, error)
	FindById(id int) (*Attributes, error)
	FindActiveList() ([]*Attributes, error)
	GetConnection() (dbConnection *pg.DB)
}

// TODO:caching because of high traffic calls clean this after proper fix
var attributeForEnforcedDeploymentTypeConfig *Attributes

func invalidateEnforcedDeploymentCache(model *Attributes) {
	if model.Key != bean.ENFORCE_DEPLOYMENT_TYPE_CONFIG {
		return
	}
	attributeForEnforcedDeploymentTypeConfig = nil
}

type AttributesRepositoryImpl struct {
	dbConnection *pg.DB
}

func NewAttributesRepositoryImpl(dbConnection *pg.DB) *AttributesRepositoryImpl {
	return &AttributesRepositoryImpl{dbConnection: dbConnection}
}

func (impl *AttributesRepositoryImpl) GetConnection() (dbConnection *pg.DB) {
	return impl.dbConnection
}

func (repo AttributesRepositoryImpl) Save(model *Attributes, tx *pg.Tx) (*Attributes, error) {
	err := tx.Insert(model)
	if err != nil {
		return model, err
	}
	// reset cached data
	invalidateEnforcedDeploymentCache(model)
	return model, nil
}

func (repo AttributesRepositoryImpl) Update(model *Attributes, tx *pg.Tx) error {
	err := tx.Update(model)
	if err != nil {
		return err
	}
	// reset cached data
	invalidateEnforcedDeploymentCache(model)
	return nil
}

func (repo AttributesRepositoryImpl) FindByKey(key string) (*Attributes, error) {
	// use cached data if existing
	if key == bean.ENFORCE_DEPLOYMENT_TYPE_CONFIG &&
		attributeForEnforcedDeploymentTypeConfig != nil {
		return attributeForEnforcedDeploymentTypeConfig, nil
	}
	model := &Attributes{}
	err := repo.dbConnection.
		Model(model).
		Where("key = ?", key).
		Where("active = ?", true).
		Select()
	if err != nil {
		return model, err
	}
	// update cached data if not existing
	if key == bean.ENFORCE_DEPLOYMENT_TYPE_CONFIG {
		attributeForEnforcedDeploymentTypeConfig = model
	}
	return model, nil
}

func (repo AttributesRepositoryImpl) FindById(id int) (*Attributes, error) {
	model := &Attributes{}
	err := repo.dbConnection.Model(model).Where("id = ?", id).Where("active = ?", true).
		Select()
	return model, err
}

func (repo AttributesRepositoryImpl) FindActiveList() ([]*Attributes, error) {
	var models []*Attributes
	err := repo.dbConnection.Model(&models).Where("active = ?", true).
		Select()
	return models, err
}
