CREATE TABLE "public"."validators_epoch_checkpoint" (
     "validator_balance_epoch" serial NOT NULL,
     "validators_active" int4 NOT NULL,
     "validators_balances_recorded" int4 NOT NULL DEFAULT 0,
     "validators_balances_remaining" int4 NOT NULL GENERATED ALWAYS AS (validators_active - validators_balances_recorded) STORED
)
;
ALTER TABLE "public"."validators_epoch_checkpoint" ADD CONSTRAINT "validator_balances_remaining_at_epoch_pkey" PRIMARY KEY ("validator_balance_epoch","validators_balances_remaining");

CREATE OR REPLACE FUNCTION "public"."validators_active_at_epoch"(epoch int4)
RETURNS "pg_catalog"."int4" AS
$BODY$
DECLARE
BEGIN
    RETURN (SELECT COUNT(*) FROM validators WHERE validator.activation_epoch <= epoch AND epoch < validator.exit_epoch);
END;
$BODY$
LANGUAGE plpgsql VOLATILE
COST 100;