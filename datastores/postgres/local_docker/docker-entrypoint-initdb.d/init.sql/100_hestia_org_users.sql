CREATE TABLE "public"."users" (
    "user_id" int8 NOT NULL DEFAULT next_id(),
    "email" text,
    "first_name" text,
    "last_name" text,
    "phone_number" text,
    "metadata" jsonb NOT NULL
);
ALTER TABLE "public"."users" ADD CONSTRAINT "user_pk" PRIMARY KEY ("user_id");
CREATE INDEX users_email_idx ON users (email);

CREATE TABLE "public"."orgs" (
  "org_id" int8 NOT NULL DEFAULT next_id(),
  "name" text NOT NULL,
  "metadata" jsonb
);
ALTER TABLE "public"."orgs" ADD CONSTRAINT "org_pk" PRIMARY KEY ("org_id");
ALTER TABLE "public"."orgs" ADD CONSTRAINT "org_name_uniq" UNIQUE ("name");

CREATE TABLE "public"."org_users" (
 "org_id" int8 NOT NULL REFERENCES orgs(org_id),
 "user_id" int8 NOT NULL REFERENCES users(user_id)
);
CREATE INDEX org_users_user_id_idx ON org_users(user_id);
CREATE INDEX org_users_org_id_idx ON org_users(org_id);