CREATE TABLE public.ai_json_schema_definitions(
    schema_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    is_obj_array bool NOT NULL DEFAULT false,
    schema_name text NOT NULL,
    schema_group text NOT NULL DEFAULT 'default'
);
CREATE INDEX ai_json_schema_name_ind ON public.ai_json_schema_definitions(schema_name);
CREATE INDEX ai_json_schema_gname_ind ON public.ai_json_schema_definitions(schema_group);
CREATE INDEX ai_json_schema_definitions_org_ind ON public.ai_json_schema_definitions(org_id);
ALTER TABLE "public"."ai_json_schema_definitions" ADD CONSTRAINT "ai_json_schema_definitions_name_org_uniq" UNIQUE ("org_id", "schema_name");

CREATE TABLE public.ai_json_schema_fields(
    schema_id BIGINT NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
    field_name text NOT NULL,
    field_description text NOT NULL,
    data_type text NOT NULL,
    PRIMARY KEY (schema_id, field_name)
);
CREATE INDEX idx_json_schema_field_id ON public.ai_json_schema_fields(schema_id);

CREATE TABLE public.ai_json_task_schemas(
    schema_id BIGINT NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
    task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id),
    PRIMARY KEY (schema_id, task_id)
);

CREATE INDEX idx_task_json_schema_id ON public.ai_json_task_schemas(task_id);

CREATE TABLE public.ai_json_eval_schemas(
    schema_id BIGINT NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    PRIMARY KEY (schema_id, eval_id)
);

CREATE INDEX idx_eval_id_json_schema_id ON public.ai_json_eval_schemas(eval_id);

-- Redefine the ai_json_eval_schemas table
CREATE TABLE public.ai_json_eval_metric_schemas (
     eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
     schema_id BIGINT NOT NULL REFERENCES ai_json_schema_definitions(schema_id),
     field_name text NOT NULL,
     eval_metric_id BIGINT NOT NULL REFERENCES eval_metrics(eval_metric_id),
     PRIMARY KEY (eval_id, schema_id, field_name),
     FOREIGN KEY (schema_id, field_name) REFERENCES ai_json_schema_fields(schema_id, field_name)
);

-- Create an index for the eval_id field
CREATE INDEX idx_eval_id_on_ai_json_eval_metric_schemas ON public.ai_json_eval_metric_schemas(eval_id);
