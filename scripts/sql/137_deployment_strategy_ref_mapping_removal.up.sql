/*
 * Copyright (c) 2024. Devtron Inc.
 */

UPDATE global_strategy_metadata_chart_ref_mapping SET active=false WHERE chart_ref_id in(SELECT id FROM chart_ref WHERE location LIKE '%deployment-chart_%') AND global_strategy_metadata_id=(SELECT id FROM global_strategy_metadata WHERE name='RECREATE');
