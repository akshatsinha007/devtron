package bean

import (
	"encoding/json"
	"errors"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/pkg/cluster"
	repository1 "github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/globalPolicy/bean"
	bean2 "github.com/devtron-labs/devtron/pkg/policyGovernance/artifactPromotion/constants"
	"github.com/devtron-labs/devtron/util"
	"log"
	"time"
)

type ArtifactPromotionRequest struct {
	SourceName         string                   `json:"sourceName"`
	SourceType         bean2.SourceTypeStr      `json:"sourceType"`
	Action             string                   `json:"action"`
	PromotionRequestId int                      `json:"promotionRequestId"`
	ArtifactId         int                      `json:"artifactId"`
	AppName            string                   `json:"appName"`
	EnvironmentNames   []string                 `json:"destinationObjectNames"`
	UserId             int32                    `json:"-"`
	WorkflowId         int                      `json:"workflowId"`
	AppId              int                      `json:"appId"`
	EnvNameIdMap       map[string]int           `json:"-"`
	EnvIdNameMap       map[int]string           `json:"-"`
	SourcePipelineId   int                      `json:"-"`
	SourceCdPipeline   *pipelineConfig.Pipeline `json:"-"`
}

type ArtifactPromotionApprovalResponse struct {
	Id                      int
	PolicyId                int
	PolicyEvaluationAuditId int
	ArtifactId              int
	SourceType              bean2.SourceType
	SourcePipelineId        int
	DestinationPipelineId   int
	Status                  bean2.ArtifactPromotionRequestStatus
}

type PromotionApprovalMetaData struct {
	ApprovalRequestId    int                         `json:"approvalRequestId"`
	ApprovalRuntimeState string                      `json:"approvalRuntimeState"`
	ApprovalUsersData    []PromotionApprovalUserData `json:"approvalUsersData"`
	RequestedUserData    PromotionApprovalUserData   `json:"requestedUserData"`
	PromotedFrom         string                      `json:"promotedFrom"`
	PromotedFromType     string                      `json:"promotedFromType"`
	Policy               PromotionPolicy             `json:"policy" validate:"dive"`
}

func (promotionApprovalMetaData PromotionApprovalMetaData) GetApprovalUserIds() []int32 {
	approvalUserIds := make([]int32, len(promotionApprovalMetaData.ApprovalUsersData))
	for _, approvalUserData := range promotionApprovalMetaData.ApprovalUsersData {
		approvalUserIds = append(approvalUserIds, approvalUserData.UserId)
	}
	return approvalUserIds
}

func (promotionApprovalMetaData PromotionApprovalMetaData) GetRequestedUserId() int32 {
	return promotionApprovalMetaData.RequestedUserData.UserId
}

type PromotionPolicyMetaRequest struct {
	Search    string `json:"search"`
	SortBy    string `json:"sortBy" validate:"oneof=ASC DESC"`
	SortOrder string `json:"sortOrder" validate:"oneof=policyName approverCount"`
}

type PromotionApprovalUserData struct {
	UserId         int32     `json:"userId"`
	UserEmail      string    `json:"userEmail"`
	UserActionTime time.Time `json:"userActionTime"`
}

type EnvironmentPromotionMetaData struct {
	Name                       string                         `json:"name"` // environment name
	ApprovalCount              int                            `json:"approvalCount,omitempty"`
	PromotionValidationMessage string                         `json:"promotionEvaluationMessage"`
	PromotionValidationState   bean2.PromotionValidationState `json:"promotionEvaluationState"`
	PromotionPossible          bool                           `json:"promotionPossible"`
	IsVirtualEnvironment       bool                           `json:"isVirtualEnvironment"`
}

type EnvironmentApprovalMetadata struct {
	Name            string   `json:"name"` // environment name
	ApprovalAllowed bool     `json:"approvalAllowed"`
	Reasons         []string `json:"reason"`
}

type PromotionPolicy struct {
	Id                 int                      `json:"id" `
	Name               string                   `json:"name" devtronSearchableField:"name" validate:"min=3,max=50,global-entity-name"`
	Description        string                   `json:"description" validate:"max=300"`
	PolicyEvaluationId int                      `json:"-"`
	Conditions         []util.ResourceCondition `json:"conditions" validate:"omitempty,min=1"`
	ApprovalMetaData   ApprovalMetaData         `json:"approvalMetadata" validate:"dive"`
	IdentifierCount    *int                     `json:"identifierCount,omitempty"`
}

func (p *PromotionPolicy) CanBePromoted(approvalsGot int) bool {
	return approvalsGot >= p.ApprovalMetaData.ApprovalCount
}

func (p *PromotionPolicy) CanApprove(requestedUserId, imageBuiltByUserId, approvingUserId int32) bool {
	return (p.ApprovalMetaData.AllowRequesterFromApprove || requestedUserId == approvingUserId) && (p.ApprovalMetaData.AllowImageBuilderFromApprove || imageBuiltByUserId == approvingUserId)
}

func (policy *PromotionPolicy) ConvertToGlobalPolicyBaseModal(userId int32) (*bean.GlobalPolicyBaseModel, error) {
	jsonPolicyBytes, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}
	return &bean.GlobalPolicyBaseModel{
		PolicyOf:      bean.GLOBAL_POLICY_TYPE_IMAGE_PROMOTION_POLICY,
		Name:          policy.Name,
		Description:   policy.Description,
		Enabled:       true, // all the policies are by default enabled
		PolicyVersion: bean.GLOBAL_POLICY_VERSION_V1,
		Active:        true,
		UserId:        userId,
		JsonData:      string(jsonPolicyBytes),
	}, nil
}

func (policy *PromotionPolicy) ConvertToGlobalPolicyDataModel(userId int32) (*bean.GlobalPolicyDataModel, error) {
	baseModel, err := policy.ConvertToGlobalPolicyBaseModal(userId)
	if err != nil {
		return nil, err
	}
	searchKeys := util.GetSearchableFields(*policy)
	approvalSearchKeys := util.GetSearchableFields(policy.ApprovalMetaData)
	searchKeys = append(searchKeys, approvalSearchKeys...)
	return &bean.GlobalPolicyDataModel{
		GlobalPolicyBaseModel: *baseModel,
		SearchableFields:      searchKeys,
	}, nil
}

func (policy *PromotionPolicy) UpdateWithGlobalPolicy(rawPolicy *bean.GlobalPolicyBaseModel) error {
	err := json.Unmarshal([]byte(rawPolicy.JsonData), policy)
	if err != nil {
		log.Printf("error in unmarshalling global policy json into promotionPolicy object, globalPolicy:%v,  err:%v", rawPolicy, err)
		return errors.New("unable to extract promotion policies")
	}
	policy.Name = rawPolicy.Name
	policy.Id = rawPolicy.Id
	policy.Description = rawPolicy.Description
	return nil
}

type ApprovalMetaData struct {
	ApprovalCount                int  `json:"approverCount" devtronSearchableField:"approver_count" validate:"min=0"`
	AllowImageBuilderFromApprove bool `json:"allowImageBuilderFromApprove"`
	AllowRequesterFromApprove    bool `json:"allowRequesterFromApprove"`
	AllowApproverFromDeploy      bool `json:"allowApproverFromDeploy"`
}

type EnvironmentListingResponse struct {
	CiSource     CiSourceMetaData               `json:"ciSource"`
	Environments []EnvironmentPromotionMetaData `json:"environments"`
}

type CiSourceMetaData struct {
	Id   int                 `json:"id"`
	Name string              `json:"name"`
	Type bean2.SourceTypeStr `json:"type"`
}

// rename to appworkflow metadata
type WorkflowMetaData struct {
	WorkflowId   int
	AppName      string
	AppId        int
	EnvMap       map[string]repository1.Environment
	CiSourceData CiSourceMetaData
}

type PipelinesMetaData struct {
	activeAuthorisedPipelineIds            []int
	activeAuthorisedPipelineIdVsEnvNameMap map[int]string
	activeAuthorisedPipelineIdDaoMap       map[int]*pipelineConfig.Pipeline
	pipelineEnvIds                         []int
}

type SourceMetaData struct {
	id               int
	typeStr          bean2.SourceTypeStr
	name             string
	sourceWorkflowId int
	cdPipeline       *pipelineConfig.Pipeline
}

func (s *SourceMetaData) WithSourceWorkflowId(sourceWorkflowId int) *SourceMetaData {
	s.sourceWorkflowId = sourceWorkflowId
	return s
}

func (s *SourceMetaData) WithId(id int) *SourceMetaData {
	s.id = id
	return s
}

func (s *SourceMetaData) WithType(typeStr bean2.SourceTypeStr) *SourceMetaData {
	s.typeStr = typeStr
	return s
}
func (s *SourceMetaData) WithName(name string) *SourceMetaData {
	s.name = name
	return s
}
func (s *SourceMetaData) WithCdPipeline(cdPipeline *pipelineConfig.Pipeline) *SourceMetaData {
	s.cdPipeline = cdPipeline
	return s
}

func (s *SourceMetaData) GetCiSourceMeta() CiSourceMetaData {
	return CiSourceMetaData{
		Id:   s.id,
		Type: s.typeStr,
		Name: s.name,
	}
}

type RequestMetaData struct {
	activeEnvIdNameMap          map[int]string
	activeEnvNameIdMap          map[string]int
	userEnvNames                []string
	authorisedEnvMap            map[string]bool
	activeEnvironments          []*cluster.EnvironmentBean
	activeEnvironmentsMap       map[string]*cluster.EnvironmentBean
	destinationPipelineMetaData *PipelinesMetaData
	activeEnvIds                []int
	activeEnvNames              []string
	activeAuthorisedEnvNames    []string
	activeAuthorisedEnvIds      []int
	sourceMetaData              *SourceMetaData
	promotableEnvs              []string
	appId                       int

	ciArtifact *repository.CiArtifact
}

func (r *RequestMetaData) GetSourceMetaData() *SourceMetaData {
	return r.sourceMetaData
}

func (r *RequestMetaData) WithCiArtifact(ciArtifact *repository.CiArtifact) *RequestMetaData {
	r.ciArtifact = ciArtifact
	return r
}

func (r *RequestMetaData) WithAppId(appId int) *RequestMetaData {
	r.appId = appId
	return r
}

func (r *RequestMetaData) WithPromotableEnvMap(promotableEnvs []string) *RequestMetaData {
	r.promotableEnvs = promotableEnvs
	return r
}

func (r *RequestMetaData) SetSourceMetaData(sourceMetaData *SourceMetaData) {
	r.sourceMetaData = sourceMetaData
}

func (r *RequestMetaData) SetDestinationPipelineMetaData(activeAuthorisedPipelines []*pipelineConfig.Pipeline) {
	pipelineIds := make([]int, 0, len(activeAuthorisedPipelines))
	pipelineIdEnvNameMap := make(map[int]string)
	pipelineIdPipelineDaoMap := make(map[int]*pipelineConfig.Pipeline)
	pipelineEnvIds := make([]int, 0, len(activeAuthorisedPipelines))
	for _, pipeline := range activeAuthorisedPipelines {
		pipelineIds = append(pipelineIds, pipeline.Id)
		pipelineIdEnvNameMap[pipeline.Id] = pipeline.Environment.Name
		pipelineIdPipelineDaoMap[pipeline.Id] = pipeline
		pipelineEnvIds = append(pipelineEnvIds, pipeline.EnvironmentId)
	}

	pipelineMetaData := &PipelinesMetaData{
		activeAuthorisedPipelineIds:            pipelineIds,
		activeAuthorisedPipelineIdDaoMap:       pipelineIdPipelineDaoMap,
		activeAuthorisedPipelineIdVsEnvNameMap: pipelineIdEnvNameMap,
	}
	r.destinationPipelineMetaData = pipelineMetaData
}

func (r *RequestMetaData) SetActiveEnvironments(userGivenEnvNames []string, authorizedEnvironmentsMap map[string]bool, activeEnvs []*cluster.EnvironmentBean) {
	r.userEnvNames = userGivenEnvNames
	r.authorisedEnvMap = authorizedEnvironmentsMap
	r.activeEnvironments = activeEnvs
	activeEnvironmentsMap := make(map[string]*cluster.EnvironmentBean)
	activeEnvNames := make([]string, 0, len(r.activeEnvironments))
	authorisedEnvNames := make([]string, 0, len(r.authorisedEnvMap))
	activeAuthorisedEnvIds := make([]int, 0, len(r.authorisedEnvMap))
	activeEnvIds := make([]int, 0, len(r.activeEnvironments))
	activeEnvIdNameMap := make(map[int]string)
	activeEnvNameIdMap := make(map[string]int)
	for _, env := range r.activeEnvironments {
		activeEnvNames = append(activeEnvNames, env.Environment)
		activeEnvIds = append(activeEnvIds, env.Id)
		activeEnvironmentsMap[env.Environment] = env
		activeEnvIdNameMap[env.Id] = env.Environment
		activeEnvNameIdMap[env.Environment] = env.Id
		if r.authorisedEnvMap[env.Environment] {
			authorisedEnvNames = append(authorisedEnvNames, env.Environment)
			activeAuthorisedEnvIds = append(activeAuthorisedEnvIds, env.Id)
		}
	}

	r.activeEnvironmentsMap = activeEnvironmentsMap
	r.activeEnvNames = activeEnvNames
	r.activeAuthorisedEnvNames = authorisedEnvNames
	r.activeAuthorisedEnvIds = activeAuthorisedEnvIds
	r.activeEnvIds = activeAuthorisedEnvIds
	r.activeEnvIdNameMap = activeEnvIdNameMap
	r.activeEnvNameIdMap = activeEnvNameIdMap
	r.activeEnvNames = activeEnvNames
}

func (r *RequestMetaData) GetActiveEnvNames() []string {
	return r.activeEnvNames
}

func (r *RequestMetaData) GetActiveAuthorisedEnvNames() []string {
	return r.activeAuthorisedEnvNames
}

func (r *RequestMetaData) GetActiveAuthorisedEnvIds() []int {
	return r.activeAuthorisedEnvIds
}

func (r *RequestMetaData) GetPipelineById(id int) *pipelineConfig.Pipeline {
	return r.destinationPipelineMetaData.activeAuthorisedPipelineIdDaoMap[id]
}

func (r *RequestMetaData) GetWorkflowId() int {
	return r.sourceMetaData.sourceWorkflowId
}

func (r *RequestMetaData) GetSourceType() bean2.SourceTypeStr {
	return r.sourceMetaData.typeStr
}

func (r *RequestMetaData) GetSourcePipelineId() int {
	return r.sourceMetaData.id
}

func (r *RequestMetaData) GetSourceName() string {
	return r.sourceMetaData.name
}

func (r *RequestMetaData) GetSourceCdPipeline() *pipelineConfig.Pipeline {
	pipeline := *r.sourceMetaData.cdPipeline
	return &pipeline
}

func (r *RequestMetaData) GetActiveAuthorisedPipelineIds() []int {
	return r.destinationPipelineMetaData.activeAuthorisedPipelineIds
}

func (r *RequestMetaData) GetActiveAuthorisedPipelineIdEnvMap() map[int]string {
	return r.destinationPipelineMetaData.activeAuthorisedPipelineIdVsEnvNameMap
}

func (r *RequestMetaData) GetActiveAuthorisedPipelineDaoMap() map[int]*pipelineConfig.Pipeline {
	return r.destinationPipelineMetaData.activeAuthorisedPipelineIdDaoMap
}

func (r *RequestMetaData) GetActiveAuthorisedPipelineEnvIds() []int {
	return r.destinationPipelineMetaData.pipelineEnvIds
}

func (r *RequestMetaData) GetUserGivenEnvNames() []string {
	return r.userEnvNames
}

func (r *RequestMetaData) GetAuthorisedEnvMap() map[string]bool {
	return r.authorisedEnvMap
}

func (r *RequestMetaData) GetCiArtifact() *repository.CiArtifact {
	artifact := *r.ciArtifact
	return &artifact
}

func (r *RequestMetaData) GetActiveEnvironmentsMap() map[string]*cluster.EnvironmentBean {
	return r.activeEnvironmentsMap
}

func (r *RequestMetaData) GetAppId() int {
	return r.appId
}
