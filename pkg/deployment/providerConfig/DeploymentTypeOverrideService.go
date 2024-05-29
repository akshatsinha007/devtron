/*
 * Copyright (c) 2024. Devtron Inc.
 */

package providerConfig

import (
	"fmt"
	util2 "github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/pkg/attributes"
	"github.com/devtron-labs/devtron/pkg/cluster"
	"github.com/devtron-labs/devtron/util"
	"go.uber.org/zap"
	"net/http"
)

type DeploymentTypeOverrideService interface {
	// ValidateAndOverrideDeploymentAppType : Set deployment application (helm/argo) types based on the enforcement configurations
	ValidateAndOverrideDeploymentAppType(deploymentType string, isGitOpsConfigured bool, environmentId int) (overrideDeploymentType string, err error)
}

type DeploymentTypeOverrideServiceImpl struct {
	logger             *zap.SugaredLogger
	deploymentConfig   *util.DeploymentServiceTypeConfig
	attributesService  attributes.AttributesService
	environmentService cluster.EnvironmentService
}

func NewDeploymentTypeOverrideServiceImpl(logger *zap.SugaredLogger,
	envVariables *util.EnvironmentVariables,
	attributesService attributes.AttributesService,
	environmentService cluster.EnvironmentService) *DeploymentTypeOverrideServiceImpl {
	return &DeploymentTypeOverrideServiceImpl{
		logger:             logger,
		deploymentConfig:   envVariables.DeploymentServiceTypeConfig,
		attributesService:  attributesService,
		environmentService: environmentService,
	}
}

func (impl *DeploymentTypeOverrideServiceImpl) ValidateAndOverrideDeploymentAppType(deploymentType string, isGitOpsConfigured bool, environmentId int) (overrideDeploymentType string, err error) {
	// initialise OverrideDeploymentType to the given DeploymentType
	overrideDeploymentType = deploymentType
	isVirtualEnvironment, err := impl.environmentService.IsVirtualEnvironmentById(environmentId)
	if err != nil {
		impl.logger.Errorw("error in fetching environment by id", "envId", environmentId, "err", err)
		return overrideDeploymentType, err
	}
	// if no deployment app type sent from user then we'll not validate
	deploymentTypeValidationConfig, err := impl.attributesService.GetDeploymentEnforcementConfig(environmentId)
	if err != nil {
		impl.logger.Errorw("error in getting enforcement config for deployment", "err", err)
		return overrideDeploymentType, err
	}
	// by default both deployment app type are allowed
	AllowedDeploymentAppTypes := map[string]bool{
		util2.PIPELINE_DEPLOYMENT_TYPE_ACD:  true,
		util2.PIPELINE_DEPLOYMENT_TYPE_HELM: true,
	}
	for k, v := range deploymentTypeValidationConfig {
		// rewriting allowed deployment types based on config provided by user
		AllowedDeploymentAppTypes[k] = v
	}
	if !impl.deploymentConfig.ExternallyManagedDeploymentType {
		if isVirtualEnvironment && len(overrideDeploymentType) != 0 && !util2.IsManifestPush(overrideDeploymentType) && !util2.IsManifestDownload(overrideDeploymentType) {
			impl.logger.Errorw("invalid deployment type for a virtual environment", "deploymentType", overrideDeploymentType)
			err = &util2.ApiError{
				HttpStatusCode:  http.StatusBadRequest,
				InternalMessage: fmt.Sprintf("Deployment type '%s' is not supported on virtual cluster", overrideDeploymentType),
				UserMessage:     fmt.Sprintf("Deployment type '%s' is not supported on virtual cluster", overrideDeploymentType),
			}
			return overrideDeploymentType, err
		}
		if !isVirtualEnvironment {
			if isGitOpsConfigured && AllowedDeploymentAppTypes[util2.PIPELINE_DEPLOYMENT_TYPE_ACD] {
				overrideDeploymentType = util2.PIPELINE_DEPLOYMENT_TYPE_ACD
			} else if AllowedDeploymentAppTypes[util2.PIPELINE_DEPLOYMENT_TYPE_HELM] {
				overrideDeploymentType = util2.PIPELINE_DEPLOYMENT_TYPE_HELM
			}
		}
	}
	if deploymentType == "" {
		if isVirtualEnvironment {
			overrideDeploymentType = util2.PIPELINE_DEPLOYMENT_TYPE_MANIFEST_DOWNLOAD
		} else if isGitOpsConfigured && AllowedDeploymentAppTypes[util2.PIPELINE_DEPLOYMENT_TYPE_ACD] {
			overrideDeploymentType = util2.PIPELINE_DEPLOYMENT_TYPE_ACD
		} else if AllowedDeploymentAppTypes[util2.PIPELINE_DEPLOYMENT_TYPE_HELM] {
			overrideDeploymentType = util2.PIPELINE_DEPLOYMENT_TYPE_HELM
		}
	}
	if err = impl.validateDeploymentAppType(overrideDeploymentType, deploymentTypeValidationConfig); err != nil {
		impl.logger.Errorw("validation error for the given deployment type", "deploymentType", deploymentType, "err", err)
		return overrideDeploymentType, err
	}
	if !isGitOpsConfigured && util2.IsAcdApp(overrideDeploymentType) {
		impl.logger.Errorw("GitOps not configured but selected as a deployment app type")
		err = &util2.ApiError{
			HttpStatusCode:  http.StatusBadRequest,
			InternalMessage: "GitOps integration is not installed/configured. Please install/configure GitOps or use helm option.",
			UserMessage:     "GitOps integration is not installed/configured. Please install/configure GitOps or use helm option.",
		}
		return overrideDeploymentType, err
	}
	return overrideDeploymentType, nil
}

func (impl *DeploymentTypeOverrideServiceImpl) validateDeploymentAppType(deploymentType string, deploymentConfig map[string]bool) error {

	// Config value doesn't exist in attribute table
	if deploymentConfig == nil {
		return nil
	}
	//Config value found to be true for ArgoCD and Helm both
	if allDeploymentConfigTrue(deploymentConfig) {
		return nil
	}
	//Case : {ArgoCD : false, Helm: true, HGF : true}
	if validDeploymentConfigReceived(deploymentConfig, deploymentType) {
		return nil
	}

	err := &util2.ApiError{
		HttpStatusCode:  http.StatusBadRequest,
		InternalMessage: "Received deployment app type doesn't match with the allowed deployment app type for this environment.",
		UserMessage:     "Received deployment app type doesn't match with the allowed deployment app type for this environment.",
	}
	return err
}

func allDeploymentConfigTrue(deploymentConfig map[string]bool) bool {
	for _, value := range deploymentConfig {
		if !value {
			return false
		}
	}
	return true
}

func validDeploymentConfigReceived(deploymentConfig map[string]bool, deploymentTypeSent string) bool {
	for key, value := range deploymentConfig {
		if value && key == deploymentTypeSent {
			return true
		}
	}
	return false
}
