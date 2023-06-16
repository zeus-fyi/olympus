CREATE TABLE "public"."org_secret_references" (
    "secret_id" int8 NOT NULL DEFAULT next_id(),
    "secret_name" text NOT NULL,
    "org_id" int8 NOT NULL REFERENCES orgs(org_id)
);

ALTER TABLE "public"."org_secret_references" ADD CONSTRAINT "org_secret_references_pk" PRIMARY KEY ("secret_id");
ALTER TABLE "public"."org_secret_references" ADD CONSTRAINT "uniq_secret_name" UNIQUE ("secret_name","org_id");

CREATE TABLE "public"."org_secret_key_val_references" (
    "secret_id" int8 NOT NULL REFERENCES org_secret_references(secret_id),
    "secret_env_var_ref" text NOT NULL,
    "secret_key_ref" text NOT NULL,
    "secret_name_ref" text NOT NULL
);
ALTER TABLE "public"."org_secret_key_val_references" ADD CONSTRAINT "org_secret_key_val_references_pk" PRIMARY KEY ("secret_id");
ALTER TABLE "public"."org_secret_key_val_references" ADD CONSTRAINT "uniq_secret_name_ref" UNIQUE ("secret_key_ref","secret_name_ref","secret_id");

CREATE TABLE "public"."topology_system_components_secrets" (
    "topology_system_component_id" int8 NOT NULL REFERENCES topology_system_components (topology_system_component_id),
    "secret_id" int8 NOT NULL REFERENCES org_secret_references (secret_id)
);
ALTER TABLE "public"."topology_system_components_secrets" ADD CONSTRAINT "topology_system_components_secrets_pk" PRIMARY KEY ("topology_system_component_id");
