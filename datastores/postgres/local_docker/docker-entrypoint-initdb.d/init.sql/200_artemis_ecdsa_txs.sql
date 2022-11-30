CREATE TABLE "public"."eth_ecdsa_txs" (
    "tx_id" int8 NOT NULL DEFAULT next_id(),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "nonce" int8 NOT NULL DEFAULT 0,
    "gas_price_gwei" int8,
    "gas_price_gwei_decimals" int8 DEFAULT 0,
    "gas_limit_gwei" int8 DEFAULT 21000,
    "gas_limit_gwei_decimals" int8 DEFAULT 0,
    "amount_gwei" int8 DEFAULT 0,
    "amount_gwei_decimals" int8 DEFAULT 0,
    "public_key_type_id" int8 NOT NULL REFERENCES key_types(key_type_id) CHECK (public_key_type_id IN (5)),
    "public_key" text NOT NULL REFERENCES users_keys(public_key),
    "r" text NOT NULL,
    "s" text NOT NULL,
    "v" text NOT NULL,
    "to" text,
    "tx_hash" text,
    "payload" bytea
);
ALTER TABLE "public"."eth_ecdsa_txs" ADD CONSTRAINT "eth_tx_id_pk" PRIMARY KEY ("tx_id");
ALTER TABLE "public"."eth_ecdsa_txs" ADD CONSTRAINT "key_network_nonce_unique" UNIQUE ("public_key", "protocol_network_id", "nonce");
CREATE INDEX eth_ecdsa_nonce_tx_hash_nulls ON "public"."eth_ecdsa_txs" ("nonce", "tx_hash" DESC NULLS LAST);


