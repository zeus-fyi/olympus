CREATE TABLE public.ai_task_library (
    task_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    max_tokens_per_task BIGINT NOT NULL DEFAULT 0,
    task_type TEXT NOT NULL,
    task_name TEXT NOT NULL,
    task_group TEXT NOT NULL DEFAULT 'default',
    token_overflow_strategy TEXT NOT NULL DEFAULT 'deduce',
    model TEXT NOT NULL,
    prompt TEXT NOT NULL
);

ALTER TABLE "public"."ai_task_library" ADD CONSTRAINT "ai_task_library_org_task_group_name_uniq" UNIQUE ("org_id", "task_group", "task_name");
ALTER TABLE "public"."ai_task_library" ADD CONSTRAINT "ai_task_library_org_group_tn_names_uniq" UNIQUE ("org_id", "task_name", "task_type");
ALTER TABLE "public"."ai_task_library" ADD CONSTRAINT "ai_task_library_org_group_tt_names_uniq" UNIQUE ("org_id", "task_group", "task_type", "task_name");
CREATE INDEX ai_task_library_task_type_idx ON public.ai_task_library("task_type");
CREATE INDEX ai_task_library_task_group_idx ON public.ai_task_library("task_group");
CREATE INDEX ai_task_library_org_idx ON public.ai_task_library("org_id");
CREATE INDEX ai_task_library_user_idx ON public.ai_task_library("user_id");

CREATE TABLE public.ai_workflow_template(
    workflow_template_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    workflow_name TEXT NOT NULL,
    fundamental_period BIGINT NOT NULL DEFAULT 1 CHECK ( fundamental_period > 0 ),
    fundamental_period_time_unit TEXT NOT NULL CHECK (fundamental_period_time_unit IN ('seconds', 'minutes', 'hours', 'days', 'weeks', 'months', 'years'))
);
CREATE INDEX workflow_name_index ON public.ai_workflow_template (workflow_name);

CREATE TABLE public.ai_workflow_component(
    component_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY
);

CREATE TABLE public.ai_workflow_template_components(
    workflow_template_id BIGINT NOT NULL REFERENCES ai_workflow_template(workflow_template_id),
    component_id BIGINT NOT NULL REFERENCES ai_workflow_component(component_id)
);
CREATE INDEX workflow_component_index ON public.ai_workflow_template_components(component_id);

CREATE TABLE public.ai_workflow_template_component_task(
    component_id BIGINT NOT NULL REFERENCES ai_workflow_component(component_id) PRIMARY KEY,
    task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id),
    cycle_count BIGINT NOT NULL DEFAULT 1 CHECK (cycle_count > 0 )
);
CREATE INDEX ai_workflow_template_component_task_idx ON public.ai_workflow_template_component_task(task_id);

CREATE TABLE public.ai_workflow_component_dependency(
    component_id BIGINT NOT NULL REFERENCES ai_workflow_component(component_id) PRIMARY KEY,
    component_dependency_id BIGINT NOT NULL REFERENCES ai_workflow_component(component_id)
);
ALTER TABLE "public"."ai_workflow_component_dependency" ADD CONSTRAINT "ai_workflow_component_dependency_uniq" UNIQUE ("component_id", "component_dependency_id");

CREATE TABLE public.ai_retrieval_library (
    retrieval_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    retrieval_name TEXT NOT NULL,
    retrieval_group TEXT NOT NULL,
    instructions jsonb NOT NULL
);
CREATE INDEX ai_retrieval_library_inst_idx ON public.ai_retrieval_library USING GIN (instructions);
CREATE INDEX ai_retrieval_library_name_idx ON public.ai_retrieval_library("retrieval_name");
CREATE INDEX ai_retrieval_library_group_idx ON public.ai_retrieval_library("retrieval_group");
CREATE INDEX ai_retrieval_library_org_idx ON public.ai_retrieval_library("org_id");
CREATE INDEX ai_retrieval_library_user_idx ON public.ai_retrieval_library("user_id");
ALTER TABLE "public"."ai_retrieval_library" ADD CONSTRAINT "ai_retrieval_library_org_gret_name_uniq" UNIQUE ("org_id", "retrieval_name");

CREATE TABLE public.ai_workflow_template_component_retrieval(
    component_id BIGINT NOT NULL REFERENCES ai_workflow_component(component_id) PRIMARY KEY,
    retrieval_id BIGINT NOT NULL REFERENCES ai_retrieval_library(retrieval_id)
);
CREATE INDEX ai_workflow_template_component_retrieval_idx ON public.ai_workflow_template_component_retrieval(retrieval_id);
