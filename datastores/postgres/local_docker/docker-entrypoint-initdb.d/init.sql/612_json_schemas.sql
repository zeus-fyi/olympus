CREATE TABLE public.json_schema_definitions(
    schema_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id int8 NOT NULL REFERENCES orgs(org_id),
    is_obj_array bool NOT NULL DEFAULT false,
    schema_name text NOT NULL
);
CREATE INDEX json_schema_definitions_org_ind ON public.json_schema_definitions(org_id);
ALTER TABLE "public"."json_schema_definitions" ADD CONSTRAINT "json_org_schema_names_uniq" UNIQUE ("org_id", "schema_name");

CREATE TABLE public.ai_task_json_schema_fields(
    schema_id int8 NOT NULL REFERENCES json_schema_definitions(schema_id),
    field_name text NOT NULL,
    data_type text NOT NULL,
    PRIMARY KEY (schema_id, field_name)
);
CREATE INDEX idx_json_schema_field_id ON public.ai_task_json_schema_fields(schema_id);

CREATE TABLE public.ai_task_json_schemas(
    schema_id int8 NOT NULL REFERENCES json_schema_definitions(schema_id),
    task_id int8 NOT NULL REFERENCES ai_task_library(task_id),
    PRIMARY KEY (schema_id, task_id)
);

CREATE INDEX idx_task_json_schema_id ON public.ai_task_json_schemas(task_id);
