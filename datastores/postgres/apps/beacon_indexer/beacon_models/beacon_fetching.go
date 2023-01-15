package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func SelectValidatorsToQueryBeaconForBalanceUpdates(ctx context.Context, batchSize int) (ValidatorBalancesEpoch, error) {
	log.Info().Msg("SelectValidatorsToQueryBeaconForBalanceUpdates")

	var selectedValidatorBalances ValidatorBalancesEpoch
	query := fmt.Sprintf(`SELECT MAX(epoch) AS max_epoch, MAX(epoch)+1 AS next_epoch, 32*(MAX(epoch)+1) AS next_epoch_slot, validator_index
								 FROM validator_balances_at_epoch
								 WHERE epoch + 1 < (SELECT mainnet_finalized_epoch())
					             GROUP by validator_index
								 LIMIT %d`, batchSize)
	log.Debug().Interface("SelectValidatorsToQueryBeaconForBalanceUpdates: Query: ", query)

	rows, err := apps.Pg.Query(ctx, query)
	log.Err(err).Interface("SelectValidatorsToQueryBeaconForBalanceUpdates: Query: ", query)
	if err != nil {
		return selectedValidatorBalances, err
	}
	defer rows.Close()
	for rows.Next() {
		var vb ValidatorBalanceEpoch
		rowErr := rows.Scan(&vb.Epoch, &vb.NextEpochToQuery, &vb.NextSlotToQuery, &vb.Index)
		if rowErr != nil {
			log.Err(rowErr).Interface("SelectValidatorsToQueryBeaconForBalanceUpdates: Query: ", query)
			return selectedValidatorBalances, rowErr
		}
		selectedValidatorBalances.ValidatorBalances = append(selectedValidatorBalances.ValidatorBalances, vb)
	}
	log.Info().Interface("SelectValidatorsToQueryBeaconForBalanceUpdates: selectedValidatorBalances: ", selectedValidatorBalances)
	return selectedValidatorBalances, nil
}

func FindValidatorIndexes(ctx context.Context, batchSize, networkID int) (Validators, error) {
	log.Info().Msg("FindValidatorIndexes")
	query := fmt.Sprintf(`
	SELECT
	generate_series FROM GENERATE_SERIES(
		(SELECT COALESCE(MIN(index),0) from validators WHERE protocol_network_id = $1), (SELECT COALESCE(MAX(index)+%d,+%d) from validators WHERE protocol_network_id = $1)
	)
	WHERE NOT EXISTS(SELECT index FROM validators WHERE index = generate_series AND protocol_network_id = $1) LIMIT %d`, batchSize, batchSize, batchSize)

	var validatorsToQueryState Validators
	log.Debug().Interface("FindValidatorIndexes: Query: ", query)
	rows, err := apps.Pg.Query(ctx, query, networkID)
	if err != nil {
		return validatorsToQueryState, err
	}
	defer rows.Close()
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index)
		if rowErr != nil {
			log.Err(rowErr).Interface("FindValidatorIndexes: Query: ", query)
			return validatorsToQueryState, rowErr
		}
		validatorsToQueryState.Validators = append(validatorsToQueryState.Validators, validator)
	}
	log.Err(err).Interface("FindValidatorIndexes: Query: ", query)
	return validatorsToQueryState, err
}

func SelectValidatorsQueryOngoingStates(ctx context.Context, batchSize, networkID int) (Validators, error) {
	log.Info().Msg("SelectValidatorsQueryOngoingStates")

	query := fmt.Sprintf(`
	SELECT index FROM validators WHERE protocol_network_id = $1 ORDER BY updated_at LIMIT %d `, batchSize)
	log.Debug().Interface("SelectValidatorsQueryOngoingStates: Query: ", query)

	var validatorsToQueryState Validators
	rows, err := apps.Pg.Query(ctx, query, networkID)
	if err != nil {
		log.Err(err).Msg("SelectValidatorsQueryOngoingStates")
		return validatorsToQueryState, err
	}
	defer rows.Close()
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index)
		if rowErr != nil {
			log.Err(err).Interface("SelectValidatorsQueryOngoingStates: Query: ", query)
			return validatorsToQueryState, rowErr
		}
		validatorsToQueryState.Validators = append(validatorsToQueryState.Validators, validator)
	}
	log.Err(err).Interface("SelectValidatorsQueryOngoingStates: Query: ", query)
	return validatorsToQueryState, err
}

func FindNewValidatorsToQueryBeaconURLEncoded(ctx context.Context, batchSize, networkID int) (string, error) {
	log.Ctx(ctx).Info().Msg("FindNewValidatorsToQueryBeaconURLEncoded")

	vs, err := FindValidatorIndexes(ctx, batchSize, networkID)
	if err != nil {
		return "", err
	}
	return vs.formatValidatorStateIndexesToURLList(), nil
}

func SelectValidatorsQueryOngoingStatesIndexesURLEncoded(ctx context.Context, batchSize, networkID int) (string, error) {
	log.Ctx(ctx).Info().Msg("SelectValidatorsQueryOngoingStatesIndexesURLEncoded")

	vs, err := SelectValidatorsQueryOngoingStates(ctx, batchSize, networkID)
	if err != nil {
		return "", err
	}
	return vs.formatValidatorStateIndexesToURLList(), nil
}
