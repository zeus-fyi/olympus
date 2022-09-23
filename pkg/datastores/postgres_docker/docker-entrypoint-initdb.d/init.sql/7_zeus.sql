
-- class types skeleton, infrastructure, configuration, base, cluster, matrix, system
CREATE TABLE "public"."topology_class_types" (
    "topology_class_type_id" int8 NOT NULL,
    "topology_class_type_name" text
);
ALTER TABLE "public"."topology_class_types" ADD CONSTRAINT "topology_class_types_pk" PRIMARY KEY ("topology_class_type_id");

-- specific use class names. eg. eth_validator_client
CREATE TABLE "public"."topology_classes" (
   "topology_class_id" int8 NOT NULL,
   "topology_class_type_id" int8 NOT NULL,
   "topology_class_name" text NOT NULL
);
ALTER TABLE "public"."topology_classes" ADD CONSTRAINT "topology_class_pk" PRIMARY KEY ("topology_class_id");

-- specific sw names. eg. prysm_validator_client
CREATE TABLE "public"."topologies" (
   "topology_id" int8 NOT NULL,
   "name" text NOT NULL
);
ALTER TABLE "public"."topologies" ADD CONSTRAINT "topology_pk" PRIMARY KEY ("topology_id");

-- links components to build higher level topologies eg. beacon + exec = full eth2 beacon cluster
CREATE TABLE "public"."topology_dependent_components" (
    "topology_class_id" int8 NOT NULL REFERENCES topology_classes(topology_class_id),
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id)
);

-- kns id table
CREATE TABLE "public"."kns" (
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "context" text NOT NULL,
    "namespace" text NOT NULL,
    "env" text NOT NULL
);
ALTER TABLE "public"."kns" ADD CONSTRAINT "kns_pk" PRIMARY KEY ("context","namespace", "env");

-- specific sw names. eg. prysm_validator_client
CREATE TABLE "public"."deployed_topologies" (
   "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
   "org_id" int8 NOT NULL REFERENCES orgs(org_id),
   "user_id" int8 NOT NULL REFERENCES users(user_id),
   "name" text
);