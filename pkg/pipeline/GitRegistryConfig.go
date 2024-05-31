/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package pipeline

import (
	"context"
	"github.com/devtron-labs/devtron/client/gitSensor"
	"github.com/devtron-labs/devtron/internal/constants"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/pipeline/types"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/juju/errors"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type GitRegistryConfig interface {
	Create(request *types.GitRegistry) (*types.GitRegistry, error)
	GetAll() ([]types.GitRegistry, error)
	FetchAllGitProviders() ([]types.GitRegistry, error)
	FetchOneGitProvider(id string) (*types.GitRegistry, error)
	Update(request *types.GitRegistry) (*types.GitRegistry, error)
	Delete(request *types.GitRegistry) error
}
type GitRegistryConfigImpl struct {
	logger              *zap.SugaredLogger
	gitProviderRepo     repository.GitProviderRepository
	GitSensorGrpcClient gitSensor.Client
}

func NewGitRegistryConfigImpl(logger *zap.SugaredLogger, gitProviderRepo repository.GitProviderRepository,
	GitSensorClient gitSensor.Client) *GitRegistryConfigImpl {
	return &GitRegistryConfigImpl{
		logger:              logger,
		gitProviderRepo:     gitProviderRepo,
		GitSensorGrpcClient: GitSensorClient,
	}
}

func (impl GitRegistryConfigImpl) Create(request *types.GitRegistry) (*types.GitRegistry, error) {
	impl.logger.Debugw("get repo create request", "req", request)
	exist, err := impl.gitProviderRepo.ProviderExists(request.Url)
	if err != nil {
		impl.logger.Errorw("error in fetch ", "url", request.Url, "err", err)
		err = &util.ApiError{
			//Code:            constants.GitProviderCreateFailed,
			InternalMessage: "git provider creation failed, error in fetching by url",
			UserMessage:     "git provider creation failed, error in fetching by url",
		}
		return nil, err
	}
	if exist {
		impl.logger.Warnw("repo already exists", "url", request.Url)
		err = &util.ApiError{
			Code:            constants.GitProviderCreateFailedAlreadyExists,
			InternalMessage: "git provider already exists",
			UserMessage:     "git provider already exists",
		}
		return nil, errors.NewAlreadyExists(err, request.Url)
	}
	provider := &repository.GitProvider{
		Name:          request.Name,
		Url:           request.Url,
		Id:            request.Id,
		AuthMode:      request.AuthMode,
		Password:      request.Password,
		Active:        request.Active,
		AccessToken:   request.AccessToken,
		SshPrivateKey: request.SshPrivateKey,
		UserName:      request.UserName,
		AuditLog:      sql.AuditLog{CreatedBy: request.UserId, CreatedOn: time.Now(), UpdatedOn: time.Now(), UpdatedBy: request.UserId},
		GitHostId:     request.GitHostId,
	}
	provider.SshPrivateKey = ModifySshPrivateKey(provider.SshPrivateKey, provider.AuthMode)
	err = impl.gitProviderRepo.Save(provider)
	if err != nil {
		impl.logger.Errorw("error in saving git repo config", "data", provider, "err", err)
		err = &util.ApiError{
			Code:            constants.GitProviderCreateFailedInDb,
			InternalMessage: "git provider failed to create in db",
			UserMessage:     "git provider failed to create in db",
		}
		return nil, err
	}
	err = impl.UpdateGitSensor(provider)
	if err != nil {
		impl.logger.Errorw("error in updating git repo config on sensor", "data", provider, "err", err)
		err = &util.ApiError{
			Code:            constants.GitProviderUpdateFailedInSync,
			InternalMessage: err.Error(),
			UserMessage:     "git provider failed to update in sync",
		}
		return nil, err
	}
	request.Id = provider.Id
	return request, nil
}

// get all active git providers
func (impl GitRegistryConfigImpl) GetAll() ([]types.GitRegistry, error) {
	impl.logger.Debug("get all provider request")
	providers, err := impl.gitProviderRepo.FindAllActiveForAutocomplete()
	if err != nil {
		impl.logger.Errorw("error in fetch all git providers", "err", err)
		return nil, err
	}
	var gitProviders []types.GitRegistry
	for _, provider := range providers {
		providerRes := types.GitRegistry{
			Id:        provider.Id,
			Name:      provider.Name,
			Url:       provider.Url,
			GitHostId: provider.GitHostId,
			AuthMode:  provider.AuthMode,
		}
		gitProviders = append(gitProviders, providerRes)
	}
	return gitProviders, err
}

func (impl GitRegistryConfigImpl) FetchAllGitProviders() ([]types.GitRegistry, error) {
	impl.logger.Debug("fetch all git providers from db")
	providers, err := impl.gitProviderRepo.FindAll()
	if err != nil {
		impl.logger.Errorw("error in fetch all git providers", "err", err)
		return nil, err
	}
	var gitProviders []types.GitRegistry
	for _, provider := range providers {
		providerRes := types.GitRegistry{
			Id:            provider.Id,
			Name:          provider.Name,
			Url:           provider.Url,
			UserName:      provider.UserName,
			Password:      "",
			AuthMode:      provider.AuthMode,
			AccessToken:   "",
			SshPrivateKey: "",
			Active:        provider.Active,
			UserId:        provider.CreatedBy,
			GitHostId:     provider.GitHostId,
		}
		gitProviders = append(gitProviders, providerRes)
	}
	return gitProviders, err
}

func (impl GitRegistryConfigImpl) FetchOneGitProvider(providerId string) (*types.GitRegistry, error) {
	impl.logger.Debug("fetch git provider by ID from db")
	provider, err := impl.gitProviderRepo.FindOne(providerId)
	if err != nil {
		impl.logger.Errorw("error in fetch all git providers", "err", err)
		return nil, err
	}

	providerRes := &types.GitRegistry{
		Id:            provider.Id,
		Name:          provider.Name,
		Url:           provider.Url,
		UserName:      provider.UserName,
		Password:      provider.Password,
		AuthMode:      provider.AuthMode,
		AccessToken:   provider.AccessToken,
		SshPrivateKey: provider.SshPrivateKey,
		Active:        provider.Active,
		UserId:        provider.CreatedBy,
		GitHostId:     provider.GitHostId,
	}

	return providerRes, err
}

func (impl GitRegistryConfigImpl) Update(request *types.GitRegistry) (*types.GitRegistry, error) {
	impl.logger.Debugw("get repo create request", "req", request)

	/*
		exist, err := impl.gitProviderRepo.ProviderExists(request.RedirectionUrl)
		if err != nil {
			impl.logger.Errorw("error in fetch ", "url", request.RedirectionUrl, "err", err)
			return nil, err
		}
		if exist {
			impl.logger.Infow("repo already exists", "url", request.RedirectionUrl)
			return nil, errors.NewAlreadyExists(err, request.RedirectionUrl)
		}
	*/

	providerId := strconv.Itoa(request.Id)
	existingProvider, err0 := impl.gitProviderRepo.FindOne(providerId)
	if err0 != nil {
		impl.logger.Errorw("No matching entry found for update.", "err", err0)
		err0 = &util.ApiError{
			Code:            constants.GitProviderUpdateProviderNotExists,
			InternalMessage: "git provider update failed, provider does not exist",
			UserMessage:     "git provider update failed, provider does not exist",
		}
		return nil, err0
	}
	if request.Password == "" {
		request.Password = existingProvider.Password
	}
	if request.SshPrivateKey == "" {
		request.SshPrivateKey = existingProvider.SshPrivateKey
	}
	if request.AccessToken == "" {
		request.AccessToken = existingProvider.AccessToken
	}
	provider := &repository.GitProvider{
		Name:          request.Name,
		Url:           request.Url,
		Id:            request.Id,
		AuthMode:      request.AuthMode,
		Password:      request.Password,
		Active:        request.Active,
		AccessToken:   request.AccessToken,
		SshPrivateKey: request.SshPrivateKey,
		UserName:      request.UserName,
		GitHostId:     request.GitHostId,
		AuditLog:      sql.AuditLog{CreatedBy: existingProvider.CreatedBy, CreatedOn: existingProvider.CreatedOn, UpdatedOn: time.Now(), UpdatedBy: request.UserId},
	}
	provider.SshPrivateKey = ModifySshPrivateKey(provider.SshPrivateKey, provider.AuthMode)
	err := impl.gitProviderRepo.Update(provider)
	if err != nil {
		impl.logger.Errorw("error in updating git repo config", "data", provider, "err", err)
		err = &util.ApiError{
			Code:            constants.GitProviderUpdateFailedInDb,
			InternalMessage: "git provider failed to update in db",
			UserMessage:     "git provider failed to update in db",
		}
		return nil, err
	}
	request.Id = provider.Id
	err = impl.UpdateGitSensor(provider)
	if err != nil {
		impl.logger.Errorw("error in updating git repo config on sensor", "data", provider, "err", err)
		err = &util.ApiError{
			Code:            constants.GitProviderUpdateFailedInSync,
			InternalMessage: err.Error(),
			UserMessage:     "git provider failed to update in sync",
		}
		return nil, err
	}
	return request, nil
}

func (impl GitRegistryConfigImpl) Delete(request *types.GitRegistry) error {
	providerId := strconv.Itoa(request.Id)
	gitProviderConfig, err := impl.gitProviderRepo.FindOne(providerId)
	if err != nil {
		impl.logger.Errorw("No matching entry found for delete.", "id", request.Id, "err", err)
		return err
	}
	deleteReq := gitProviderConfig
	deleteReq.UpdatedOn = time.Now()
	deleteReq.UpdatedBy = request.UserId
	err = impl.gitProviderRepo.MarkProviderDeleted(&deleteReq)
	if err != nil {
		impl.logger.Errorw("err in deleting git account", "id", request.Id, "err", err)
		return err
	}
	deleteReq.Active = false
	err = impl.UpdateGitSensor(&deleteReq)
	if err != nil {
		impl.logger.Errorw("error in updating git repo config on sensor after deleting", "deleteReq", deleteReq, "err", err)
		err = &util.ApiError{
			Code:            constants.GitProviderUpdateFailedInSync,
			InternalMessage: err.Error(),
			UserMessage:     "git provider failed to update in sync",
		}
		return err
	}
	return nil
}

func (impl GitRegistryConfigImpl) UpdateGitSensor(provider *repository.GitProvider) error {
	sensorGitProvider := &gitSensor.GitProvider{
		Id:            provider.Id,
		Url:           provider.Url,
		Name:          provider.Name,
		UserName:      provider.UserName,
		AccessToken:   provider.AccessToken,
		Password:      provider.Password,
		Active:        provider.Active,
		SshPrivateKey: provider.SshPrivateKey,
		AuthMode:      provider.AuthMode,
	}
	return impl.GitSensorGrpcClient.SaveGitProvider(context.Background(), sensorGitProvider)
}

// Modifying Ssh Private Key because Ssh key authentication requires a new-line at the end of string & there are chances that user skips sending \n
func ModifySshPrivateKey(sshPrivateKey string, authMode repository.AuthMode) string {
	if authMode == repository.AUTH_MODE_SSH {
		if !strings.HasSuffix(sshPrivateKey, "\n") {
			sshPrivateKey += "\n"
		}
	}
	return sshPrivateKey
}
