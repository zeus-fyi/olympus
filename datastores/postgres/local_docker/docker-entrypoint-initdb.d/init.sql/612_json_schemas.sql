CREATE TABLE public.ai_json_schema_definitions(
    schema_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id int8 NOT NULL REFERENCES orgs(org_id),
    is_obj_array bool NOT NULL DEFAULT false,
    schema_name text NOT NULL,
    schema_group text NOT NULL DEFAULT 'default'
);
CREATE INDEX ai_json_schema_name_ind ON public.ai_json_schema_definitions(schema_name);
CREATE INDEX ai_json_schema_gname_ind ON public.ai_json_schema_definitions(schema_group);
CREATE INDEX ai_json_schema_definitions_org_ind ON public.ai_json_schema_definitions(org_id);
ALTER TABLE "public"."ai_json_schema_definitions" ADD CONSTRAINT "ai_json_schema_definitions_name_org_uniq" UNIQUE ("org_id", "schema_name");

CREATE TABLE public.ai_json_schema_fields(
    schema_id int8 NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
    field_name text NOT NULL,
    field_description text NOT NULL,
    data_type text NOT NULL,
    PRIMARY KEY (schema_id, field_name)
);
CREATE INDEX idx_json_schema_field_id ON public.ai_json_schema_fields(schema_id);

CREATE TABLE public.ai_json_task_schemas(
    schema_id int8 NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
    task_id int8 NOT NULL REFERENCES ai_task_library(task_id),
    PRIMARY KEY (schema_id, task_id)
);

CREATE INDEX idx_task_json_schema_id ON public.ai_json_task_schemas(task_id);
