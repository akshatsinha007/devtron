/*
 * Copyright (c) 2024. Devtron Inc.
 */

ALTER TABLE "public"."docker_artifact_store" ADD COLUMN "connection" varchar(250);

ALTER TABLE "public"."docker_artifact_store" ADD COLUMN "cert" text;