CREATE OR REPLACE FUNCTION update_or_insert_validator_status()
    RETURNS trigger AS
$$
DECLARE
    current_epoch int8 := (SELECT mainnet_head_epoch());
    FAR_FUTURE_EPOCH numeric := 2^64-1;
    status validator_status := 'unknown';
    substatus validator_substatus := 'unknown';
BEGIN
    IF validators.activation_eligibility_epoch > current_epoch THEN
        status := 'pending';
        CASE
            WHEN (validators.activation_eligibility_epoch == FAR_FUTURE_EPOCH) THEN
                substatus := 'pending_initialized';
            WHEN (validators.activation_eligibility_epoch < FAR_FUTURE_EPOCH) AND (validators.activation_epoch > current_epoch) THEN
                substatus := 'pending_queued';
            ELSE
                substatus := 'unknown';
            END CASE;
    ELSIF validators.activation_eligibility_epoch <= current_epoch AND current_epoch < validators.exit_epoch THEN
        status := 'active';
        CASE
            WHEN (validators.activation_epoch <= current_epoch) AND (validators.exit_epoch == FAR_FUTURE_EPOCH) THEN
                substatus := 'active_ongoing';
            WHEN (validators.activation_epoch <= current_epoch) AND (current_epoch < validators.exit_epoch AND validators.exit_epoch < FAR_FUTURE_EPOCH) AND (NOT validators.slashed) THEN
                substatus := 'active_exiting';
            WHEN (validators.activation_epoch <= current_epoch) AND (current_epoch < validators.exit_epoch AND validators.exit_epoch < FAR_FUTURE_EPOCH) AND validators.slashed THEN
                substatus := 'active_slashed';
            ELSE
                substatus := 'unknown';
            END CASE;
    ELSIF validators.exit_epoch <= current_epoch AND current_epoch < validators.withdrawable_epoch THEN
        status := 'exited';
        CASE
            WHEN (validators.exit_epoch <= current_epoch AND current_epoch < validators.withdrawable_epoch) AND (NOT validators.slashed) THEN
                substatus := 'exited_unslashed';
            WHEN (validators.exit_epoch <= current_epoch AND current_epoch < validators.withdrawable_epoch) AND validators.slashed THEN
                substatus := 'exited_slashed';
            ELSE
                substatus := 'unknown';
            END CASE;
    ELSIF validators.withdrawable_epoch <= current_epoch THEN
        status := 'withdrawal';
        CASE
            WHEN (validators.withdrawable_epoch <= current_epoch) AND (validators.balance != 0) THEN
                substatus := 'withdrawal_possible';
            WHEN (validators.withdrawable_epoch <= current_epoch) AND (validators.balance == 0) THEN
                substatus := 'withdrawal_done';
            ELSE
                substatus := 'unknown';
            END CASE;
    ELSE
        status := 'unknown';
    END IF;

    IF (TG_OP = 'INSERT') AND (NEW.activation_epoch IS NOT NULL) THEN
        CASE
            WHEN NEW.activation_epoch < FAR_FUTURE_EPOCH THEN
                INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei)
                VALUES (NEW.activation_epoch, NEW.index, 32000000000, 0);
            ELSE
            END CASE;
    END IF;

    IF (TG_OP = 'UPDATE') AND (NEW.activation_epoch IS NOT NULL) THEN
        CASE
            WHEN (OLD.activation_epoch IS NULL) AND (NEW.activation_epoch < FAR_FUTURE_EPOCH) THEN
                INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei)
                VALUES (NEW.activation_epoch, NEW.index, 32000000000, 0);
            WHEN (NEW.activation_epoch < FAR_FUTURE_EPOCH) AND (OLD.activation_epoch != NEW.activation_epoch) THEN
                INSERT INTO validator_balances_at_epoch (epoch, validator_index, total_balance_gwei, current_epoch_yield_gwei)
                VALUES (NEW.activation_epoch, NEW.index, 32000000000, 0);
            ELSE
            END CASE;
    END IF;

    NEW.status=status;
    NEW.substatus=substatus;

    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE TRIGGER trigger_update_validator_status
AFTER INSERT OR UPDATE OF balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed
ON validators
FOR EACH ROW EXECUTE PROCEDURE update_or_insert_validator_status();