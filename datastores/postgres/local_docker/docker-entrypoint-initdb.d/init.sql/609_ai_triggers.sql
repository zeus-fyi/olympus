CREATE TABLE public.ai_trigger_actions(
    trigger_id BIGINT NOT NULL DEFAULT next_id() PRIMARY KEY,
    org_id BIGINT NOT NULL REFERENCES orgs(org_id),
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    trigger_name text NOT NULL,
    trigger_group text NOT NULL
);

CREATE INDEX ai_trigger_actions_oid_ind ON public.ai_trigger_actions("org_id");
CREATE INDEX ai_trigger_actions_uid_ind ON public.ai_trigger_actions("user_id");
CREATE INDEX ai_trigger_actions_name_ind ON public.ai_trigger_actions("trigger_name");
CREATE INDEX ai_trigger_actions_group_ind ON public.ai_trigger_actions("trigger_group");

ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_name_uniq" UNIQUE ("org_id", "trigger_name");
ALTER TABLE "public"."ai_trigger_actions" ADD CONSTRAINT "ai_trigger_actions_group_names_uniq" UNIQUE ("org_id", "trigger_group", "trigger_name");

CREATE TABLE public.ai_eval_trigger_actions(
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    trigger_id BIGINT NOT NULL REFERENCES ai_trigger_actions(trigger_id),
    eval_trigger_state text NOT NULL,
    eval_results_trigger_on text NOT NULL,
    PRIMARY KEY (eval_id, trigger_id)
);
ALTER TABLE public.ai_eval_trigger_actions
ADD CONSTRAINT unique_eval_triggers UNIQUE (eval_id, trigger_id);
