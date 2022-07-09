CREATE TABLE "public"."validator_balances_at_epoch" (
 "epoch" serial NOT NULL,
 "validator_index" int4 NOT NULL REFERENCES validators(index),
 "total_balance_gwei" int8 NOT NULL,
 "current_epoch_yield_gwei" int8 NOT NULL,
 "yield_to_date_gwei" int8 NOT NULL GENERATED ALWAYS AS (total_balance_gwei - 32000000000) STORED
)
;
CREATE INDEX scans_epoch_val_range_idx ON validator_balances_at_epoch (validator_index, epoch);

