/*
 * Copyright (c) 2024. Devtron Inc.
 */

CREATE INDEX "pco_ci_artifact_id" ON "public"."pipeline_config_override" USING BTREE ("ci_artifact_id");