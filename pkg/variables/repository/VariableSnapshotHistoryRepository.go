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

package repository

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

type VariableSnapshotHistoryRepository interface {
	SaveVariableSnapshots(variableSnapshotHistories []*VariableSnapshotHistory) error
	CheckIfVariableSnapshotExists(historyReference HistoryReference) (bool, error)
	GetVariableSnapshots(historyReferences []HistoryReference) ([]*VariableSnapshotHistory, error)
}

func (impl VariableSnapshotHistoryRepositoryImpl) SaveVariableSnapshots(variableSnapshotHistories []*VariableSnapshotHistory) error {
	err := impl.dbConnection.Insert(&variableSnapshotHistories)
	if err != nil {
		impl.logger.Errorw("err in saving variable snapshot history", "err", err)
		return err
	}
	return nil
}

func (impl VariableSnapshotHistoryRepositoryImpl) CheckIfVariableSnapshotExists(historyReference HistoryReference) (bool, error) {
	var variableSnapshotHistory VariableSnapshotHistory
	exists, err := impl.dbConnection.Model(&variableSnapshotHistory).
		Where("history_reference_id = ?", historyReference.HistoryReferenceId).
		Where("history_reference_type = ?", historyReference.HistoryReferenceType).
		Exists()
	return exists, err
}

func (impl VariableSnapshotHistoryRepositoryImpl) GetVariableSnapshots(historyReferences []HistoryReference) ([]*VariableSnapshotHistory, error) {
	variableSnapshotHistories := make([]*VariableSnapshotHistory, 0)

	err := impl.dbConnection.Model(&variableSnapshotHistories).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			for _, historyReference := range historyReferences {
				q = q.WhereOr("history_reference_id = ? AND history_reference_type = ?", historyReference.HistoryReferenceId, historyReference.HistoryReferenceType)
			}
			return q, nil
		}).Select()
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("err in getting variables for entities", "err", err)
		return nil, err
	}

	return variableSnapshotHistories, nil
}

func NewVariableSnapshotHistoryRepository(logger *zap.SugaredLogger, dbConnection *pg.DB) *VariableSnapshotHistoryRepositoryImpl {
	return &VariableSnapshotHistoryRepositoryImpl{
		logger:       logger,
		dbConnection: dbConnection,
	}
}

type VariableSnapshotHistoryRepositoryImpl struct {
	logger       *zap.SugaredLogger
	dbConnection *pg.DB
}
