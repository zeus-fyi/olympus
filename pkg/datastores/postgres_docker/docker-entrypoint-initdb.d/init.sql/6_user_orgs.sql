CREATE TABLE "public"."users" (
    "user_id" int8 NOT NULL,
    "metadata" jsonb NOT NULL
);
ALTER TABLE "public"."users" ADD CONSTRAINT "user_pk" PRIMARY KEY ("user_id");

CREATE TABLE "public"."orgs" (
  "org_id" int8 NOT NULL,
  "metadata" jsonb NOT NULL
);
ALTER TABLE "public"."orgs" ADD CONSTRAINT "org_pk" PRIMARY KEY ("org_id");

CREATE TABLE "public"."user_orgs" (
 "org_id" int8 NOT NULL REFERENCES orgs(org_id),
 "user_id" int8 NOT NULL REFERENCES users(user_id)
);

