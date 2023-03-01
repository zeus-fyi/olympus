CREATE TABLE "public"."eth_scheduled_delivery" (
    "delivery_id" int8 NOT NULL DEFAULT next_id(),
    "public_key" text NOT NULL REFERENCES users_keys(public_key),
    "delivery_schedule_type" text NOT NULL,
    "protocol_network_id" int8 NOT NULL REFERENCES protocol_networks(protocol_network_id),
    "amount" int8 NOT NULL,
    "units" text NOT NULL
);
ALTER TABLE "public"."eth_scheduled_delivery" ADD CONSTRAINT "eth_scheduled_delivery_pk" PRIMARY KEY ("delivery_id");
