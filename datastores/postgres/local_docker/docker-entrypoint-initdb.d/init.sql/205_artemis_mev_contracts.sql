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
