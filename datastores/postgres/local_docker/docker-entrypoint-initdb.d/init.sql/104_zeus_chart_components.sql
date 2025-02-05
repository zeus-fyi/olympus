-- chart component kind: e.g. service, statefulset, etc
CREATE TABLE "public".chart_component_resources(
    "chart_component_resource_id" int8 NOT NULL DEFAULT next_id(),
    "chart_component_kind_name" text NOT NULL,
    "chart_component_api_version" text NOT NULL
);
ALTER TABLE "public".chart_component_resources ADD CONSTRAINT "chart_component_kind_types_pk" PRIMARY KEY (chart_component_resource_id);
ALTER TABLE "public".chart_component_resources ADD CONSTRAINT "chart_component_kinds_api_version_pk" UNIQUE (chart_component_kind_name, "chart_component_api_version");

-- synthetic helm package, eg eth_validator_client_package
CREATE TABLE "public"."chart_packages" (
    "chart_package_id" int8 NOT NULL DEFAULT next_id(),
    "chart_name" text NOT NULL,
    "chart_version" text NOT NULL,
    "chart_description" text
);
ALTER TABLE "public"."chart_packages" ADD CONSTRAINT "chart_package_pk" PRIMARY KEY ("chart_package_id");
ALTER TABLE "public"."chart_packages" ADD CONSTRAINT "chart_package_unique" UNIQUE("chart_name", "chart_version");

-- use naming to match k8s conventional naming: eg. StatefulSetSpec
-- synthetic chart subcomponent classes, eg container, ports, volume_mounts, etc
-- synthetic chart subcomponent classes are flexible, eg. ports as a component of a container
-- this allows hierarchy stacking
-- spec:
--  serviceName: "nginx"
--  replicas: 2
--  selector:
--    matchLabels:
--      app: nginx
--  template: (chart_subcomponent_parent_class_types)
--    metadata: (chart_subcomponent_child_class_types)
--      labels: (chart_subcomponent_child_class_types)
--        app: nginx (chart_subcomponents_child_values)
--    spec: (chart_subcomponent_child_class_types)
--      containers: (chart_subcomponent_child_class_types)
--      - name: nginx (chart_subcomponents_child_values)
--      image: registry.k8s.io/nginx-slim:0.8 (chart_subcomponents_child_values)
--      ports: (chart_subcomponent_child_class_types)
--       - containerPort: 80 (chart_subcomponents_child_values)
--        name: web
--      volumeMounts:
--       - name: www
--       mountPath: /usr/share/nginx/html

---------------------------------------------------------------------------------------------------------------------------------------------------------------------
-- eg. spec, metadata (top parent fields) -> (chart_subcomponent_child_class_types)
CREATE TABLE "public"."chart_subcomponent_parent_class_types" (
    "chart_package_id" int8 NOT NULL REFERENCES chart_packages(chart_package_id),
    "chart_component_resource_id" int8 NOT NULL REFERENCES chart_component_resources(chart_component_resource_id),
    "chart_subcomponent_parent_class_type_id" int8 DEFAULT next_id(),
    "chart_subcomponent_parent_class_type_name" text NOT NULL
);
ALTER TABLE "public"."chart_subcomponent_parent_class_types" ADD CONSTRAINT "chart_subcomponent_parent_class_types_pk" PRIMARY KEY ("chart_subcomponent_parent_class_type_id");

-- flexible keys schema ----------------------------------------------------------------------------------------------------------------------------------------------------
-- to add/change fields as needed without validation for now
-- chart_subcomponent_child_class_types
-- e.g. deploymentSpec, statefulSetSpec
CREATE TABLE "public"."chart_subcomponent_child_class_types" (
    "chart_subcomponent_parent_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_parent_class_types(chart_subcomponent_parent_class_type_id),
    "chart_subcomponent_child_class_type_id" int8 DEFAULT next_id(),
    "chart_subcomponent_child_class_type_name" text NOT NULL
);
ALTER TABLE "public"."chart_subcomponent_child_class_types" ADD CONSTRAINT "chart_subcomponent_child_class_types_pk" PRIMARY KEY ("chart_subcomponent_child_class_type_id");

-- flexible schema values mapping text k-v and jsonb -------------------------------------------------------------------------------------------------------------------------------------
-- link to synthetic chart subcomponent child key->values, use bool toggle to generate a controller for the package
CREATE TABLE "public"."chart_subcomponents_child_values" (
    "chart_subcomponent_child_values_id" int8 DEFAULT next_id(),
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "chart_subcomponent_chart_package_template_injection" bool NOT NULL DEFAULT false,
    "chart_subcomponent_key_name" text NOT NULL,
    "chart_subcomponent_value" text NOT NULL
);
ALTER TABLE "public"."chart_subcomponents_child_values" ADD CONSTRAINT "chart_subcomponents_child_values_pk" PRIMARY KEY ("chart_subcomponent_child_values_id");

---------------------------------------------------------------------------------------------------------------------------------------------------------------------
-- tables to jump links
-- links package to chart subcomponents
CREATE TABLE "public"."chart_package_components" (
   "chart_package_id" int8 NOT NULL REFERENCES chart_packages(chart_package_id),
   "chart_subcomponent_parent_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_parent_class_types(chart_subcomponent_parent_class_type_id)
);
