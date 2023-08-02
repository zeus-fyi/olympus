CREATE EXTENSION pgcrypto;

CREATE TABLE "public"."key_types" (
    "key_type_id" int8 NOT NULL DEFAULT next_id(),
    "key_type_name" text NOT NULL
);
ALTER TABLE "public"."key_types" ADD CONSTRAINT "key_type_pk" PRIMARY KEY ("key_type_id");

CREATE TABLE "public"."users_keys" (
    "user_id" int8 NOT NULL REFERENCES users(user_id),
    "public_key_type_id" int8 NOT NULL REFERENCES key_types(key_type_id),
    "created_at" timestamptz  NOT NULL DEFAULT NOW(),
    "public_key_name" text NOT NULL NOT NULL DEFAULT '',
    "public_key" text NOT NULL,
    "public_key_verified" bool NOT NULL DEFAULT false
);
ALTER TABLE "public"."users_keys" ADD CONSTRAINT "users_keys_pk" PRIMARY KEY ("public_key");
CREATE INDEX ON users_keys (public_key, user_id);
CREATE INDEX users_keys_public_key_type_id_idx ON users_keys (public_key_type_id);
CREATE INDEX users_keys_user_id_idx ON users_keys (user_id);

CREATE TABLE "public"."users_key_groups" (
    "user_id" int8 NOT NULL REFERENCES users(user_id),
    "key_group_id" int8 NOT NULL DEFAULT next_id(),
    "updated_at" timestamptz  NOT NULL DEFAULT NOW(),
    "public_key" text NOT NULL REFERENCES users_keys(public_key),
    "key_group_name" text NOT NULL DEFAULT ''
);
ALTER TABLE "public"."users_key_groups" ADD CONSTRAINT "user_key_pk" PRIMARY KEY ("key_group_id");

CREATE TRIGGER set_timestamp_on_user_key_groups
    BEFORE UPDATE ON users_key_groups
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "public"."services" (
    "service_id" int8 NOT NULL DEFAULT next_id(),
    "service_name" text NOT NULL
);
ALTER TABLE "public"."services" ADD CONSTRAINT "services_pk" PRIMARY KEY ("service_id");
ALTER TABLE "public"."services" ADD CONSTRAINT "service_name_uniq" UNIQUE ("service_name");

CREATE TABLE "public"."users_key_services" (
    "service_id" int8 NOT NULL REFERENCES services(service_id),
    "public_key" text NOT NULL REFERENCES users_keys(public_key)
);
ALTER TABLE "public"."users_key_services" ADD CONSTRAINT "users_key_services_pk" PRIMARY KEY ("service_id", "public_key");
