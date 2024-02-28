package policyGovernance

import "github.com/devtron-labs/devtron/pkg/globalPolicy/bean"

const NO_POLICY = "NA"

type PathVariablePolicyType string

const PathVariablePolicyTypeVariable string = "policyType"
const ImagePromotion PathVariablePolicyType = "artifact-promotion"

var ExistingPolicyTypes = []PathVariablePolicyType{ImagePromotion}
var PathPolicyTypeGlobalPolicyTypeMap = map[PathVariablePolicyType]bean.GlobalPolicyType{
	ImagePromotion: bean.IMAGE_PROMOTION_POLICY,
}

type AppEnvPolicyContainer struct {
	AppId      int    `json:"-"`
	EnvId      int    `json:"-"`
	PolicyId   int    `json:"-"`
	AppName    string `json:"appName"`
	EnvName    string `json:"envName"`
	PolicyName string `json:"policyName,omitempty"`
}

type AppEnvPolicyMappingsListFilter struct {
	PolicyType  bean.GlobalPolicyType `json:"-"`
	AppNames    []string              `json:"appNames"`
	EnvNames    []string              `json:"envNames"`
	PolicyNames []string              `json:"policyNames"`
	SortBy      string                `json:"sortBy,omitempty" validate:"omitempty,oneof=appName environmentName"`
	SortOrder   string                `json:"sortOrder,omitempty" validate:"omitempty,oneof=ASC DESC"`
	Offset      int                   `json:"offset,omitempty" validate:"omitempty,min=0"`
	Size        int                   `json:"size,omitempty" validate:"omitempty,min=0"`
}

type BulkPromotionPolicyApplyRequest struct {
	PolicyType              bean.GlobalPolicyType          `json:"-"`
	ApplicationEnvironments []AppEnvPolicyContainer        `json:"applicationEnvironments"`
	ApplyToPolicyName       string                         `json:"applyToPolicyName" validate:"min=3"`
	AppEnvPolicyListFilter  AppEnvPolicyMappingsListFilter `json:"appEnvPolicyListFilter" validate:"dive"`
}
