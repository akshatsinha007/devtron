/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE ci_build_config ADD COLUMN IF NOT EXISTS "use_root_context" bool;