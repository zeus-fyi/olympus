CREATE TABLE "public"."validator_balances_at_epoch" (
 "epoch" serial NOT NULL,
 "validator_index" int4 NOT NULL REFERENCES validators(index),
 "total_balance_gwei" int8 NOT NULL,
 "current_epoch_yield_gwei" int8 NOT NULL,
 "yield_to_date_gwei" int8 NOT NULL GENERATED ALWAYS AS (total_balance_gwei - 32000000000) STORED
)
;
ALTER TABLE "public"."validator_balances_at_epoch" ADD CONSTRAINT "validator_balances_at_epoch_pkey" PRIMARY KEY ("validator_index","epoch");
ALTER TABLE "public"."validator_balances_at_epoch" ADD CONSTRAINT "validator_balances_at_epoch_validator_index_fkey" FOREIGN KEY ("validator_index") REFERENCES "public"."validators" ("index") ON DELETE NO ACTION ON UPDATE NO ACTION;
