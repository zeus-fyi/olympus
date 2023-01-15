CREATE TABLE "public"."validator_service_org_group" (
    "group_name" text NOT NULL,
    "org_id" int8 NOT NULL REFERENCES orgs(org_id),
    "pubkey" text NOT NULL CHECK(LENGTH(pubkey)=98),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "fee_recipient" text NOT NULL,
    "enabled" bool NOT NULL DEFAULT false
);

ALTER TABLE "public"."validator_service_org_group" ADD CONSTRAINT "validators_org_group_pubkey_org_uniq" PRIMARY KEY ("org_id", "pubkey");
ALTER TABLE "public"."validator_service_org_group" ADD CONSTRAINT "validators_org_group_validator_pubkey_uniq" UNIQUE ("pubkey");
ALTER TABLE "public"."validator_service_org_group" ADD CONSTRAINT "validators_org_group_validator_pubkey_network_uniq" UNIQUE ("pubkey", "protocol_network_id");
CREATE INDEX "org_group_index" ON "public"."validator_service_org_group" ("group_name", "org_id", "protocol_network_id");
