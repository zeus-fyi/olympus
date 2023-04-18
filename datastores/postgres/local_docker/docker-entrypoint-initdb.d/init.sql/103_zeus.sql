
-- class types skeleton, infrastructure, configuration, test_suites_base, cluster, matrix, system
CREATE TABLE "public"."topology_class_types" (
    "topology_class_type_id" int8 NOT NULL DEFAULT next_id(),
    "topology_class_type_name" text
);
ALTER TABLE "public"."topology_class_types" ADD CONSTRAINT "topology_class_types_pk" PRIMARY KEY ("topology_class_type_id");

-- specific sw names. eg. prysm_validator_client
CREATE TABLE "public"."topologies" (
   "topology_id" int8 NOT NULL DEFAULT next_id(),
   "name" text NOT NULL
);
ALTER TABLE "public"."topologies" ADD CONSTRAINT "topology_pk" PRIMARY KEY ("topology_id");

-- kns id table
CREATE TABLE "public"."topologies_kns" (
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "cloud_provider" text NOT NULL,
    "context" text NOT NULL,
    "region" text NOT NULL,
    "namespace" text NOT NULL,
    "env" text NOT NULL
);
ALTER TABLE "public"."topologies_kns" ADD CONSTRAINT "kns_pk" PRIMARY KEY ("topology_id", "cloud_provider", "region", "context","namespace", "env");
CREATE INDEX topologies_kns_topology_id_idx ON topologies_kns (topology_id);


-- specific deployed topology to user (statuses can be pending, terminated, etc)
CREATE TABLE "public"."topologies_deployed" (
   "deployment_id" int8 NOT NULL DEFAULT next_id(),
   "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
   "topology_status" text NOT NULL,
   "updated_at" timestamptz  NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."topologies_deployed" ADD CONSTRAINT "topologies_deployed_pk" PRIMARY KEY ("topology_id","deployment_id");
CREATE INDEX status_idx ON topologies_deployed (topology_status);
CREATE INDEX topologies_deployed_idx ON topologies_deployed (topology_id);

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