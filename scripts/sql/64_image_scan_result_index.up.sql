/*
 * Copyright (c) 2024. Devtron Inc.
 */

--- Create index on image_scan_execution_history_id in image_scan_execution_result
CREATE INDEX IF NOT EXISTS image_scan_execution_history_id_IX ON public.image_scan_execution_result (image_scan_execution_history_id);