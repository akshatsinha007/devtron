/*
 * Copyright (c) 2024. Devtron Inc.
 */

DROP index image_scan_deploy_info_unique;
CREATE INDEX image_scan_deploy_info_unique ON public.image_scan_deploy_info USING btree (scan_object_meta_id, object_type);
