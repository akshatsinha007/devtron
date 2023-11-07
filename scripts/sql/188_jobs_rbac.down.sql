ALTER TABLE roles DROP COLUMN workflow;
DELETE from rbac_role_resource_detail where resource='workflow';
UPDATE rbac_policy_resource_detail set eligible_entity_access_types = ARRAY['apps/devtron-app','apps/helm-app'] where resource='project' OR resource ='global-environment' OR resource='terminal';
UPDATE rbac_role_resource_detail set eligible_entity_access_types = ARRAY['apps/devtron-app','apps/helm-app'] where resource ='project' OR resource ='environment';
DELETE FROM rbac_policy_resource_detail where resource='jobEnv';
DELETE FROM rbac_policy_resource_detail where resource='workflow';
DELETE FROM "public"."rbac_role_data" where entity='jobs';
DELETE FROM "public"."rbac_policy_data" where entity='jobs';
