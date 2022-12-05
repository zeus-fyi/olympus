--- WIP BELOW
-- links topology config to infra
-- CREATE TABLE "public"."topology_configuration_class" (
--     "topology_configuration_class_id" int8 DEFAULT next_id(),
--     "topology_infrastructure_component_id" int8 NOT NULL REFERENCES topology_infrastructure_components(topology_infrastructure_component_id)
-- );
-- ALTER TABLE "public"."topology_configuration_class" ADD CONSTRAINT "topology_configuration_class_pk" PRIMARY KEY ("topology_configuration_class_id");
--
-- -- links topology to child_value
-- CREATE TABLE "public"."topology_configuration_child_values_overrides" (
--     "topology_configuration_child_values_override_id" int8 DEFAULT next_id(),
--     "topology_configuration_class_id" int8 NOT NULL REFERENCES topology_configuration_class(topology_configuration_class_id),
--     "chart_subcomponent_child_values_id" int8 NOT NULL REFERENCES chart_subcomponents_child_values(chart_subcomponent_child_values_id),
--     "chart_subcomponent_override_value" text NOT NULL
-- );
--
-- ALTER TABLE "public"."topology_configuration_child_values_overrides" ADD CONSTRAINT "topology_configuration_child_values_overrides_pk" PRIMARY KEY ("topology_configuration_child_values_override_id");
-- ALTER TABLE "public"."topology_configuration_child_values_overrides" ADD CONSTRAINT "config_override_unique" UNIQUE("topology_configuration_class_id", "chart_subcomponent_child_values_id");
--
--- WIP ABOVE
-- links cluster+ topologies to bases

-- class types skeleton, infrastructure, configuration, test_suites_base, cluster, matrix, system
-- CREATE TABLE "public"."topology_class_types" (
--                                                  "topology_class_type_id" int8 NOT NULL DEFAULT next_id(),
--                                                  "topology_class_type_name" text
-- );
-- ALTER TABLE "public"."topology_class_types" ADD CONSTRAINT "topology_class_types_pk" PRIMARY KEY ("topology_class_type_id");

-- use the distributed id gen to also function as the version
CREATE TABLE "public"."topology_system_components" (
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "topology_system_component_id" int8 DEFAULT next_id(),
    "topology_class_type_id" int8 NOT NULL REFERENCES topology_class_types(topology_class_type_id) CHECK (topology_class_type_id >=4) NOT NULL,
    "topology_system_component_name" text NOT NULL
);
ALTER TABLE "public"."topology_system_components" ADD CONSTRAINT "topology_system_components_pk" PRIMARY KEY ("topology_system_component_id");
ALTER TABLE "public"."topology_system_components" ADD CONSTRAINT "system_component_name_unique_to_org" UNIQUE("topology_system_component_name", "topology_class_type_id", "org_id");

-- links components to build higher level topologies eg. beacon + exec = full eth2 beacon cluster

-- enforce the class type is base on entry, use the distributed id gen to also function as the version
CREATE TABLE "public"."topology_base_components" (
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "topology_system_component_id" int8 NOT NULL REFERENCES topology_system_components(topology_system_component_id),
    "topology_class_type_id" int8 NOT NULL REFERENCES topology_class_types(topology_class_type_id) CHECK (topology_class_type_id <= 3 AND topology_class_type_id >=3) NOT NULL DEFAULT 3,
    "topology_base_component_id" int8 DEFAULT next_id(),
    "topology_base_name" text NOT NULL
);
ALTER TABLE "public"."topology_base_components" ADD CONSTRAINT "topology_base_components_pk" PRIMARY KEY ("topology_base_component_id");
ALTER TABLE "public"."topology_base_components" ADD CONSTRAINT "base_component_name_unique_to_org" UNIQUE("topology_base_name", "org_id");


-- enforce the class type is base on entry, use the distributed id gen to also function as the version
CREATE TABLE "public"."topology_skeleton_base_components" (
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "topology_base_component_id" int8 NOT NULL REFERENCES topology_base_components(topology_base_component_id) NOT NULL,
    "topology_class_type_id" int8 NOT NULL REFERENCES topology_class_types(topology_class_type_id) CHECK (topology_class_type_id <= 2 AND topology_class_type_id >=2) NOT NULL DEFAULT 2,
    "topology_skeleton_base_id" int8 DEFAULT next_id(),
    "topology_skeleton_base_name" text NOT NULL
);

ALTER TABLE "public"."topology_skeleton_base_components" ADD CONSTRAINT "topology_skeleton_base_components_pk" PRIMARY KEY ("topology_base_component_id");
ALTER TABLE "public"."topology_skeleton_base_components" ADD CONSTRAINT "topology_skeleton_base_component_name_unique_to_org" UNIQUE("topology_skeleton_base_name", "topology_skeleton_base_id", "topology_class_type_id", "org_id");
ALTER TABLE "public"."topology_skeleton_base_components" ADD CONSTRAINT "topology_skeleton_base_version_id_uniq" UNIQUE("topology_skeleton_base_id");

-- links topology to kubernetes package
CREATE TABLE "public"."topology_infrastructure_components" (
    "topology_skeleton_base_id" int8 NOT NULL REFERENCES topology_skeleton_base_components(topology_skeleton_base_id) NOT NULL,
    "topology_infrastructure_component_id" int8 DEFAULT next_id(),
    "topology_id" int8 NOT NULL REFERENCES topologies(topology_id),
    "chart_package_id" int8 NOT NULL REFERENCES chart_packages(chart_package_id)
);
ALTER TABLE "public"."topology_infrastructure_components" ADD CONSTRAINT "topology_infrastructure_components_pk" PRIMARY KEY ("topology_infrastructure_component_id");

-- ALTER TABLE topology_infrastructure_components ADD COLUMN topology_skeleton_base_id int8;
-- ALTER TABLE topology_infrastructure_components ALTER COLUMN topology_skeleton_base_id SET NOT NULL;
-- ALTER TABLE topology_infrastructure_components ADD CONSTRAINT topology_infrastructure_components_fk FOREIGN KEY (topology_skeleton_base_id) REFERENCES topology_skeleton_base_components (topology_skeleton_base_id) MATCH FULL;
