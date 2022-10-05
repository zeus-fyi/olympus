package beacon_models

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

var insertValidatorsFromBeaconAPI = `INSERT INTO validators (index, pubkey, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed, withdrawal_credentials) VALUES `

func (vs *Validators) InsertValidatorsFromBeaconAPI(ctx context.Context) error {
	log.Info().Msg("Validators: InsertValidatorsFromBeaconAPI")

	vs.RowSetting.RowsToInclude = "beacon_state"
	querySuffix := ` ON CONFLICT ON CONSTRAINT validators_pkey DO UPDATE SET index = EXCLUDED.index, pubkey = EXCLUDED.pubkey `
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
	if len(vs.Validators) == 0 {
		return 0, errors.New("no Validators were supplied")
	}
	validators := string_utils.DelimitedSliceStrBuilderSQLRows("", vs.GetManyRowValues())
	query := fmt.Sprintf(`
	UPDATE validators SET
		balance = x.balance,
		effective_balance = x.effective_balance,
		activation_eligibility_epoch = x.activation_eligibility_epoch,
		activation_epoch = x.activation_epoch,
		exit_epoch = x.exit_epoch,
		withdrawable_epoch = x.withdrawable_epoch,
		slashed = x.slashed
	FROM (VALUES %s ) AS x(index, balance, effective_balance, activation_eligibility_epoch, activation_epoch, exit_epoch, withdrawable_epoch, slashed)
    WHERE x.index = validators.index `, validators)

	rows, err := postgres.Pg.Exec(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("UpdateValidatorsFromBeaconAPI: Query")
		return rows.RowsAffected(), err
	}

	return rows.RowsAffected(), nil
}
