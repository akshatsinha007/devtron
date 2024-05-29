/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package pipeline

import (
	"github.com/devtron-labs/devtron/internal/constants"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/pipeline/types"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"time"
)

type GitHostConfig interface {
	GetAll() ([]types.GitHostRequest, error)
	GetById(id int) (*types.GitHostRequest, error)
	Create(request *types.GitHostRequest) (int, error)
}

type GitHostConfigImpl struct {
	logger           *zap.SugaredLogger
	gitHostRepo      repository.GitHostRepository
	attributeService attributes.AttributesService
}

func NewGitHostConfigImpl(gitHostRepo repository.GitHostRepository, logger *zap.SugaredLogger, attributeService attributes.AttributesService) *GitHostConfigImpl {
	return &GitHostConfigImpl{
		logger:           logger,
		gitHostRepo:      gitHostRepo,
		attributeService: attributeService,
	}
}

// get all git hosts
func (impl GitHostConfigImpl) GetAll() ([]types.GitHostRequest, error) {
	impl.logger.Debug("get all hosts request")
	hosts, err := impl.gitHostRepo.FindAll()
	if err != nil {
		impl.logger.Errorw("error in fetching all git hosts", "err", err)
		return nil, err
	}
	var gitHosts []types.GitHostRequest
	for _, host := range hosts {
		hostRes := types.GitHostRequest{
			Id:     host.Id,
			Name:   host.Name,
			Active: host.Active,
		}
		gitHosts = append(gitHosts, hostRes)
	}
	return gitHosts, err
}

// get git host by Id
func (impl GitHostConfigImpl) GetById(id int) (*types.GitHostRequest, error) {
	impl.logger.Debug("get hosts request for Id", id)
	host, err := impl.gitHostRepo.FindOneById(id)
	if err != nil {
		impl.logger.Errorw("error in fetching git host", "err", err)
		return nil, err
	}

	// get orchestrator host
	orchestratorHost, err := impl.attributeService.GetByKey("url")
	if err != nil {
		impl.logger.Errorw("error in fetching orchestrator host url from db", "err", err)
		return nil, err
	}

	var webhookUrlPrepend string
	if orchestratorHost == nil || len(orchestratorHost.Value) == 0 {
		webhookUrlPrepend = "{HOST_URL_PLACEHOLDER}"
	} else {
		webhookUrlPrepend = orchestratorHost.Value
	}
	webhookUrl := webhookUrlPrepend + host.WebhookUrl

	gitHost := &types.GitHostRequest{
		Id:              host.Id,
		Name:            host.Name,
		Active:          host.Active,
		WebhookUrl:      webhookUrl,
		WebhookSecret:   host.WebhookSecret,
		EventTypeHeader: host.EventTypeHeader,
		SecretHeader:    host.SecretHeader,
		SecretValidator: host.SecretValidator,
	}

	return gitHost, err
}

// Create in DB
func (impl GitHostConfigImpl) Create(request *types.GitHostRequest) (int, error) {
	impl.logger.Debugw("get git host create request", "req", request)
	exist, err := impl.gitHostRepo.Exists(request.Name)
	if err != nil {
		impl.logger.Errorw("error in fetching git host ", "name", request.Name, "err", err)
		err = &util.ApiError{
			InternalMessage: "git host creation failed, error in fetching by name",
			UserMessage:     "git host creation failed, error in fetching by name",
		}
		return 0, err
	}
	if exist {
		impl.logger.Warnw("git host already exists", "name", request.Name)
		err = &util.ApiError{
			Code:            constants.GitHostCreateFailedAlreadyExists,
			InternalMessage: "git host already exists",
			UserMessage:     "git host already exists",
		}
		return 0, errors.NewAlreadyExists(err, request.Name)
	}
	gitHost := &repository.GitHost{
		Name:     request.Name,
		Active:   request.Active,
		AuditLog: sql.AuditLog{CreatedBy: request.UserId, CreatedOn: time.Now(), UpdatedOn: time.Now(), UpdatedBy: request.UserId},
	}
	err = impl.gitHostRepo.Save(gitHost)
	if err != nil {
		impl.logger.Errorw("error in saving git host in db", "data", gitHost, "err", err)
		err = &util.ApiError{
			Code:            constants.GitHostCreateFailedInDb,
			InternalMessage: "git host failed to create in db",
			UserMessage:     "git host failed to create in db",
		}
		return 0, err
	}
	return gitHost.Id, nil
}
