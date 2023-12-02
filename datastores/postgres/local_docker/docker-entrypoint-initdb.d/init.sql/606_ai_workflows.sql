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
