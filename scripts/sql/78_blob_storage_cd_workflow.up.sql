/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE cd_workflow_runner
    ADD COLUMN IF NOT EXISTS  blob_storage_enabled boolean NOT NULL DEFAULT true;