package imageDigestPolicy

import (
	"fmt"
	"github.com/devtron-labs/devtron/pkg/cluster/repository"
	"github.com/devtron-labs/devtron/pkg/devtronResource"
	"github.com/devtron-labs/devtron/pkg/devtronResource/bean"
	"github.com/devtron-labs/devtron/pkg/resourceQualifiers"
	"github.com/devtron-labs/devtron/pkg/sql"
	"github.com/go-pg/pg"
	"go.uber.org/zap"
	"time"
)

type ImageDigestQualifierMappingService interface {
	//CreateOrDeletePolicyForPipeline created policy for enforcing pull using digest at pipeline level
	CreateOrDeletePolicyForPipeline(pipelineId int, isImageDigestEnforcedForPipeline bool, UserId int32) error

	//IsPolicyConfiguredForPipeline returns true if pipeline or env or cluster has image digest policy enabled
	IsPolicyConfiguredForPipeline(pipelineId int) (bool, error)

	//CreateOrUpdatePolicyForCluster creates or updates image digest qualifier mapping for given cluster and environments
	CreateOrUpdatePolicyForCluster(policyRequest *PolicyRequest) (*PolicyRequest, error)

	//IsPolicyConfiguredAtGlobalLevel for env or cluster or for all clusters
	IsPolicyConfiguredAtGlobalLevel(envId int, clusterId int) (bool, error)

	//GetAllConfiguredPolicies get all cluster and environment configured for pull using image digest
	GetAllConfiguredPolicies() (*PolicyRequest, error)
}

type ImageDigestQualifierMappingServiceImpl struct {
	logger                       *zap.SugaredLogger
	qualifierMappingService      resourceQualifiers.QualifierMappingService
	devtronResourceSearchableKey devtronResource.DevtronResourceSearchableKeyService
	environmentRepository        repository.EnvironmentRepository
}

func NewImageDigestQualifierMappingServiceImpl(
	logger *zap.SugaredLogger,
	qualifierMappingService resourceQualifiers.QualifierMappingService,
	devtronResourceSearchableKey devtronResource.DevtronResourceSearchableKeyService,
	environmentRepository repository.EnvironmentRepository,
) *ImageDigestQualifierMappingServiceImpl {
	return &ImageDigestQualifierMappingServiceImpl{
		logger:                       logger,
		qualifierMappingService:      qualifierMappingService,
		devtronResourceSearchableKey: devtronResourceSearchableKey,
		environmentRepository:        environmentRepository,
	}
}

func (impl ImageDigestQualifierMappingServiceImpl) CreateOrDeletePolicyForPipeline(pipelineId int, isImageDigestEnforcedInRequest bool, UserId int32) error {

	devtronResourceSearchableKeyMap := impl.devtronResourceSearchableKey.GetAllSearchableKeyNameIdMap()

	qualifierMappings, err := impl.getQualifierMappingForPipeline(pipelineId)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in fetching qualifier mappings for resourceType: imageDigest by pipelineId", "pipelineId", pipelineId)
		return err
	}

	if len(qualifierMappings) == 0 && isImageDigestEnforcedInRequest {

		qualifierMapping := &resourceQualifiers.QualifierMapping{
			ResourceId:            resourceQualifiers.ImageDigestResourceId,
			ResourceType:          resourceQualifiers.ImageDigest,
			QualifierId:           int(resourceQualifiers.PIPELINE_QUALIFIER),
			IdentifierKey:         devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_PIPELINE_ID],
			IdentifierValueInt:    pipelineId,
			Active:                true,
			IdentifierValueString: fmt.Sprintf("%d", pipelineId),
			AuditLog: sql.AuditLog{
				CreatedOn: time.Now(),
				CreatedBy: UserId,
				UpdatedOn: time.Now(),
				UpdatedBy: UserId,
			},
		}

		dbConnection := impl.qualifierMappingService.GetDbConnection()
		tx, _ := dbConnection.Begin()
		_, err := impl.qualifierMappingService.CreateQualifierMappings([]*resourceQualifiers.QualifierMapping{qualifierMapping}, tx)
		if err != nil {
			impl.logger.Errorw("error in creating image digest qualifier mapping for pipeline", "err", err)
			return err
		}
		_ = tx.Commit()

	} else if !isImageDigestEnforcedInRequest && len(qualifierMappings) > 0 {

		dbConnection := impl.qualifierMappingService.GetDbConnection()
		tx, _ := dbConnection.Begin()
		auditLog := sql.AuditLog{
			CreatedOn: time.Now(),
			CreatedBy: UserId,
			UpdatedOn: time.Now(),
			UpdatedBy: UserId,
		}
		err := impl.qualifierMappingService.DeleteAllQualifierMappingsByIdentifierKeyAndValue(devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_PIPELINE_ID], pipelineId, auditLog, tx)
		if err != nil {
			impl.logger.Errorw("error in deleting image digest policy for pipeline", "err", err, "pipeline id", pipelineId)
			return err
		}
		_ = tx.Commit()

	}
	return nil
}

func (impl ImageDigestQualifierMappingServiceImpl) IsPolicyConfiguredForPipeline(pipelineId int) (bool, error) {
	qualifierMappings, err := impl.getQualifierMappingForPipeline(pipelineId)
	if err != nil && err != pg.ErrNoRows {
		impl.logger.Errorw("error in fetching qualifier mappings for resourceType: imageDigest by pipelineId", "pipelineId", pipelineId)
		return false, err
	}
	if err == pg.ErrNoRows || len(qualifierMappings) == 0 {
		return false, nil
	}
	return true, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) getQualifierMappingForPipeline(pipelineId int) ([]*resourceQualifiers.QualifierMapping, error) {
	scope := &resourceQualifiers.Scope{PipelineId: pipelineId}
	resourceIds := []int{resourceQualifiers.ImageDigestResourceId}
	qualifierMappings, err := impl.qualifierMappingService.GetQualifierMappings(resourceQualifiers.ImageDigest, scope, resourceIds)
	if err != nil && err != pg.ErrNoRows {
		return qualifierMappings, err
	}
	return qualifierMappings, err
}

func (impl ImageDigestQualifierMappingServiceImpl) GetAllConfiguredPolicies() (*PolicyRequest, error) {

	imageDigestQualifierMappings, err := impl.qualifierMappingService.GetQualifierMappingsByResourceType(resourceQualifiers.ImageDigest)
	if err != nil {
		impl.logger.Errorw("error in fetching qualifier mappings by resourceType", "resourceType: ", resourceQualifiers.ImageDigest)
		return nil, err
	}

	imageDigestPolicies := &PolicyRequest{
		ClusterDetails: make([]*ClusterDetail, 0),
		AllClusters:    false,
	}

	for _, qualifierMapping := range imageDigestQualifierMappings {
		if qualifierMapping.QualifierId == int(resourceQualifiers.GLOBAL_QUALIFIER) {
			imageDigestPolicies.AllClusters = true
			break
		}
	}

	if imageDigestPolicies.AllClusters {
		return imageDigestPolicies, nil
	}

	ClusterIdToEnvIdsMapping, err := impl.getClusterIdToEnvIdsMapping(imageDigestQualifierMappings)
	if err != nil {
		impl.logger.Errorw("error in getting cluster id to envIds map", "err", err)
		return nil, err
	}

	for clusterId, envIds := range ClusterIdToEnvIdsMapping {
		clusterDetail := &ClusterDetail{
			ClusterId: clusterId,
		}
		if len(envIds) == 0 {
			clusterDetail.PolicyType = ALL_ENVIRONMENTS
		} else {
			clusterDetail.Environments = envIds
			clusterDetail.PolicyType = SELECTED_ENVIRONMENTS
		}
		imageDigestPolicies.ClusterDetails = append(imageDigestPolicies.ClusterDetails, clusterDetail)
	}

	return imageDigestPolicies, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) getClusterIdToEnvIdsMapping(imageDigestQualifierMappings []*resourceQualifiers.QualifierMapping) (map[ClusterId][]EnvironmentId, error) {
	ClusterIdToEnvIdsMapping := make(map[ClusterId][]EnvironmentId)
	EnvToClusterMapping, err := impl.getEnvToClusterMapping(imageDigestQualifierMappings) // map[envId]clusterId
	if err != nil {
		impl.logger.Errorw("error in fetching environments to cluster maping", "err", err)
		return nil, err
	}
	for envId, clusterId := range EnvToClusterMapping {
		if _, ok := ClusterIdToEnvIdsMapping[clusterId]; !ok {
			ClusterIdToEnvIdsMapping[clusterId] = make([]int, 0)
		} else {
			ClusterIdToEnvIdsMapping[clusterId] = append(ClusterIdToEnvIdsMapping[clusterId], envId)
		}
	}
	return ClusterIdToEnvIdsMapping, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) getEnvToClusterMapping(imageDigestQualifierMappings []*resourceQualifiers.QualifierMapping) (map[EnvironmentId]ClusterId, error) {
	EnvToClusterMapping := make(map[int]int)
	devtronResourceSearchableKeyMap := impl.devtronResourceSearchableKey.GetAllSearchableKeyNameIdMap()
	environmentIds := make([]*int, 0)
	for _, qualifierMapping := range imageDigestQualifierMappings {
		if qualifierMapping.IdentifierKey == devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_ENV_ID] {
			environmentIds = append(environmentIds, &qualifierMapping.IdentifierValueInt)
		}
	}
	environments, err := impl.environmentRepository.FindByIds(environmentIds)
	if err != nil {
		impl.logger.Errorw("error in fetching environments by environmentIds", "err", err)
		return nil, err
	}
	for _, env := range environments {
		EnvToClusterMapping[env.Id] = env.ClusterId
	}
	return EnvToClusterMapping, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) CreateOrUpdatePolicyForCluster(policyRequest *PolicyRequest) (*PolicyRequest, error) {

	dbConnection := impl.qualifierMappingService.GetDbConnection()
	tx, _ := dbConnection.Begin()

	if policyRequest.AllClusters == true {
		err := impl.handleImageDigestPolicyForAllClusters(policyRequest.UserId, tx)
		if err != nil {
			impl.logger.Errorw("Error in saving image digest policy for all clusters", "err", err)
			return nil, err
		}
		_ = tx.Commit()
		return policyRequest, nil
	}

	devtronResourceSearchableKeyMap := impl.devtronResourceSearchableKey.GetAllSearchableKeyNameIdMap()

	// fetching already configured policies
	imageDigestQualifierMappings, err := impl.qualifierMappingService.GetQualifierMappingsByResourceType(resourceQualifiers.ImageDigest)
	if err != nil {
		impl.logger.Errorw("error in fetching qualifier mappings by resourceType", "resourceType: ", resourceQualifiers.ImageDigest)
		return nil, err
	}

	// exiting cluster and environments already having imageDigest configured
	ExistingClustersWithImageDigestPolicyConfigured, ExistingEnvironmentsWithImageDigestPolicyConfigured := getExistingClustersAndEnvsWithImagePullPolicyConfigured(imageDigestQualifierMappings, devtronResourceSearchableKeyMap)

	// saving image digest policy for new clusters and environments
	newClustersWithImageDigestPolicyConfigured, newEnvironmentsWithImageDigestPolicyConfigured, err :=
		impl.SaveNewPolicies(
			policyRequest,
			ExistingClustersWithImageDigestPolicyConfigured, ExistingEnvironmentsWithImageDigestPolicyConfigured, devtronResourceSearchableKeyMap, tx)
	if err != nil {
		impl.logger.Errorw("error in creating image digest policies", "err", err)
		return nil, err
	}

	// removing policies present in db but not present in request
	err = impl.removePoliciesNotPresentInRequest(imageDigestQualifierMappings,
		newClustersWithImageDigestPolicyConfigured,
		newEnvironmentsWithImageDigestPolicyConfigured,
		devtronResourceSearchableKeyMap,
		policyRequest.UserId,
		tx)
	if err != nil {
		impl.logger.Errorw("error in deleting policies not present in request but present in DB", "err", err)
		return nil, err
	}

	_ = tx.Commit()

	return policyRequest, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) handleImageDigestPolicyForAllClusters(userId int32, tx *pg.Tx) error {

	// step1: create image digest policy for all clusters by setting qualifierId = int(resourceQualifiers.GLOBAL_QUALIFIER)
	// step2: Delete individual cluster and env level image digest policy mappings

	globalQualifierMapping := &resourceQualifiers.QualifierMapping{
		ResourceId:   resourceQualifiers.ImageDigestResourceId,
		ResourceType: resourceQualifiers.ImageDigest,
		QualifierId:  int(resourceQualifiers.GLOBAL_QUALIFIER),
		Active:       true,
		AuditLog: sql.AuditLog{
			CreatedOn: time.Time{},
			CreatedBy: userId,
			UpdatedOn: time.Time{},
			UpdatedBy: userId,
		},
	}

	// creating image digest policy at global level
	_, err := impl.qualifierMappingService.CreateQualifierMappings([]*resourceQualifiers.QualifierMapping{globalQualifierMapping}, tx)
	if err != nil {
		impl.logger.Errorw("error in creating global image digest policy", "err", err)
		return err
	}

	// deleting all cluster and env level policies
	err = impl.qualifierMappingService.DeleteAllByResourceTypeAndQualifierId(
		resourceQualifiers.ImageDigest,
		resourceQualifiers.ImageDigestResourceId,
		[]int{int(resourceQualifiers.CLUSTER_QUALIFIER), int(resourceQualifiers.ENV_QUALIFIER)},
		userId,
		tx)
	if err != nil {
		impl.logger.Errorw("error in deleting resource by resource type, id and qualifier id", "err", err)
		return err
	}
	return nil
}

func getExistingClustersAndEnvsWithImagePullPolicyConfigured(imageDigestQualifierMappings []*resourceQualifiers.QualifierMapping, devtronResourceSearchableKeyMap map[bean.DevtronResourceSearchableKeyName]int) (map[ClusterId]bool, map[EnvironmentId]bool) {
	ExistingClustersWithImageDigestPolicyConfigured := make(map[ClusterId]bool)
	ExistingEnvironmentsWithImageDigestPolicyConfigured := make(map[EnvironmentId]bool)
	for _, existingMapping := range imageDigestQualifierMappings {
		if existingMapping.IdentifierKey == devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_CLUSTER_ID] {
			ExistingClustersWithImageDigestPolicyConfigured[existingMapping.IdentifierValueInt] = true
		} else if existingMapping.IdentifierKey == devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_ENV_ID] {
			ExistingEnvironmentsWithImageDigestPolicyConfigured[existingMapping.IdentifierValueInt] = true
		}
	}
	return ExistingClustersWithImageDigestPolicyConfigured, ExistingEnvironmentsWithImageDigestPolicyConfigured
}

func (impl ImageDigestQualifierMappingServiceImpl) SaveNewPolicies(
	policyRequest *PolicyRequest, ExistingClustersWithImageDigestPolicyConfigured map[ClusterId]bool,
	ExistingEnvironmentsWithImageDigestPolicyConfigured map[EnvironmentId]bool,
	devtronResourceSearchableKeyMap map[bean.DevtronResourceSearchableKeyName]int, tx *pg.Tx) (map[ClusterId]bool, map[EnvironmentId]bool, error) {

	newPolicies := make([]*resourceQualifiers.QualifierMapping, 0)
	newClustersWithImageDigestPolicyConfigured := make(map[ClusterId]bool)
	newEnvironmentsWithImageDigestPolicyConfigured := make(map[EnvironmentId]bool)

	for _, policy := range policyRequest.ClusterDetails {
		if policy.PolicyType == ALL_ENVIRONMENTS {
			if _, ok := ExistingClustersWithImageDigestPolicyConfigured[policy.ClusterId]; !ok {
				qualifierMapping := &resourceQualifiers.QualifierMapping{
					ResourceId:         resourceQualifiers.ImageDigestResourceId,
					ResourceType:       resourceQualifiers.ImageDigest,
					QualifierId:        int(resourceQualifiers.CLUSTER_QUALIFIER),
					IdentifierKey:      devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_CLUSTER_ID],
					IdentifierValueInt: policy.ClusterId,
					Active:             true,
					AuditLog: sql.AuditLog{
						CreatedOn: time.Now(),
						CreatedBy: policyRequest.UserId,
						UpdatedOn: time.Now(),
						UpdatedBy: policyRequest.UserId,
					},
				}
				newPolicies = append(newPolicies, qualifierMapping)
			}
			newClustersWithImageDigestPolicyConfigured[policy.ClusterId] = true
		} else if policy.PolicyType == SELECTED_ENVIRONMENTS {
			for _, envId := range policy.Environments {
				if _, ok := ExistingEnvironmentsWithImageDigestPolicyConfigured[policy.ClusterId]; !ok {
					qualifierMapping := &resourceQualifiers.QualifierMapping{
						ResourceId:         resourceQualifiers.ImageDigestResourceId,
						ResourceType:       resourceQualifiers.ImageDigest,
						QualifierId:        int(resourceQualifiers.ENV_QUALIFIER),
						IdentifierKey:      devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_ENV_ID],
						IdentifierValueInt: envId,
						Active:             true,
						AuditLog: sql.AuditLog{
							CreatedOn: time.Now(),
							CreatedBy: policyRequest.UserId,
							UpdatedOn: time.Now(),
							UpdatedBy: policyRequest.UserId,
						},
					}
					newPolicies = append(newPolicies, qualifierMapping)
				}
				newEnvironmentsWithImageDigestPolicyConfigured[envId] = true
			}
		}
	}
	_, err := impl.qualifierMappingService.CreateQualifierMappings(newPolicies, tx)
	if err != nil {
		impl.logger.Errorw("error in creating qualifier mappings for image digest policy", "err", err)
		return newClustersWithImageDigestPolicyConfigured, newEnvironmentsWithImageDigestPolicyConfigured, err
	}
	return newClustersWithImageDigestPolicyConfigured, newEnvironmentsWithImageDigestPolicyConfigured, nil
}

func (impl ImageDigestQualifierMappingServiceImpl) removePoliciesNotPresentInRequest(imageDigestQualifierMappings []*resourceQualifiers.QualifierMapping,
	newClustersWithImageDigestPolicyConfigured map[ClusterId]bool,
	newEnvironmentsWithImageDigestPolicyConfigured map[EnvironmentId]bool,
	devtronResourceSearchableKeyMap map[bean.DevtronResourceSearchableKeyName]int,
	UserId int32,
	tx *pg.Tx) error {

	policiesToBeRemovedIDs := make([]int, 0)

	for _, existingMapping := range imageDigestQualifierMappings {
		removePolicy := false
		if existingMapping.IdentifierKey == devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_CLUSTER_ID] {
			if _, ok := newClustersWithImageDigestPolicyConfigured[existingMapping.IdentifierValueInt]; !ok {
				removePolicy = true
			}
		} else if existingMapping.IdentifierKey == devtronResourceSearchableKeyMap[bean.DEVTRON_RESOURCE_SEARCHABLE_KEY_ENV_ID] {
			if _, ok := newEnvironmentsWithImageDigestPolicyConfigured[existingMapping.IdentifierValueInt]; !ok {
				removePolicy = true
			}
		} else if existingMapping.QualifierId == int(resourceQualifiers.GLOBAL_QUALIFIER) {
			// removing global policy because if we are here AllClusters=false in request
			removePolicy = true
		}
		if removePolicy {
			policiesToBeRemovedIDs = append(policiesToBeRemovedIDs, existingMapping.Id)
		}
	}

	if len(policiesToBeRemovedIDs) > 0 {
		err := impl.qualifierMappingService.DeleteAllByIds(policiesToBeRemovedIDs, UserId, tx)
		if err != nil {
			impl.logger.Errorw("error in deleting old policies", "err", err)
			return err
		}
	}

	return nil
}

func (impl ImageDigestQualifierMappingServiceImpl) IsPolicyConfiguredAtGlobalLevel(envId int, clusterId int) (bool, error) {
	if clusterId == 0 {
		env, err := impl.environmentRepository.FindById(envId)
		if err != nil {
			impl.logger.Errorw("error in fetching environment by envId", "err", err, "envId", envId)
		}
		clusterId = env.ClusterId
	}
	scope := &resourceQualifiers.Scope{EnvId: envId, ClusterId: clusterId}
	resourceIds := []int{resourceQualifiers.ImageDigestResourceId}
	qualifierMappings, err := impl.qualifierMappingService.GetQualifierMappings(resourceQualifiers.ImageDigest, scope, resourceIds)
	if err != nil && err != pg.ErrNoRows {
		return false, err
	}
	if len(qualifierMappings) == 0 {
		return false, nil
	}
	return true, nil
}
