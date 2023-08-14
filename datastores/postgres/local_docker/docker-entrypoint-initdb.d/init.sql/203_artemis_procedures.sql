CREATE TABLE "public"."orchestrations" (
   "orchestration_id" int8 NOT NULL DEFAULT next_id(),
   "org_id" int8 NOT NULL REFERENCES orgs(org_id),
   "active" bool NOT NULL DEFAULT false,
   "type" text NOT NULL DEFAULT 'zeus',
   "group_name" text NOT NULL DEFAULT 'zeus',
   "orchestration_name" text NOT NULL,
   "instructions" jsonb NOT NULL DEFAULT {},
   "updated_at" timestamptz  NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."orchestrations" ADD CONSTRAINT "orchestrations_pk" PRIMARY KEY ("orchestration_id");
ALTER TABLE "public"."orchestrations" ADD CONSTRAINT "orchestrations_uniq_name_to_org" UNIQUE ("org_id", "orchestration_name");
CREATE INDEX "orchestrations_last_updated_at_index" ON "public"."orchestrations" (updated_at ASC);
CREATE INDEX orchestrations_active ON orchestrations(active);
CREATE INDEX orchestrations_name_ind ON orchestrations(orchestration_name);
CREATE INDEX orchestrations_group_ind ON orchestrations(group_name);
CREATE INDEX orchestrations_type_ind ON orchestrations(type);

CREATE TRIGGER set_timestamp_on_orchestrations_updated
AFTER UPDATE ON orchestrations
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "public"."orchestrations_scheduled_to_cloud_ctx_ns" (
    "orchestration_schedule_id" int8 NOT NULL DEFAULT next_id(),
    "orchestration_id" int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    "cloud_ctx_ns_id" int8 NOT NULL REFERENCES topologies_org_cloud_ctx_ns(cloud_ctx_ns_id),
    "status" text NOT NULL DEFAULT 'Pending',
    "date_scheduled" timestamptz NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."orchestrations_scheduled_to_cloud_ctx_ns" ADD CONSTRAINT "orchestrations_scheduled_to_cloud_ctx_ns_pk" PRIMARY KEY ("orchestration_schedule_id");
