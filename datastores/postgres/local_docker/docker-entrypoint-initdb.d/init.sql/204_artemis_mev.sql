CREATE TABLE "public"."erc20_token_info" (
    "address" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "balance_of_slot_num" int8 NOT NULL,
    "name" text,
    "symbol" text,
    "decimals" int4,
    "transfer_tax_numerator" int8,
    "transfer_tax_denominator" int8,
    "trading_enabled" bool NOT NULL DEFAULT false,
);
ALTER TABLE "public"."erc20_token_info" ADD CONSTRAINT "erc20_token_info_pk" PRIMARY KEY ("address");

CREATE TABLE "public"."uniswap_pair_info" (
    "address" text NOT NULL,
    "factory_address" text NOT NULL,
    "fee" int8 NOT NULL,
    "version" text NOT NULL,
    "token0" text NOT NULL REFERENCES erc20_token_info(address),
    "token1" text NOT NULL REFERENCES erc20_token_info(address),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id) DEFAULT 1,
    "trading_enabled" bool NOT NULL DEFAULT false
);
ALTER TABLE "public"."uniswap_pair_info" ADD CONSTRAINT "uniswap_pair_info_pk" PRIMARY KEY ("address");

CREATE TABLE "public"."events" (
    "event_id" int8 NOT NULL DEFAULT next_id()
);
ALTER TABLE "public"."events" ADD CONSTRAINT "events_pk" PRIMARY KEY ("event_id");

CREATE TABLE "public"."tags"
(
    "tag_id" int8 NOT NULL DEFAULT next_id()
);
ALTER TABLE "public"."tags" ADD CONSTRAINT "tags_pk" PRIMARY KEY ("tag_id");

CREATE TABLE "public"."event_tags"
(
    "event_group_id" int8 NOT NULL DEFAULT next_id(),
    "event_id"       int8 NOT NULL REFERENCES events (event_id),
    "tag_id"         int8 NOT NULL REFERENCES tags (tag_id),
);
CREATE INDEX event_tags_index ON "public"."event_tags" ("event_group_id" DESC);
ALTER TABLE "public"."event_tags" ADD CONSTRAINT "event_tags_pk" PRIMARY KEY ("event_id");

CREATE TABLE "public"."eth_tx"
(
    "event_id"            int8 NOT NULL REFERENCES events (event_id),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks (protocol_network_id) DEFAULT 1,
    "tx_hash"             text NOT NULL
);
ALTER TABLE "public"."eth_tx" ADD CONSTRAINT "eth_tx_pk" PRIMARY KEY ("tx_hash");

CREATE TABLE "public"."eth_rx"
(
    "event_id"            int8 NOT NULL REFERENCES events (event_id),
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks (protocol_network_id) DEFAULT 1,
    "rx_hash"             text NOT NULL
);
ALTER TABLE "public"."eth_rx" ADD CONSTRAINT "eth_rx_pk" PRIMARY KEY ("rx_hash");

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

