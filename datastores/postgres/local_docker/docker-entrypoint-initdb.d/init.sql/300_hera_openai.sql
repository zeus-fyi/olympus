CREATE TABLE "public"."hera_openai_usage" (
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "tokens_remaining" int8 NOT NULL DEFAULT 0,
    "tokens_consumed" int8 NOT NULL DEFAULT 0
);
CREATE INDEX idx_oai_usage_org_cr ON public.hera_openai_usage("org_id");
ALTER TABLE "public"."hera_openai_usage" ADD CONSTRAINT "hera_openai_usage_pk" PRIMARY KEY ("org_id");

CREATE TABLE "public"."completion_responses" (
    response_id int8 NOT NULL DEFAULT next_id(),
    org_id int8 NOT NULL REFERENCES orgs(org_id),
    user_id int8 NOT NULL REFERENCES users(user_id),
    prompt_tokens int NOT NULL,
    completion_tokens int NOT NULL,
    total_tokens int NOT NULL,
    model text NOT NULL,
    completion_choices jsonb NOT NULL
);
ALTER TABLE "public"."completion_responses" ADD CONSTRAINT "completion_responses_pk" PRIMARY KEY ("response_id");
CREATE INDEX idx_ou_org_cr ON public.completion_responses("org_id");
CREATE INDEX idx_ou_user_cr ON public.completion_responses("user_id");
