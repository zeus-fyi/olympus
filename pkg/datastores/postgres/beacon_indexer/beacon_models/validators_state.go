package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var insertValidatorsFromBeaconAPI = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed, withdrawal_credentials) VALUES `

func (vs *Validators) InsertValidatorsFromBeaconAPI(ctx context.Context) error {
	log.Info().Msg("Validators: InsertValidatorsFromBeaconAPI")

	vs.RowSetting.RowsToInclude = "beacon_state"

	querySuffix := ` ON CONFLICT ON CONSTRAINT validators_pkey DO NOTHING`
	query := string_utils.PrefixAndSuffixDelimitedSliceStrBuilderSQLRows(insertValidatorsFromBeaconAPI, vs.GetManyRowValues(), querySuffix)
	r, err := postgres.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Error().Err(err).Msg("InsertValidatorsFromBeaconAPI: InsertValidatorsFromBeaconAPI")
		return err
	}
	return err
}

func (vs *Validators) UpdateValidatorsFromBeaconAPI(ctx context.Context) (int64, error) {
	vs.RowSetting.RowsToInclude = "beacon_state_update"
	validators := string_utils.MultiArraySliceStrBuilderSQL(vs.GetManyRowValues())
	query := fmt.Sprintf(`
	WITH validator_update AS (
		SELECT * FROM UNNEST(%s) AS x(index, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed)
	) 
	UPDATE validators
	SET balance = validator_update.balance, effective_balance = validator_update.effective_balance, activation_eligibility_epoch = validator_update.activation_eligibility_epoch, activation_epoch = validator_update.activation_epoch, exit_epoch = validator_update.exit_epoch, withdrawable_epoch = validator_update.withdrawable_epoch, slashed = validator_update.slashed
	JOIN validators ON validator_update.index = validators.index`, validators)

	rows, err := postgres.Pg.Exec(ctx, query)
	if err != nil {
		log.Error().Err(err).Msgf("UpdateValidatorsFromBeaconAPI: Query %s", query)
		return rows.RowsAffected(), err
	}

	return rows.RowsAffected(), nil
}
