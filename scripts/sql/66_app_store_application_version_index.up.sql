/*
 * Copyright (c) 2024. Devtron Inc.
 */

--- Create index on app_store_id in app_store_application_version
CREATE INDEX IF NOT EXISTS app_store_application_version_app_store_id_IX ON public.app_store_application_version USING btree(app_store_id);