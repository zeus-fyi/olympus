---------------
-- CONTAINERS--
---------------
-- for template podSpec in deploymentSpec, statefulsetSpec, etc

CREATE TABLE "public"."containers" (
    "container_id" int8 DEFAULT next_id(),
    "container_name" text NOT NULL,
    "container_image_id" text NOT NULL,
    "container_version_tag" text NOT NULL DEFAULT 'latest',
    "container_platform_os" text NOT NULL,
    "container_repository" text NOT NULL,
    "container_image_pull_policy" text CHECK (container_image_pull_policy IN ('IfNotPresent', 'Always', 'Never')) NOT NULL DEFAULT 'IfNotPresent',
    "is_init_container" bool NOT NULL DEFAULT false
);
ALTER TABLE "public"."containers" ADD CONSTRAINT "containers_pk" PRIMARY KEY ("container_id");
ALTER TABLE "public"."containers" ADD CONSTRAINT "containers_version_pk" UNIQUE ("container_name","container_image_id", "container_version_tag", "container_platform_os");

-- lazy sort using time -> unique increased number on order input
CREATE TABLE "public"."chart_subcomponent_spec_pod_template_containers" (
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id),
    "container_sort_order" int8 NOT NULL DEFAULT next_id()
);
ALTER TABLE "public"."chart_subcomponent_spec_pod_template_containers" ADD CONSTRAINT "containers_order_pk" UNIQUE ("chart_subcomponent_child_class_type_id", "container_id", "container_sort_order");

------------
-- PORTS----
------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

CREATE TABLE "public"."container_ports" (
    -- pk for table
    "port_id" int8 NOT NULL DEFAULT next_id(),
    "port_name" text NOT NULL,
    "container_port" int NOT NULL,
    "host_ip" text NOT NULL DEFAULT '',
    "host_port" int NOT NULL,
    "port_protocol" text CHECK (port_protocol IN ('UDP', 'TCP', 'SCTP')) NOT NULL DEFAULT 'TCP'
);
ALTER TABLE "public"."container_ports" ADD CONSTRAINT "container_ports_pk" PRIMARY KEY ("port_id");

-- tables to jump links for ports
CREATE TABLE "public"."containers_ports" (
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id),
    "port_id" int8 NOT NULL REFERENCES container_ports(port_id)
);
-- lazy unique to all pods in spec chart
ALTER TABLE "public"."containers_ports" ADD CONSTRAINT "containers_ports_pk" UNIQUE ("chart_subcomponent_child_class_type_id", "port_id");

------------
--ENV_VARS--
------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

-- tables for env variables
CREATE TABLE "public"."container_environmental_vars" (
    -- pk for table (maybe make it a hash of the key/value?)
    "env_id" int8 NOT NULL DEFAULT next_id(),
    "name" text NOT NULL,
    "value" jsonb NOT NULL DEFAULT '{}'::jsonb
);
ALTER TABLE "public"."container_environmental_vars" ADD CONSTRAINT "environmental_vars_pk" PRIMARY KEY ("env_id");

-- tables for containers_environmental_vars links
CREATE TABLE "public"."containers_environmental_vars" (
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id),
    "env_id" int8 NOT NULL REFERENCES container_environmental_vars(env_id)
);
ALTER TABLE "public"."containers_environmental_vars" ADD CONSTRAINT "container_env_pk" UNIQUE ("container_id", "env_id");

------------------
-- VOLUME_MOUNTS--
------------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

-- tables for containers_volume_mounts
CREATE TABLE "public"."container_volume_mounts" (
    "volume_mount_id" int8 NOT NULL DEFAULT next_id(),
    "volume_mount_path" text NOT NULL,
    "volume_name" text NOT NULL,
    "volume_read_only" bool NOT NULL DEFAULT false,
    "volume_sub_path" text NOT NULL DEFAULT ''
);
ALTER TABLE "public"."container_volume_mounts" ADD CONSTRAINT "container_volume_mounts_pk" PRIMARY KEY ("volume_mount_id");

-- tables for containers_volume_mounts links
CREATE TABLE "public"."containers_volume_mounts" (
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id),
    "volume_mount_id" int8 NOT NULL REFERENCES container_volume_mounts(volume_mount_id)
);
ALTER TABLE "public"."containers_volume_mounts" ADD CONSTRAINT "containers_volume_mounts_pk" UNIQUE ("container_id", "volume_mount_id");

------------
-- VOLUMES--
------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

-- tables for volumes
CREATE TABLE "public"."volumes" (
    "volume_id" int8 NOT NULL DEFAULT next_id(),
    "volume_name" text NOT NULL,
    "volume_key_values_jsonb" jsonb NOT NULL DEFAULT '{}'::jsonb
);
ALTER TABLE "public"."volumes" ADD CONSTRAINT "volumes_pk" PRIMARY KEY ("volume_id");

-- tables for containers_volumes links
CREATE TABLE "public"."containers_volumes" (
    "chart_subcomponent_child_class_type_id" int8 NOT NULL REFERENCES chart_subcomponent_child_class_types(chart_subcomponent_child_class_type_id),
    "volume_id" int8 NOT NULL REFERENCES volumes(volume_id)
);
ALTER TABLE "public"."containers_volumes" ADD CONSTRAINT "containers_volumes_pk" UNIQUE ("chart_subcomponent_child_class_type_id", "volume_id");

------------
-- PROBES--
------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

-- tables for container_probes
CREATE TABLE "public"."container_probes" (
    "probe_id" int8 NOT NULL DEFAULT next_id(),
    "probe_key_values_jsonb" jsonb NOT NULL DEFAULT '{}'::jsonb
);
ALTER TABLE "public"."container_probes" ADD CONSTRAINT "container_probes_pk" PRIMARY KEY ("probe_id");

-- tables for containers_probes links
CREATE TABLE "public"."containers_probes" (
    "probe_id" int8  NOT NULL REFERENCES container_probes(probe_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id),
    "probe_type" text CHECK (probe_type IN ('livenessProbe', 'readinessProbe', 'startupProbe')) NOT NULL
);
ALTER TABLE "public"."containers_probes" ADD CONSTRAINT "containers_probes_pk" UNIQUE ("container_id", "probe_type");

--------------
-- RESOURCES--
--------------
---------------------------------------------------------------------------------------------------------------------------------------------------------------------

-- tables for resources
CREATE TABLE "public"."container_compute_resources" (
    "compute_resources_id" int8 NOT NULL DEFAULT next_id(),
    "compute_resources_cpu_request" text NOT NULL DEFAULT '',
    "compute_resources_cpu_limit" text NOT NULL DEFAULT '',
    "compute_resources_ram_request" text NOT NULL DEFAULT '',
    "compute_resources_ram_limit" text NOT NULL DEFAULT '',
    "compute_resources_ephemeral_storage_request" text NOT NULL DEFAULT '',
    "compute_resources_ephemeral_storage_limit" text NOT NULL DEFAULT ''
);
ALTER TABLE "public"."container_compute_resources" ADD CONSTRAINT "container_compute_resources_pk" PRIMARY KEY ("compute_resources_id");

-- tables for containers_compute_resources links
CREATE TABLE "public"."containers_compute_resources" (
    "compute_resources_id" int8 NOT NULL REFERENCES container_compute_resources(compute_resources_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id)
);
ALTER TABLE "public"."containers_compute_resources" ADD CONSTRAINT "containers_compute_resources_pk" UNIQUE ("container_id", "compute_resources_id");

---------------------------------------------------------------------------------------------------------------------------------------------------------------------
--------------
-- CmdArgs--
--------------
CREATE TABLE "public"."container_command_args" (
    "command_args_id" int8 NOT NULL DEFAULT next_id(),
    "command_values" text NOT NULL DEFAULT '',
    "args_values" text NOT NULL DEFAULT ''
);
ALTER TABLE "public"."container_command_args" ADD CONSTRAINT "container_command_args_pk" PRIMARY KEY ("command_args_id");

-- tables for container_command_args links
CREATE TABLE "public"."containers_command_args" (
    "command_args_id" int8 NOT NULL REFERENCES container_command_args(command_args_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id)
);
ALTER TABLE "public"."containers_command_args" ADD CONSTRAINT "containers_command_args_pk" UNIQUE ("container_id", "command_args_id");

---------------------------------------------------------------------------------------------------------------------------------------------------------------------
--------------
-- SecurityContext--
--------------
CREATE TABLE "public"."container_security_context" (
   "container_security_context_id" int8 NOT NULL DEFAULT next_id(),
   "security_context_key_values" text NOT NULL DEFAULT ''
);
ALTER TABLE "public"."container_security_context" ADD CONSTRAINT "container_security_context_pk" PRIMARY KEY ("container_security_context_id");

-- tables for container_command_args links
CREATE TABLE "public"."containers_security_context" (
    "container_security_context_id" int8 NOT NULL REFERENCES container_security_context(container_security_context_id),
    "container_id" int8 NOT NULL REFERENCES containers(container_id)
);
ALTER TABLE "public"."containers_security_context" ADD CONSTRAINT "containers_security_context_pk" UNIQUE ("container_id", "container_security_context_id");
