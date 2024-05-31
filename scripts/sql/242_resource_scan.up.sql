/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE public.scan_tool_execution_history_mapping ADD COLUMN IF NOT EXISTS error_message  varchar NULL;
INSERT INTO public.scan_tool_step(scan_tool_id, index, step_execution_type, step_execution_sync, retry_count, execute_step_on_fail, execute_step_on_pass, render_input_data_from_step, http_input_payload, http_method_type, http_req_headers, http_query_params, cli_command, cli_output_type, deleted, created_on, created_by, updated_on, updated_by) VALUES (3,5,'CLI',true,1,-1,-1,-1,null,null,null,null,'trivy image -f json -o {{.OUTPUT_FILE_PATH}} --timeout {{.timeout}} {{.IMAGE_NAME}}', 'STATIC',false,now()::timestamp,'1',now()::timestamp,'1');
