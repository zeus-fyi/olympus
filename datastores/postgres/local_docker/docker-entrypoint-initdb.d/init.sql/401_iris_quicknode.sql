CREATE TABLE "public"."quicknode_marketplace_customer" (
    quicknode_id text NOT NULL,
    plan text NOT NULL CHECK (plan IN ('test', 'free', 'lite', 'standard', 'performance','enterprise')),
    is_test bool NOT NULL DEFAULT false,
    tutorial_on bool NOT NULL DEFAULT true,
    PRIMARY KEY (quicknode_id)
);

ALTER TABLE "public"."quicknode_marketplace_customer" ADD CONSTRAINT "mp_org_qid_uniq2" UNIQUE ("quicknode_id", "plan");
CREATE INDEX quicknode_marketplace_customer_test_users ON quicknode_marketplace_customer ("is_test");

CREATE TABLE "public"."provisioned_quicknode_services" (
    quicknode_id text NOT NULL,
    endpoint_id text NOT NULL,
    wss_url text,
    http_url text,
    chain text,
    network text,
    active bool NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    PRIMARY KEY (quicknode_id, endpoint_id)
);
ALTER TABLE "public"."provisioned_quicknode_services" ADD CONSTRAINT "org_qid_uniq" UNIQUE ("quicknode_id", "endpoint_id");
CREATE UNIQUE INDEX qn_endpoint_ind ON provisioned_quicknode_services ("endpoint_id");
CREATE INDEX provisioned_quicknode_services_http_ind ON provisioned_quicknode_services ("http_url");

CREATE TABLE "public"."provisioned_quicknode_services_referers" (
    endpoint_id text NOT NULL,
    referer text NOT NULL,
    PRIMARY KEY (endpoint_id),
    FOREIGN KEY (endpoint_id) REFERENCES provisioned_quicknode_services(endpoint_id)
);
CREATE INDEX endpoint_ref_ind ON provisioned_quicknode_services_referers ("endpoint_id");

CREATE TABLE "public"."provisioned_quicknode_services_contract_addresses" (
    endpoint_id text NOT NULL,
    contract_address text NOT NULL,
    PRIMARY KEY (endpoint_id),
    FOREIGN KEY (endpoint_id) REFERENCES provisioned_quicknode_services(endpoint_id)
);
CREATE INDEX endpoint_ca_ind ON provisioned_quicknode_services_contract_addresses ("endpoint_id");

CREATE TRIGGER set_timestamp_on_provisioned_quicknode_services
AFTER UPDATE ON provisioned_quicknode_services
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();