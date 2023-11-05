-- network id 0 is all nodes
CREATE TABLE "public"."eth_p2p_nodes" (
    "id" int8 NOT NULL DEFAULT next_id(),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 0,
    "nodes" jsonb NOT NULL
);
ALTER TABLE "public"."eth_p2p_nodes" ADD CONSTRAINT "eth_p2p_nodes_pk" PRIMARY KEY ("id");

CREATE TABLE "public"."eth_mev_tx_analysis" (
    "tx_hash" text NOT NULL REFERENCES eth_mempool_mev_tx(tx_hash),
    "pair_address" text,
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

CREATE TABLE "public"."eth_mev_address_filter" (
    "address" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1
);
ALTER TABLE "public"."eth_mev_address_filter" ADD CONSTRAINT "eth_mev_address_filter_pk" PRIMARY KEY ("address");

CREATE TABLE "public"."eth_mev_bundle" (
    "event_id" int8 NOT NULL REFERENCES events (event_id),
    "bundle_hash" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1
);
ALTER TABLE "public"."eth_mev_bundle" ADD CONSTRAINT "eth_mev_bundle_pk" PRIMARY KEY ("bundle_hash");
CREATE INDEX eth_mev_bundle_protocol_id ON "public"."eth_mev_bundle" ("protocol_network_id");

CREATE TABLE eth_tx_receipts (
    tx_hash             text NOT NULL REFERENCES eth_tx (tx_hash),
    event_id            int8 NOT NULL REFERENCES events (event_id),
    status text NOT NULL,
    gas_used int8 NOT NULL,
    effective_gas_price int8 NOT NULL,
    cumulative_gas_used int8 NOT NULL,
    block_hash text NOT NULL,
    block_number int8 NOT NULL,
    transaction_index int8 NOT NULL
);
ALTER TABLE "public"."eth_tx_receipts" ADD CONSTRAINT "eth_tx_receipts_pk" PRIMARY KEY ("tx_hash");
CREATE INDEX eth_rx_block_number ON "public"."eth_tx_receipts" ("block_number" DESC);
CREATE INDEX eth_rx_status ON "public"."eth_tx_receipts" ("status");

CREATE TABLE "public"."eth_mev_bundle_profit" (
    "bundle_hash" text  NOT NULL REFERENCES eth_mev_bundle(bundle_hash),
    "revenue" int8 NOT NULL,
    "revenue_prediction" int8 NOT NULL DEFAULT 0,
    "revenue_prediction_skew" int8 NOT NULL GENERATED ALWAYS AS (revenue - revenue_prediction) STORED,
    "costs" int8 NOT NULL,
    "profit" int8 NOT NULL GENERATED ALWAYS AS (revenue - costs) STORED
);
ALTER TABLE "public"."eth_mev_bundle_profit" ADD CONSTRAINT "eth_mev_bundle_profit_pk" PRIMARY KEY ("bundle_hash");

CREATE TABLE "public"."eth_mev_call_bundle" (
    "bundle_hash" text  NOT NULL REFERENCES eth_mev_bundle(bundle_hash),
    "revenue" int8 NOT NULL,
    "revenue_prediction" int8 NOT NULL DEFAULT 0,
    "revenue_prediction_skew" int8 NOT NULL GENERATED ALWAYS AS (revenue - revenue_prediction) STORED,
    "costs" int8 NOT NULL,
    "profit" int8 NOT NULL GENERATED ALWAYS AS (revenue - costs) STORED
);
ALTER TABLE "public"."eth_mev_bundle_profit" ADD CONSTRAINT "eth_mev_bundle_profit_pk" PRIMARY KEY ("bundle_hash");

CREATE TABLE "public"."eth_mev_call_bundle" (
    "event_id" int8 NOT NULL REFERENCES events (event_id),
    "bundle_hash" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1
);
ALTER TABLE "public"."eth_mev_call_bundle" ADD CONSTRAINT "eth_mev_call_bundle_pk" PRIMARY KEY ("event_id", "bundle_hash");
CREATE INDEX eth_mev_call_bundle_protocol_id ON "public"."eth_mev_call_bundle" ("protocol_network_id");

CREATE TABLE "public"."eth_mev_block_builders" (
    "builder_name" text NOT NULL
);
ALTER TABLE "public"."eth_mev_block_builders" ADD CONSTRAINT "mev_block_builders_pk" PRIMARY KEY ("builder_name");


CREATE TABLE "public"."eth_mev_call_bundle" (
    "event_id" int8 NOT NULL REFERENCES events (event_id),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "bundle_hash" text NOT NULL,
    "builder_name" text NOT NULL REFERENCES eth_mev_block_builders(builder_name) NOT NULL,
    "eth_call_resp_json" jsonb NOT NULL DEFAULT '{}'::jsonb
);
ALTER TABLE "public"."eth_mev_call_bundle" ADD CONSTRAINT "eth_mev_call_bundle_pk" PRIMARY KEY ("event_id", "bundle_hash");
CREATE INDEX eth_mev_call_bundle_protocol_id ON "public"."eth_mev_call_bundle" ("protocol_network_id");
