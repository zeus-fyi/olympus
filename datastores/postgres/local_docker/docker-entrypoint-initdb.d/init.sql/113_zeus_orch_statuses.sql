
CREATE TABLE "public"."orchestrations_cloud_ctx_ns_logs" (
    "log_id" int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    "orchestration_id" int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    "cloud_ctx_ns_id" int8 REFERENCES topologies_org_cloud_ctx_ns(cloud_ctx_ns_id) ON DELETE SET NULL,
    "status" text NOT NULL DEFAULT 'Pending',
    "msg" text NOT NULL DEFAULT ''
);
CREATE INDEX "orchestrations_cloud_ctx_ns_logs_orchestration_id" ON "public"."orchestrations_cloud_ctx_ns_logs"("orchestration_id");
CREATE INDEX "orchestrations_cloud_ctx_ns_logs_cloud_ctx_ns_id" ON "public"."orchestrations_cloud_ctx_ns_logs"("cloud_ctx_ns_id");