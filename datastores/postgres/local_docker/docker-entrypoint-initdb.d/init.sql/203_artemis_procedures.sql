CREATE TABLE "public"."orchestrations" (
   "orchestration_id" int8 NOT NULL DEFAULT next_id(),
   "org_id" int8 NOT NULL REFERENCES orgs(org_id),
   "orchestration_name" text NOT NULL,
   "instructions" jsonb NOT NULL DEFAULT {}
);
ALTER TABLE "public"."orchestrations" ADD CONSTRAINT "orchestrations_pk" PRIMARY KEY ("orchestration_id");
ALTER TABLE "public"."orchestrations" ADD CONSTRAINT "orchestrations_uniq_name_to_org" UNIQUE ("org_id", "orchestration_name");

CREATE TABLE "public"."orchestrations_scheduled_to_cloud_ctx_ns" (
    "orchestration_schedule_id" int8 NOT NULL DEFAULT next_id(),
    "orchestration_id" int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    "cloud_ctx_ns_id" int8 NOT NULL REFERENCES topologies_org_cloud_ctx_ns(cloud_ctx_ns_id),
    "status" text NOT NULL DEFAULT 'Pending',
    "date_scheduled" timestamptz NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."orchestrations_scheduled_to_cloud_ctx_ns" ADD CONSTRAINT "orchestrations_scheduled_to_cloud_ctx_ns_pk" PRIMARY KEY ("orchestration_schedule_id");
