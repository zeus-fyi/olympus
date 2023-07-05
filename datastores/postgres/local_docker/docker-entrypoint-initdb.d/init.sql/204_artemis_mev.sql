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

CREATE TABLE "public"."eth_mempool_mev_tx" (
    "tx_id" int8 NOT NULL DEFAULT next_id(),
    "tx_hash" text NOT NULL,
    "pair_address" text,
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

