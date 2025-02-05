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
    prompt TEXT NOT NULL,
    response_format text NOT NULL DEFAULT 'text',
    temperature FLOAT8 NOT NULL DEFAULT 1.0,
    margin_buffer FLOAT8 NOT NULL DEFAULT 0.5,
    CONSTRAINT temperature_range CHECK (temperature >= 0 AND temperature <= 2),
    CONSTRAINT margin_buffer_range CHECK (margin_buffer >= 0.2 AND margin_buffer <= 0.8)
);

ALTER TABLE "public"."ai_task_library" ADD CONSTRAINT "ai_task_library_org_task_group_name_uniq" UNIQUE ("org_id", "task_group", "task_name");
ALTER TABLE "public"."ai_task_library" ADD CONSTRAINT "ai_task_library_org_group_tt_names_uniq" UNIQUE ("org_id", "task_group", "task_type", "task_name");
CREATE INDEX ai_task_library_task_type_idx ON public.ai_task_library("task_type");
CREATE INDEX ai_task_library_task_group_idx ON public.ai_task_library("task_group");
CREATE INDEX ai_task_library_org_idx ON public.ai_task_library("org_id");
CREATE INDEX ai_task_library_user_idx ON public.ai_task_library("user_id");

CREATE TABLE public.ai_workflow_template(
    workflow_template_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    workflow_name TEXT NOT NULL,
    workflow_group TEXT NOT NULL,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    fundamental_period BIGINT NOT NULL DEFAULT 1 CHECK ( fundamental_period > 0 ),
    fundamental_period_time_unit TEXT NOT NULL CHECK (fundamental_period_time_unit IN ('seconds', 'minutes', 'hours', 'days', 'weeks', 'months', 'years'))
);
CREATE INDEX workflow_name_index ON public.ai_workflow_template (workflow_name);
CREATE INDEX workflow_gname_index ON public.ai_workflow_template (workflow_group);
CREATE INDEX ai_workflow_template_org_idx ON public.ai_workflow_template("org_id");
CREATE INDEX ai_workflow_template_user_idx ON public.ai_workflow_template("user_id");
ALTER TABLE "public"."ai_workflow_template" ADD CONSTRAINT "ai_workflow_template_org_name_uniq" UNIQUE ("org_id", "workflow_name");
ALTER TABLE "public"."ai_workflow_template" ADD CONSTRAINT "ai_workflow_template_org_gname_uniq" UNIQUE ("org_id", "workflow_group", "workflow_name");

CREATE TABLE public.ai_retrieval_library (
    retrieval_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    retrieval_name TEXT NOT NULL,
    retrieval_group TEXT NOT NULL,
    retrieval_platform TEXT NOT NULL,
    instructions jsonb NOT NULL
);
CREATE INDEX ai_retrieval_library_inst_idx ON public.ai_retrieval_library USING GIN (instructions);
CREATE INDEX ai_retrieval_library_name_idx ON public.ai_retrieval_library("retrieval_name");
CREATE INDEX ai_retrieval_library_group_idx ON public.ai_retrieval_library("retrieval_group");
CREATE INDEX ai_retrieval_library_platform_idx ON public.ai_retrieval_library("retrieval_platform");
CREATE INDEX ai_retrieval_library_org_idx ON public.ai_retrieval_library("org_id");
CREATE INDEX ai_retrieval_library_user_idx ON public.ai_retrieval_library("user_id");
ALTER TABLE "public"."ai_retrieval_library" ADD CONSTRAINT "ai_retrieval_library_org_gret_name_uniq" UNIQUE ("org_id", "retrieval_name");

CREATE TABLE public.ai_workflow_template_analysis_tasks (
    workflow_template_id BIGINT NOT NULL REFERENCES ai_workflow_template(workflow_template_id),
    task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id),
    cycle_count BIGINT NOT NULL DEFAULT 1 CHECK (cycle_count > 0),
    retrieval_id BIGINT REFERENCES ai_retrieval_library(retrieval_id) -- Allows NULL
);

ALTER TABLE public.ai_workflow_template_analysis_tasks
    ADD PRIMARY KEY (workflow_template_id, task_id);

CREATE UNIQUE INDEX idx_uniq_retrieval_id
    ON public.ai_workflow_template_analysis_tasks (workflow_template_id, task_id, retrieval_id)
    WHERE retrieval_id IS NOT NULL;

ALTER TABLE public.ai_workflow_template_analysis_tasks
    ADD CONSTRAINT unique_workflow_task_retrieval
        UNIQUE (workflow_template_id, task_id, retrieval_id);

CREATE INDEX ai_workflow_template_analysis_tasks_idx ON public.ai_workflow_template_analysis_tasks (workflow_template_id);
CREATE INDEX ai_workflow_template_analysis_tasks_idx2 ON public.ai_workflow_template_analysis_tasks (task_id);
CREATE INDEX ai_workflow_template_analysis_tasks_idx3 ON public.ai_workflow_template_analysis_tasks (retrieval_id);

CREATE TABLE public.ai_workflow_template_agg_tasks(
    workflow_template_id BIGINT NOT NULL REFERENCES ai_workflow_template(workflow_template_id),
    agg_task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id) CHECK (agg_task_id != analysis_task_id),
    analysis_task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id) CHECK (agg_task_id != analysis_task_id),
    cycle_count BIGINT NOT NULL DEFAULT 1 CHECK ( cycle_count > 0 )
);
ALTER TABLE "public"."ai_workflow_template_agg_tasks" ADD CONSTRAINT "ai_workflow_template_agg_tasks_link_uniq" UNIQUE ("workflow_template_id", "agg_task_id", "analysis_task_id");

CREATE INDEX ai_workflow_template_agg_tasks_idx ON public.ai_workflow_template_agg_tasks (workflow_template_id);
CREATE INDEX ai_workflow_template_agg_tasks_idx2 ON public.ai_workflow_template_agg_tasks (agg_task_id);
CREATE INDEX ai_workflow_template_agg_tasks_idx3 ON public.ai_workflow_template_agg_tasks (analysis_task_id);

ALTER TABLE public.ai_workflow_template_agg_tasks
DROP CONSTRAINT ai_workflow_template_agg_tasks_link_uniq,
ADD PRIMARY KEY (workflow_template_id, agg_task_id, analysis_task_id);
