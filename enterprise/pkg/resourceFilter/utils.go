package resourceFilter

import (
	"encoding/json"
	"github.com/devtron-labs/devtron/pkg/devtronResource/bean"
)

const (
	NoResourceFiltersFound           = "no active resource filters found"
	AppAndEnvSelectorRequiredMessage = "both application and environment selectors are required"
	InvalidExpressions               = "one or more expressions are invalid"
	FilterNotFound                   = "filter not found"
)

// util methods
func getJsonStringFromResourceCondition(resourceConditions []ResourceCondition) (string, error) {

	jsonBytes, err := json.Marshal(resourceConditions)
	return string(jsonBytes), err
}

func getResourceConditionFromJsonString(conditionExpression string) ([]ResourceCondition, error) {
	res := make([]ResourceCondition, 0)
	err := json.Unmarshal([]byte(conditionExpression), &res)
	return res, err
}

func extractResourceConditions(resourceConditionJson string) ([]ResourceCondition, error) {
	var resourceConditions []ResourceCondition
	err := json.Unmarshal([]byte(resourceConditionJson), &resourceConditions)
	return resourceConditions, err
}

func convertToResponseBeans(resourceFilters []*ResourceFilter) ([]*FilterMetaDataBean, error) {
	var filterResponseBeans []*FilterMetaDataBean
	for _, resourceFilter := range resourceFilters {
		filterResponseBean, err := convertToFilterBean(resourceFilter)
		if err != nil {
			return filterResponseBeans, err
		}
		filterResponseBeans = append(filterResponseBeans, filterResponseBean.FilterMetaDataBean)
	}
	return filterResponseBeans, nil
}

func convertToFilterBean(resourceFilter *ResourceFilter) (*FilterRequestResponseBean, error) {
	var err error
	resourceConditions, err := extractResourceConditions(resourceFilter.ConditionExpression)
	if err != nil {
		return nil, err
	}
	filterResponseBean := &FilterRequestResponseBean{
		FilterMetaDataBean: &FilterMetaDataBean{
			Id:           resourceFilter.Id,
			TargetObject: resourceFilter.TargetObject,
			Description:  resourceFilter.Description,
			Name:         resourceFilter.Name,
			Conditions:   resourceConditions,
		},
	}
	return filterResponseBean, nil
}

func GetIdentifierKey(identifierType IdentifierType, searchableKeyNameIdMap map[bean.DevtronResourceSearchableKeyName]int) int {
	switch identifierType {
	case AppIdentifier:
		return searchableKeyNameIdMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_APP_ID]
	case ClusterIdentifier:
		return searchableKeyNameIdMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_CLUSTER_ID]
	case EnvironmentIdentifier:
		return searchableKeyNameIdMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_ENV_ID]
	case ProjectIdentifier:
		return searchableKeyNameIdMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_PROJECT_ID]
	default:
		//TODO: revisit
		return -1
	}
}

func getJsonStringFromFilterHistoryObjects(filterHistoryObjects []*FilterHistoryObject) (string, error) {
	jsonBytes, err := json.Marshal(filterHistoryObjects)
	return string(jsonBytes), err
}

func getFilterHistoryObjectsFromJsonString(jsonStr string) ([]*FilterHistoryObject, error) {
	filterHistoryObjects := make([]*FilterHistoryObject, 0)
	if jsonStr == "" {
		return filterHistoryObjects, nil
	}
	err := json.Unmarshal([]byte(jsonStr), &filterHistoryObjects)
	return filterHistoryObjects, err
}