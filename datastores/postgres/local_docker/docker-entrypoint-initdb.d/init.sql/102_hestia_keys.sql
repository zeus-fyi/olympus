CREATE TABLE "public"."key_types" (
    "key_type_id" int8 NOT NULL DEFAULT next_id(),
    "key_type_name" text NOT NULL
);
ALTER TABLE "public"."key_types" ADD CONSTRAINT "key_type_pk" PRIMARY KEY ("key_type_id");

CREATE TABLE "public"."users_keys" (
    "public_key" text NOT NULL,
    "user_id" int8 NOT NULL REFERENCES users(user_id),
    "public_key_name" text NOT NULL NOT NULL DEFAULT '',
    "public_key_verified" bool NOT NULL DEFAULT false,
    "public_key_type_id" int8 NOT NULL REFERENCES key_types(key_type_id),
    "created_at" timestamptz  NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."users_keys" ADD CONSTRAINT "users_keys_pk" PRIMARY KEY ("public_key");

CREATE TABLE "public"."users_key_groups" (
    "user_id" int8 NOT NULL REFERENCES users(user_id),
    "public_key" text NOT NULL REFERENCES users_keys(public_key),
    "key_group_id" int8 NOT NULL DEFAULT next_id(),
    "key_group_name" text NOT NULL DEFAULT '',
    "updated_at" timestamptz  NOT NULL DEFAULT NOW()
);
ALTER TABLE "public"."users_key_groups" ADD CONSTRAINT "user_key_pk" PRIMARY KEY ("key_group_id");

CREATE TRIGGER set_timestamp_on_user_key_groups
    BEFORE UPDATE ON users_key_groups
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
