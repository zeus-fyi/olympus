CREATE TABLE "public"."eth_ecdsa_rxs" (
    "rx_id" int8 NOT NULL DEFAULT next_id(),
    "tx_id" int8 NOT NULL REFERENCES eth_ecdsa_txs(tx_id),
    "gas_used_gwei" int8 DEFAULT 0,
    "gas_used_gwei_decimals" int8 DEFAULT 0,
    "gas_used_cumulative_gwei" int8 DEFAULT 0,
    "gas_used_cumulative_gwei_decimals" int8 DEFAULT 0,
    "block_number" int8,
    "block_timestamp" timestamptz,
    "tx_index" int,
    "contract_address" text,
    "status" text
);
ALTER TABLE "public"."eth_ecdsa_rxs" ADD CONSTRAINT "eth_rx_id_pk" PRIMARY KEY ("rx_id");
ALTER TABLE "public"."eth_ecdsa_rxs" ADD CONSTRAINT "eth_rx_tx_id" UNIQUE ("tx_id");