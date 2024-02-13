CREATE TABLE public.ai_trigger_actions(
    trigger_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    trigger_name text NOT NULL,
    trigger_group text NOT NULL,
    trigger_action text NOT NULL DEFAULT 'social-media-engagement',
    expires_after_seconds BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX ai_trigger_actions_oid_ind ON public.ai_trigger_actions("org_id");
CREATE INDEX ai_trigger_actions_uid_ind ON public.ai_trigger_actions("user_id");
CREATE INDEX ai_trigger_actions_name_ind ON public.ai_trigger_actions("trigger_name");
CREATE INDEX ai_trigger_actions_group_ind ON public.ai_trigger_actions("trigger_group");

ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_name_uniq" UNIQUE ("org_id", "trigger_name");
ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_group_names_uniq" UNIQUE ("org_id", "trigger_group", "trigger_name");

CREATE TABLE public.ai_trigger_eval(
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    eval_trigger_state text NOT NULL,
    eval_results_trigger_on text NOT NULL,
    PRIMARY KEY (trigger_id)
);
CREATE INDEX ai_trigger_eval_trigger_index ON public.ai_trigger_eval("trigger_id");

CREATE TABLE public.ai_trigger_actions_evals(
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    PRIMARY KEY (eval_id, trigger_id)
);
CREATE INDEX ai_trigger_actions_evals_indx ON public.ai_trigger_actions_evals("eval_id");
CREATE INDEX ai_trigger_actions_trg_indx ON public.ai_trigger_actions_evals("trigger_id");

ALTER TABLE public.ai_trigger_actions_evals
    ADD CONSTRAINT eval_trigger_uniq UNIQUE (eval_id, trigger_id);


CREATE TABLE public.ai_trigger_actions_approvals(
    approval_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    workflow_result_id BIGINT NOT NULL REFERENCES ai_workflow_analysis_results(workflow_result_id),
    approval_state text NOT NULL DEFAULT 'pending',
    request_summary text NOT NULL,
    expires_at timestamptz,
    updated_at timestamptz  NOT NULL DEFAULT NOW()
);
CREATE INDEX ai_trigger_actions_approval_eval_source ON public.ai_trigger_actions_approvals("eval_id");
CREATE INDEX ai_trigger_actions_approval_trigger_id ON public.ai_trigger_actions_approvals("trigger_id");
CREATE INDEX ai_trigger_actions_approval_trigger_state ON public.ai_trigger_actions_approvals("approval_state");
CREATE INDEX ai_trigger_actions_approval_wf_analysis_source ON public.ai_trigger_actions_approvals("workflow_result_id");
ALTER TABLE public.ai_trigger_actions_approvals
    ADD CONSTRAINT unique_eval_triggers_actions UNIQUE (approval_id, eval_id, trigger_id, workflow_result_id);

CREATE TRIGGER set_timestamp_on_trigger_actions_approval
BEFORE UPDATE ON ai_trigger_actions_approvals
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
