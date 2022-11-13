-- links topology config to infra
CREATE TABLE "public"."topology_configuration_class" (
    "topology_configuration_class_id" int8 DEFAULT next_id(),
    "topology_infrastructure_component_id" int8 NOT NULL REFERENCES topology_infrastructure_components(topology_infrastructure_component_id)
);
ALTER TABLE "public"."topology_configuration_class" ADD CONSTRAINT "topology_configuration_class_pk" PRIMARY KEY ("topology_configuration_class_id");

-- links topology to child_value
CREATE TABLE "public"."topology_configuration_child_values_overrides" (
    "topology_configuration_child_values_override_id" int8 DEFAULT next_id(),
    "topology_configuration_class_id" int8 NOT NULL REFERENCES topology_configuration_class(topology_configuration_class_id),
    "chart_subcomponent_child_values_id" int8 NOT NULL REFERENCES chart_subcomponents_child_values(chart_subcomponent_child_values_id),
    "chart_subcomponent_override_value" text NOT NULL
);

ALTER TABLE "public"."topology_configuration_child_values_overrides" ADD CONSTRAINT "topology_configuration_child_values_overrides_pk" PRIMARY KEY ("topology_configuration_child_values_override_id");
ALTER TABLE "public"."topology_configuration_child_values_overrides" ADD CONSTRAINT "config_override_unique" UNIQUE("topology_configuration_class_id", "chart_subcomponent_child_values_id");


