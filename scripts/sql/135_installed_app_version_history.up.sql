/*
 * Copyright (c) 2024. Devtron Inc.
 */

alter table installed_app_version_history add column started_on timestamp with time zone, add column finished_on timestamp with time zone;