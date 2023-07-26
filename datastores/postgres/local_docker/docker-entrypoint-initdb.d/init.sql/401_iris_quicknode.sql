CREATE TABLE "public"."provisioned_quicknode_services" (
    org_id int8 NOT NULL REFERENCES orgs(org_id),
    quicknode_id text NOT NULL,
    endpoint_id text NOT NULL,
    wss_url text,
    http_url text,
    chain text,
    network text,
    plan text NOT NULL CHECK (plan IN ('lite', 'standard', 'performance')),
    active bool NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (quicknode_id)
);
ALTER TABLE "public"."provisioned_quicknode_services" ADD CONSTRAINT "org_qid_uniq" UNIQUE ("org_id", "quicknode_id");

CREATE TABLE "public"."provisioned_quicknode_services_referers" (
    quicknode_id text NOT NULL REFERENCES provisioned_quicknode_services(quicknode_id),
    referer text NOT NULL,
    PRIMARY KEY (quicknode_id)
);
ALTER TABLE "public"."provisioned_quicknode_services_referers" ADD CONSTRAINT "org_qid_referer_uniq" UNIQUE ("quicknode_id", "referer");

CREATE TABLE "public"."provisioned_quicknode_services_contract_addresses" (
    quicknode_id text NOT NULL REFERENCES provisioned_quicknode_services(quicknode_id),
    contract_address text NOT NULL,
    PRIMARY KEY (quicknode_id)
);
ALTER TABLE "public"."provisioned_quicknode_services_contract_addresses" ADD CONSTRAINT "org_qid_contract_addr_uniq" UNIQUE ("quicknode_id", "contract_address");

CREATE TRIGGER set_timestamp_on_provisioned_quicknode_services
AFTER UPDATE ON provisioned_quicknode_services
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();