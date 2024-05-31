/*
 * Copyright (c) 2020-2024. Devtron Inc.
 */

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	bean2 "github.com/devtron-labs/devtron/pkg/attributes/bean"
	"github.com/devtron-labs/devtron/pkg/module"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	pubsub "github.com/devtron-labs/common-lib/pubsub-lib"
	"github.com/devtron-labs/devtron/api/bean"
	"github.com/devtron-labs/devtron/client/gitSensor"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	util "github.com/devtron-labs/devtron/util/event"
	"go.uber.org/zap"
)

type EventClientConfig struct {
	DestinationURL string `env:"EVENT_URL" envDefault:"http://localhost:3000/notify"`
}

func GetEventClientConfig() (*EventClientConfig, error) {
	cfg := &EventClientConfig{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, errors.New("could not get event service url")
	}
	return cfg, err
}

type EventClient interface {
	WriteNotificationEvent(event Event) (bool, error)
	WriteNatsEvent(channel string, payload interface{}) error
	SendAnyEvent(event map[string]interface{}) (bool, error)
}

type Event struct {
	EventTypeId        int               `json:"eventTypeId"`
	EventName          string            `json:"eventName"`
	PipelineId         int               `json:"pipelineId"`
	PipelineType       string            `json:"pipelineType"`
	CorrelationId      string            `json:"correlationId"`
	Payload            *Payload          `json:"payload"`
	EventTime          string            `json:"eventTime"`
	TeamId             int               `json:"teamId"`
	AppId              int               `json:"appId"`
	EnvId              int               `json:"envId"`
	CdWorkflowType     bean.WorkflowType `json:"cdWorkflowType,omitempty"`
	CdWorkflowRunnerId int               `json:"cdWorkflowRunnerId"`
	CiWorkflowRunnerId int               `json:"ciWorkflowRunnerId"`
	CiArtifactId       int               `json:"ciArtifactId"`
	BaseUrl            string            `json:"baseUrl"`
	UserId             int               `json:"-"`
}

type Payload struct {
	AppName                          string               `json:"appName"`
	EnvName                          string               `json:"envName"`
	PipelineName                     string               `json:"pipelineName"`
	Source                           string               `json:"source"`
	DockerImageUrl                   string               `json:"dockerImageUrl"`
	TriggeredBy                      string               `json:"triggeredBy"`
	Stage                            string               `json:"stage"`
	DeploymentHistoryLink            string               `json:"deploymentHistoryLink"`
	AppDetailLink                    string               `json:"appDetailLink"`
	DownloadLink                     string               `json:"downloadLink"`
	BuildHistoryLink                 string               `json:"buildHistoryLink"`
	MaterialTriggerInfo              *MaterialTriggerInfo `json:"material"`
	ApprovedByEmail                  []string             `json:"approvedByEmail"`
	FailureReason                    string               `json:"failureReason"`
	Providers                        []*Provider          `json:"providers"`
	ImageTagNames                    []string             `json:"imageTagNames"`
	ImageComment                     string               `json:"imageComment"`
	ImageApprovalLink                string               `json:"imageApprovalLink"`
	ProtectConfigFileType            string               `json:"protectConfigFileType"`
	ProtectConfigFileName            string               `json:"protectConfigFileName"`
	ProtectConfigComment             string               `json:"protectConfigComment"`
	ProtectConfigLink                string               `json:"protectConfigLink"`
	ApprovalLink                     string               `json:"approvalLink"`
	TimeWindowComment                string               `json:"timeWindowComment"`
	ImageScanExecutionInfo           json.RawMessage      `json:"imageScanExecutionInfo"`
	ArtifactPromotionRequestViewLink string               `json:"artifactPromotionRequestViewLink"`
	ArtifactPromotionApprovalLink    string               `json:"artifactPromotionApprovalLink"`
	PromotionArtifactSource          string               `json:"promotionArtifactSource"`
	ScoopNotificationConfig          interface{}          `json:"scoopNotificationConfig"`
}

type CiPipelineMaterialResponse struct {
	Id              int                    `json:"id"`
	GitMaterialId   int                    `json:"gitMaterialId"`
	GitMaterialUrl  string                 `json:"gitMaterialUrl"`
	GitMaterialName string                 `json:"gitMaterialName"`
	Type            string                 `json:"type"`
	Value           string                 `json:"value"`
	Active          bool                   `json:"active"`
	History         []*gitSensor.GitCommit `json:"history,omitempty"`
	LastFetchTime   time.Time              `json:"lastFetchTime"`
	IsRepoError     bool                   `json:"isRepoError"`
	RepoErrorMsg    string                 `json:"repoErrorMsg"`
	IsBranchError   bool                   `json:"isBranchError"`
	BranchErrorMsg  string                 `json:"branchErrorMsg"`
	Url             string                 `json:"url"`
}

type MaterialTriggerInfo struct {
	GitTriggers map[int]pipelineConfig.GitCommit `json:"gitTriggers"`
	CiMaterials []CiPipelineMaterialResponse     `json:"ciMaterials"`
}

type EventRESTClientImpl struct {
	logger               *zap.SugaredLogger
	client               *http.Client
	config               *EventClientConfig
	pubsubClient         *pubsub.PubSubClientServiceImpl
	ciPipelineRepository pipelineConfig.CiPipelineRepository
	pipelineRepository   pipelineConfig.PipelineRepository
	attributesRepository repository.AttributesRepository
	moduleService        module.ModuleService
}

func NewEventRESTClientImpl(logger *zap.SugaredLogger, client *http.Client, config *EventClientConfig, pubsubClient *pubsub.PubSubClientServiceImpl,
	ciPipelineRepository pipelineConfig.CiPipelineRepository, pipelineRepository pipelineConfig.PipelineRepository,
	attributesRepository repository.AttributesRepository, moduleService module.ModuleService) *EventRESTClientImpl {
	return &EventRESTClientImpl{logger: logger, client: client, config: config, pubsubClient: pubsubClient,
		ciPipelineRepository: ciPipelineRepository, pipelineRepository: pipelineRepository,
		attributesRepository: attributesRepository, moduleService: moduleService}
}

func (impl *EventRESTClientImpl) buildFinalPayload(event Event, cdPipeline *pipelineConfig.Pipeline, ciPipeline *pipelineConfig.CiPipeline) *Payload {
	payload := event.Payload
	if payload == nil {
		payload = &Payload{}
	}
	if event.PipelineType == string(util.CD) {
		if cdPipeline != nil {
			payload.AppName = cdPipeline.App.AppName
			payload.EnvName = cdPipeline.Environment.Name
			payload.PipelineName = cdPipeline.Name
		}
		payload.Stage = string(event.CdWorkflowType)
		payload.DeploymentHistoryLink = fmt.Sprintf("/dashboard/app/%d/cd-details/%d/%d/%d/source-code", event.AppId, event.EnvId, event.PipelineId, event.CdWorkflowRunnerId)
		payload.AppDetailLink = fmt.Sprintf("/dashboard/app/%d/details/%d/pod", event.AppId, event.EnvId)
		if event.CdWorkflowType != bean.CD_WORKFLOW_TYPE_DEPLOY {
			payload.DownloadLink = fmt.Sprintf("/orchestrator/app/cd-pipeline/workflow/download/%d/%d/%d/%d", event.AppId, event.EnvId, event.PipelineId, event.CdWorkflowRunnerId)
		}
	} else if event.PipelineType == string(util.CI) {
		if ciPipeline != nil {
			payload.AppName = ciPipeline.App.AppName
			payload.PipelineName = ciPipeline.Name
		}
		payload.BuildHistoryLink = fmt.Sprintf("/dashboard/app/%d/ci-details/%d/%d/artifacts", event.AppId, event.PipelineId, event.CiWorkflowRunnerId)
	}
	return payload
}

func (impl *EventRESTClientImpl) WriteNotificationEvent(event Event) (bool, error) {
	// if notification integration is not installed then do not send the notification
	moduleInfo, err := impl.moduleService.GetModuleInfo(module.ModuleNameNotification)
	if err != nil {
		impl.logger.Errorw("error while getting notification module status", "err", err)
		return false, err
	}
	if moduleInfo.Status != module.ModuleStatusInstalled {
		impl.logger.Warnw("Notification module is not installed, hence skipping sending notification", "currentModuleStatus", moduleInfo.Status)
		return false, nil
	}

	var cdPipeline *pipelineConfig.Pipeline
	var ciPipeline *pipelineConfig.CiPipeline
	if event.PipelineId > 0 {
		if event.PipelineType == string(util.CD) {
			cdPipeline, err = impl.pipelineRepository.FindById(event.PipelineId)
			if err != nil {
				impl.logger.Errorw("error while fetching pipeline", "err", err)
				return false, err
			}
			if cdPipeline != nil {
				event.TeamId = cdPipeline.App.TeamId
			}
		} else if event.PipelineType == string(util.CI) {
			ciPipeline, err = impl.ciPipelineRepository.FindById(event.PipelineId)
			if err != nil {
				impl.logger.Errorw("error while fetching pipeline", "err", err)
				return false, err
			}
			if ciPipeline != nil {
				event.TeamId = ciPipeline.App.TeamId
			}
		}
	}

	payload := impl.buildFinalPayload(event, cdPipeline, ciPipeline)
	event.Payload = payload

	isPreStageExist := false
	isPostStageExist := false
	if cdPipeline != nil && len(cdPipeline.PreStageConfig) > 0 {
		isPreStageExist = true
	}
	if cdPipeline != nil && len(cdPipeline.PostStageConfig) > 0 {
		isPostStageExist = true
	}

	attribute, err := impl.attributesRepository.FindByKey(bean2.HostUrlKey)
	if err != nil {
		impl.logger.Errorw("there is host url configured", "ci pipeline", ciPipeline)
		return false, err
	}
	if attribute != nil {
		event.BaseUrl = attribute.Value
	}
	if event.CdWorkflowType == "" {
		_, err = impl.sendEvent(event)
	} else if event.CdWorkflowType == bean.CD_WORKFLOW_TYPE_PRE {
		if event.EventTypeId == int(util.Success) {
			impl.logger.Debug("skip - will send from deployment or post stage")
		} else {
			_, err = impl.sendEvent(event)
		}
	} else if event.CdWorkflowType == bean.CD_WORKFLOW_TYPE_DEPLOY {
		if isPreStageExist && event.EventTypeId == int(util.Trigger) {
			impl.logger.Debug("skip - already sent from pre stage")
		} else if isPostStageExist && event.EventTypeId == int(util.Success) {
			impl.logger.Debug("skip - will send from post stage")
		} else {
			_, err = impl.sendEvent(event)
		}
	} else if event.CdWorkflowType == bean.CD_WORKFLOW_TYPE_POST {
		if event.EventTypeId == int(util.Trigger) {
			impl.logger.Debug("skip - already sent from pre or deployment stage")
		} else {
			_, err = impl.sendEvent(event)
		}
	}
	return true, err
}

// do not call this method if notification module is not installed
func (impl *EventRESTClientImpl) sendEvent(event Event) (bool, error) {
	impl.logger.Debugw("event before send", "event", event)
	body, err := json.Marshal(event)
	if err != nil {
		impl.logger.Errorw("error while marshaling event request ", "err", err)
		return false, err
	}
	var reqBody = []byte(body)
	req, err := http.NewRequest(http.MethodPost, impl.config.DestinationURL, bytes.NewBuffer(reqBody))
	if err != nil {
		impl.logger.Errorw("error while writing event", "err", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := impl.client.Do(req)
	if err != nil {
		impl.logger.Errorw("error while UpdateJiraTransition request ", "err", err)
		return false, err
	}
	defer resp.Body.Close()
	impl.logger.Debugw("event completed", "event resp", resp)
	return true, err
}

func (impl *EventRESTClientImpl) WriteNatsEvent(topic string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = impl.pubsubClient.Publish(topic, string(body))
	return err
}

// do not call this method if notification module is not installed
func (impl *EventRESTClientImpl) SendAnyEvent(event map[string]interface{}) (bool, error) {
	impl.logger.Debugw("event before send", "event", event)
	body, err := json.Marshal(event)
	if err != nil {
		impl.logger.Errorw("error while marshaling event request ", "err", err)
		return false, err
	}
	var reqBody = []byte(body)
	req, err := http.NewRequest(http.MethodPost, impl.config.DestinationURL, bytes.NewBuffer(reqBody))
	if err != nil {
		impl.logger.Errorw("error while writing event", "err", err)
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := impl.client.Do(req)
	if err != nil {
		impl.logger.Errorw("error while UpdateJiraTransition request ", "err", err)
		return false, err
	}
	defer resp.Body.Close()
	impl.logger.Debugw("event completed", "event resp", resp)
	return true, err
}
