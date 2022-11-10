
-- class types skeleton, infrastructure, configuration, base, cluster, matrix, system
CREATE TABLE "public"."topology_class_types" (
    "topology_class_type_id" int8 NOT NULL DEFAULT next_id(),
    "topology_class_type_name" text
);
ALTER TABLE "public"."topology_class_types" ADD CONSTRAINT "topology_class_types_pk" PRIMARY KEY ("topology_class_type_id");

-- specific use class names. eg. eth_validator_client
CREATE TABLE "public"."topology_classes" (
   "topology_class_id" int8 NOT NULl DEFAULT next_id(),
   "topology_class_type_id" int8 NOT NULL DEFAULT next_id(),
   "topology_class_name" text NOT NULL
);
ALTER TABLE "public"."topology_classes" ADD CONSTRAINT "topology_class_pk" PRIMARY KEY ("topology_class_id");

-- specific sw names. eg. prysm_validator_client
CREATE TABLE "public"."topologies" (
   "topology_id" int8 NOT NULL DEFAULT next_id(),
   "name" text NOT NULL
);
ALTER TABLE "public"."topologies" ADD CONSTRAINT "topology_pk" PRIMARY KEY ("topology_id");

-- links components to build higher level topologies eg. beacon + exec = full eth2 beacon cluster
CREATE TABLE "public"."topology_dependent_components" (
    "topology_class_id" int8 NOT NULL REFERENCES topology_classes(topology_class_id),
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id)
);

-- kns id table
CREATE TABLE "public"."topologies_kns" (
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "cloud_provider" text NOT NULL,
    "context" text NOT NULL,
    "region" text NOT NULL,
    "namespace" text NOT NULL,
    "env" text NOT NULL
);
ALTER TABLE "public"."topologies_kns" ADD CONSTRAINT "kns_pk" PRIMARY KEY ("cloud_provider", "region", "context","namespace", "env");

-- specific deployed topology to user (statuses can be pending, terminated, etc)
CREATE TABLE "public"."topologies_deployed" (
   "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
   "topology_status" text NOT NULL,
   "updated_at" timestamptz  NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."topologies_deployed" ADD CONSTRAINT "topologies_deployed_unique_key" UNIQUE  ("topology_id", "topology_status");

-- if needed again per different schema. eg zeus.DB, vs eth.DB
-- CREATE OR REPLACE FUNCTION trigger_set_timestamp()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = NOW();
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_on_topologies_deployed
    BEFORE UPDATE ON topologies_deployed
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();