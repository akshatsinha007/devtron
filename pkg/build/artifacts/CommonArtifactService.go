/*
 * Copyright (c) 2024. Devtron Inc.
 */

package artifacts

import (
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/pkg/sql"
	"go.uber.org/zap"
	"time"
)

//this service is created with the thought of as a common interaction point for artifact related operations required
//for build and deployment. Can be updated/removed as per evolution patterns.
//ques can be asked - why this is not outside /build/artifact? Because the place used for storing artifacts is still known as ci_artifact and a lot of current impl
//has duplicated logic, until it is extracted and the design pattern is made easily justifiable, have placed this inside /build/artifact.

type CommonArtifactService interface {
	SavePluginArtifacts(ciArtifact *repository.CiArtifact, pluginArtifactsDetail map[string][]string,
		pipelineId int, stage string, triggeredBy int32) ([]*repository.CiArtifact, error)
}

type CommonArtifactServiceImpl struct {
	logger               *zap.SugaredLogger
	ciArtifactRepository repository.CiArtifactRepository
}

func NewCommonArtifactServiceImpl(logger *zap.SugaredLogger,
	ciArtifactRepository repository.CiArtifactRepository) *CommonArtifactServiceImpl {
	return &CommonArtifactServiceImpl{
		logger:               logger,
		ciArtifactRepository: ciArtifactRepository,
	}
}

func (impl *CommonArtifactServiceImpl) SavePluginArtifacts(ciArtifact *repository.CiArtifact, pluginArtifactsDetail map[string][]string,
	pipelineId int, stage string, triggeredBy int32) ([]*repository.CiArtifact, error) {
	saveArtifacts, err := impl.ciArtifactRepository.GetArtifactsByDataSourceAndComponentId(stage, pipelineId)
	if err != nil {
		return nil, err
	}
	PipelineArtifacts := make(map[string]bool)
	for _, artifact := range saveArtifacts {
		PipelineArtifacts[artifact.Image] = true
	}
	var parentCiArtifactId int
	if ciArtifact.ParentCiArtifact > 0 {
		parentCiArtifactId = ciArtifact.ParentCiArtifact
	} else {
		parentCiArtifactId = ciArtifact.Id
	}
	var CDArtifacts []*repository.CiArtifact
	for registry, artifacts := range pluginArtifactsDetail {
		// artifacts are list of images
		for _, artifact := range artifacts {
			_, artifactAlreadySaved := PipelineArtifacts[artifact]
			if artifactAlreadySaved {
				continue
			}
			pluginArtifact := &repository.CiArtifact{
				Image:                 artifact,
				ImageDigest:           ciArtifact.ImageDigest,
				MaterialInfo:          ciArtifact.MaterialInfo,
				DataSource:            stage,
				ComponentId:           pipelineId,
				CredentialsSourceType: repository.GLOBAL_CONTAINER_REGISTRY,
				CredentialSourceValue: registry,
				AuditLog: sql.AuditLog{
					CreatedOn: time.Now(),
					CreatedBy: triggeredBy,
					UpdatedOn: time.Now(),
					UpdatedBy: triggeredBy,
				},
				ParentCiArtifact: parentCiArtifactId,
			}
			CDArtifacts = append(CDArtifacts, pluginArtifact)
		}
	}
	_, err = impl.ciArtifactRepository.SaveAll(CDArtifacts)
	if err != nil {
		impl.logger.Errorw("Error in saving artifacts metadata generated by plugin")
		return CDArtifacts, err
	}
	return CDArtifacts, nil
}
