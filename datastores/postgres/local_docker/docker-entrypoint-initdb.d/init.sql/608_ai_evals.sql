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


CREATE TABLE public.eval_metrics(
    eval_metric_id BIGINT PRIMARY KEY,
    eval_id BIGINT NOT NULL REFERENCES public.eval_fns(eval_id),
    eval_model_prompt text NOT NULL,
    eval_metric_name text NOT NULL,
    eval_metric_result text NOT NULL,
    eval_comparison_boolean boolean,
    eval_comparison_number BIGINT,
    eval_comparison_string text,
    eval_metric_data_type text NOT NULL,
    eval_operator text NOT NULL,
    eval_state text NOT NULL
);
ALTER TABLE "public"."eval_metrics" ADD CONSTRAINT "eval_metrics_fn_uniq" UNIQUE ("eval_id", "eval_metric_id");
