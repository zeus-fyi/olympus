CREATE TABLE "public"."hera_openai_usage" (
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "tokens_remaining" int8 NOT NULL DEFAULT 0,
    "tokens_consumed" int8 NOT NULL DEFAULT 0
);

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
