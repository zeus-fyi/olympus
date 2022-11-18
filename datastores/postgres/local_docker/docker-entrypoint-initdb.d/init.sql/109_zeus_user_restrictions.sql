CREATE TABLE "public"."topologies_org_cloud_ctx_ns" (
    "cloud_ctx_ns_id" int8 NOT NULL DEFAULT next_id(),
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "cloud_provider" text NOT NULL,
    "context" text NOT NULL,
    "region" text NOT NULL,
    "namespace" text NOT NULL,
    "created_at" timestamptz  NOT NULL DEFAULT NOW()
);

ALTER TABLE "public"."topologies_org_cloud_ctx_ns" ADD CONSTRAINT "cloud_ctx_ns_id_pk" PRIMARY KEY ("cloud_ctx_ns_id");
ALTER TABLE "public"."topologies_org_cloud_ctx_ns" ADD CONSTRAINT "cloud_ctx_ns_unique" UNIQUE ("cloud_provider", "context", "region", "namespace");
