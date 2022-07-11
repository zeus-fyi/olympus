package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/databases/postgres"
)

func SelectValidatorsToQueryBeaconForBalanceUpdates(ctx context.Context, batchSize int) (ValidatorBalancesEpoch, error) {
	log.Info().Msg("SelectValidatorsToQueryBeaconForBalanceUpdates")

	var selectedValidatorBalances ValidatorBalancesEpoch
	query := fmt.Sprintf(`SELECT MAX(epoch) as max_epoch, validator_index
								 FROM validator_balances_at_epoch
								 WHERE epoch + 1 < (SELECT mainnet_finalized_epoch())
					             GROUP by validator_index
								 LIMIT %d`, batchSize)

	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return selectedValidatorBalances, err
	}
	for rows.Next() {
		var vb ValidatorBalanceEpoch
		rowErr := rows.Scan(&vb.Epoch, &vb.Index)
		if rowErr != nil {
			return selectedValidatorBalances, rowErr
		}
		selectedValidatorBalances.ValidatorBalance = append(selectedValidatorBalances.ValidatorBalance, vb)
	}
	return selectedValidatorBalances, nil
}

func FindValidatorIndexes(ctx context.Context, batchSize int) (Validators, error) {
	query := fmt.Sprintf(`
	SELECT
	generate_series FROM GENERATE_SERIES(
		(SELECT MIN(index) from validators), (SELECT MAX(index)+%d from validators)
	)
	WHERE NOT EXISTS(SELECT index FROM validators WHERE index = generate_series)`, batchSize)

	var validatorsToQueryState Validators
	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return validatorsToQueryState, err
	}
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index)
		if rowErr != nil {
			return validatorsToQueryState, rowErr
		}
		validatorsToQueryState.Validators = append(validatorsToQueryState.Validators, validator)
	}
	return validatorsToQueryState, err
}

func SelectValidatorsQueryOngoingStates(ctx context.Context, batchSize int) (Validators, error) {
	query := fmt.Sprintf(`
	SELECT index FROM validators ORDER BY updated_at LIMIT %d `, batchSize)

	var validatorsToQueryState Validators
	rows, err := postgres.Pg.Query(ctx, query)
	if err != nil {
		return validatorsToQueryState, err
	}
	for rows.Next() {
		var validator Validator
		rowErr := rows.Scan(&validator.Index)
		if rowErr != nil {
			return validatorsToQueryState, rowErr
		}
		validatorsToQueryState.Validators = append(validatorsToQueryState.Validators, validator)
	}
	return validatorsToQueryState, err
}

func SelectValidatorsToQueryBalancesByEpoch(ctx context.Context, batchSize int) (map[int64][]ValidatorBalanceEpoch, error) {
	epochMap := make(map[int64][]ValidatorBalanceEpoch, 1)
	vbal, err := SelectValidatorsToQueryBeaconForBalanceUpdates(ctx, batchSize)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SelectValidatorsToQueryBalancesURL: had nil ValidatorBalancesEpoch")
		return epochMap, err
	}
	for _, vb := range vbal.ValidatorBalance {
		epochMap[vb.Epoch] = append(epochMap[vb.Epoch], vb)
	}
	return epochMap, err
}

func SelectValidatorsToQueryBalancesURL(ctx context.Context, batchSize int) (string, error) {
	vbal, err := SelectValidatorsToQueryBeaconForBalanceUpdates(ctx, batchSize)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("SelectValidatorsToQueryBalancesURL: had nil ValidatorBalancesEpoch")
		return "", err
	}
	return vbal.FormatValidatorBalancesEpochIndexesToURLList(), nil
}

func FindNewValidatorsToQueryBeaconURLEncoded(ctx context.Context, batchSize int) (string, error) {
	log.Ctx(ctx).Info().Msg("FindNewValidatorsToQueryBeaconURLEncoded")

	vs, err := FindValidatorIndexes(ctx, batchSize)
	if err != nil {
		return "", err
	}
	return vs.formatValidatorStateIndexesToURLList(), nil
}

func SelectValidatorsQueryOngoingStatesIndexesURLEncoded(ctx context.Context, batchSize int) (string, error) {
	log.Ctx(ctx).Info().Msg("SelectValidatorsQueryOngoingStatesIndexesURLEncoded")

	vs, err := SelectValidatorsQueryOngoingStates(ctx, batchSize)
	if err != nil {
		return "", err
	}
	return vs.formatValidatorStateIndexesToURLList(), nil
}
