CREATE TABLE public.eval_fns(
    eval_id BIGINT PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    eval_name text NOT NULL,
    eval_type text NOT NULL,
    eval_group_name text NOT NULL,
    eval_model text,
    eval_format text NOT NULL
);

CREATE INDEX eval_fns_oid_ind ON public.eval_fns("org_id");
CREATE INDEX eval_fns_uid_ind ON public.eval_fns("user_id");
CREATE INDEX eval_fns_name_ind ON public.eval_fns("eval_name");
CREATE INDEX eval_fns_type_ind ON public.eval_fns("eval_type");

ALTER TABLE "public"."eval_fns" ADD CONSTRAINT "ai_eval_fns_name_uniq" UNIQUE ("org_id", "eval_name");
ALTER TABLE "public"."eval_fns" ADD CONSTRAINT "ai_eval_fns_group_name_uniq" UNIQUE ("org_id", "eval_group_name", "eval_name");

CREATE TABLE public.eval_metrics(
    eval_metric_id BIGINT PRIMARY KEY,
    eval_id BIGINT NOT NULL REFERENCES public.eval_fns(eval_id),
    eval_model_prompt text NOT NULL,
    eval_metric_name text NOT NULL,
    eval_metric_result text NOT NULL,
    eval_comparison_boolean boolean,
    eval_comparison_number float8,
    eval_comparison_string text,
    eval_metric_data_type text NOT NULL,
    eval_operator text NOT NULL,
    eval_state text NOT NULL
);
ALTER TABLE "public"."eval_metrics" ADD CONSTRAINT "eval_metrics_fn_uniq" UNIQUE ("eval_id", "eval_metric_id");

CREATE TABLE public.eval_metrics_results(
    eval_metrics_result_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    source_task_id int8 NOT NULL REFERENCES ai_task_library(task_id),
    eval_metric_id BIGINT NOT NULL REFERENCES public.eval_metrics(eval_metric_id),
    running_cycle_number int8 NOT NULL DEFAULT 1,
    search_window_unix_start int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    search_window_unix_end int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    eval_result_outcome boolean NOT NULL,
    eval_metadata jsonb
);

CREATE INDEX eval_result_outcome_idx ON public.eval_metrics_results("eval_result_outcome");
CREATE INDEX eval_result_metric_idx ON public.eval_metrics_results("eval_metric_id");
CREATE INDEX eval_result_orch_id_idx ON public.eval_metrics_results("orchestration_id");
CREATE INDEX eval_result_cycle_idx ON public.eval_metrics_results("running_cycle_number");
CREATE INDEX eval_result_source_search_start_idx ON public.eval_metrics_results("search_window_unix_start");
CREATE INDEX eval_result_source_search_end_idx ON public.eval_metrics_results("search_window_unix_end");

ALTER TABLE public.eval_metrics_results
ADD CONSTRAINT unique_eval_metrics_combination UNIQUE (eval_metric_id, source_task_id, orchestration_id, running_cycle_number);

CREATE TABLE public.ai_workflow_template_eval_task_relationships(
    task_eval_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    workflow_template_id BIGINT NOT NULL REFERENCES ai_workflow_template(workflow_template_id),
    task_id BIGINT NOT NULL REFERENCES ai_task_library(task_id),
    cycle_count BIGINT NOT NULL DEFAULT 1 CHECK ( cycle_count > 0 ),
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id)
);
ALTER TABLE "public"."ai_workflow_template_eval_task_relationships" ADD CONSTRAINT "ai_workflow_template_eval_task_relationships_uniq" UNIQUE ("workflow_template_id", "task_id", "eval_id");
CREATE INDEX ai_workflow_template_eval_task_relationships_wf_id ON public.ai_workflow_template_eval_task_relationships("workflow_template_id");
CREATE INDEX ai_workflow_template_eval_task_relationships_task_id ON public.ai_workflow_template_eval_task_relationships("task_id");
CREATE INDEX ai_workflow_template_eval_task_relationships_eval_id ON public.ai_workflow_template_eval_task_relationships("eval_id");
