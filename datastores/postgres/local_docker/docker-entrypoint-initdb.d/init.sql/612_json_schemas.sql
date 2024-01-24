CREATE TABLE public.ai_json_schema_definitions(
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    schema_id BIGINT NOT NULL REFERENCES ai_schemas(schema_id),
    is_obj_array bool NOT NULL DEFAULT false,
    schema_name text NOT NULL,
    schema_group text NOT NULL DEFAULT 'default'
);
CREATE INDEX ai_json_schema_definitions_org_ind ON public.ai_json_schema_definitions(org_id);
CREATE INDEX ai_json_schema_name_ind ON public.ai_json_schema_definitions(schema_name);
CREATE INDEX ai_json_schema_gname_ind ON public.ai_json_schema_definitions(schema_group);
ALTER TABLE "public"."ai_json_schema_definitions" ADD CONSTRAINT "ai_json_schema_definitions_sid_sn_uniq" UNIQUE ("schema_id", "schema_name");
ALTER TABLE "public"."ai_json_schema_definitions" ADD CONSTRAINT "ai_json_schema_definitions_name_org_uniq" UNIQUE ("org_id", "schema_name");

CREATE TABLE public.ai_task_schemas(
    schema_id BIGINT NOT NULL REFERENCES ai_schemas(schema_id),
    task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id),
    PRIMARY KEY (schema_id, task_id)
);
CREATE INDEX idx_task_json_schema_id ON public.ai_task_schemas(task_id);
CREATE INDEX idx_ai_task_schemas_id ON public.ai_task_schemas(schema_id);
