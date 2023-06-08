CREATE TABLE "public"."eth_mempool_mev_tx" (
    "tx_id" int8 NOT NULL DEFAULT next_id(),
    "tx_hash" text NOT NULL,
    "nonce" int8 NOT NULL DEFAULT 0,
    "from" text NOT NULL,
    "to" text NOT NULL,
    "block_number" int8 NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "tx" jsonb NOT NULL,
    "tx_flow_prediction" jsonb NOT NULL
);
ALTER TABLE "public"."eth_mempool_mev_tx" ADD CONSTRAINT "eth_mempool_mev_tx_pk" PRIMARY KEY ("tx_hash");
CREATE INDEX eth_mempool_mev_tx_ordering ON "public"."eth_mempool_mev_tx" ("block_number", "nonce" DESC);

-- network id 0 is all nodes
CREATE TABLE "public"."eth_p2p_nodes" (
  "id" int8 NOT NULL DEFAULT next_id(),
  "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 0,
  "nodes" jsonb NOT NULL
);
ALTER TABLE "public"."eth_p2p_nodes" ADD CONSTRAINT "eth_p2p_nodes_pk" PRIMARY KEY ("id");

CREATE TABLE "public"."erc20_token_info" (
    "address" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "balance_of_slot_num" int8 NOT NULL
);
ALTER TABLE "public"."erc20_token_info" ADD CONSTRAINT "erc20_token_info_pk" PRIMARY KEY ("address");

CREATE TABLE "public"."eth_mev_tx_analysis" (
    "tx_hash" text NOT NULL REFERENCES eth_mempool_mev_tx(tx_hash),
    "trade_method" text NOT NULL,
    "rx_block_number" int8 NOT NULL,
    "end_reason" text NOT NULL,
    "amount_in_addr" text NOT NULL,
    "amount_in" text NOT NULL,
    "amount_out_addr" text NOT NULL,
    "amount_out" text NOT NULL,
    "expected_profit_amount_out" text NOT NULL,
    "actual_profit_amount_out" text NOT NULL,
    "gas_used_wei" text NOT NULL,
    "metadata" jsonb NOT NULL
);
ALTER TABLE "public"."eth_mev_tx_analysis" ADD CONSTRAINT "eth_mev_tx_analysis_pk" PRIMARY KEY ("tx_hash");
CREATE INDEX eth_mev_tx_analysis_trade_method ON "public"."eth_mev_tx_analysis" ("trade_method");
CREATE INDEX eth_mev_tx_analysis_end_reason ON "public"."eth_mev_tx_analysis" ("end_reason");
CREATE INDEX eth_mev_tx_analysis_amount_in_addr ON "public"."eth_mev_tx_analysis" ("amount_in_addr");
CREATE INDEX eth_mev_tx_analysis_amount_out_addr ON "public"."eth_mev_tx_analysis" ("amount_out_addr");
