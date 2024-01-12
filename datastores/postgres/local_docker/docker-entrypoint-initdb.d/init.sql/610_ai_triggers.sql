CREATE TABLE public.ai_trigger_actions(
    trigger_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    trigger_name text NOT NULL,
    trigger_group text NOT NULL,
    trigger_env text NOT NULL DEFAULT 'social-media-io-text'
);

CREATE INDEX ai_trigger_actions_oid_ind ON public.ai_trigger_actions("org_id");
CREATE INDEX ai_trigger_actions_uid_ind ON public.ai_trigger_actions("user_id");
CREATE INDEX ai_trigger_actions_name_ind ON public.ai_trigger_actions("trigger_name");
CREATE INDEX ai_trigger_actions_group_ind ON public.ai_trigger_actions("trigger_group");

ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_name_uniq" UNIQUE ("org_id", "trigger_name");
ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_group_names_uniq" UNIQUE ("org_id", "trigger_group", "trigger_name");

CREATE TABLE public.ai_trigger_actions_evals(
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    eval_trigger_state text NOT NULL,
    eval_results_trigger_on text NOT NULL,
    PRIMARY KEY (eval_id, trigger_id)
);
ALTER TABLE public.ai_trigger_actions_evals
ADD CONSTRAINT unique_eval_triggers UNIQUE (eval_id, trigger_id);

CREATE TABLE public.ai_trigger_actions_approval(
    approval_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    workflow_result_id BIGINT NOT NULL REFERENCES ai_workflow_analysis_results(workflow_result_id),
    approval_state text NOT NULL DEFAULT 'pending',
    request_summary text NOT NULL,
    updated_at timestamptz  NOT NULL DEFAULT NOW()
);

CREATE INDEX ai_trigger_actions_approval_eval_source ON public.ai_trigger_actions_approval("eval_id");
CREATE INDEX ai_trigger_actions_approval_trigger_id ON public.ai_trigger_actions_approval("trigger_id");
CREATE INDEX ai_trigger_actions_approval_trigger_state ON public.ai_trigger_actions_approval("approval_state");

CREATE INDEX ai_trigger_actions_approval_wf_analysis_source ON public.ai_trigger_actions_approval("workflow_result_id");
ALTER TABLE public.ai_trigger_actions_approval
    ADD CONSTRAINT unique_eval_triggers_actions UNIQUE (eval_id, trigger_id, workflow_result_id);

CREATE TRIGGER set_timestamp_on_trigger_actions_approval
BEFORE UPDATE ON ai_trigger_actions_approval
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
