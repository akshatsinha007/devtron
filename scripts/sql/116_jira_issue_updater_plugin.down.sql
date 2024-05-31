/*
 * Copyright (c) 2024. Devtron Inc.
 */

DELETE FROM plugin_step_variable
WHERE plugin_step_id = (SELECT ps.id FROM plugin_metadata p inner JOIN plugin_step ps on ps.plugin_id=p.id WHERE p.name='Jira Issue Updater' and ps."index"=1 and ps.deleted=false);