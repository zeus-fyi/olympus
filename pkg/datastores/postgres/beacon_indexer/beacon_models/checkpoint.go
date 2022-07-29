package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
)

type ValidatorsEpochCheckpoint struct {
	Epoch                       int
	ValidatorsActive            int
	ValidatorsBalancesRecorded  int
	ValidatorsBalancesRemaining int
}

func InsertEpochCheckpoint(ctx context.Context, epoch int) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validator_balance_epoch, validators_active) VALUES (%d, (SELECT validators_active_at_epoch(%d)) )`, epoch, epoch)
	if epoch == 0 {
		query = fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validator_balance_epoch, validators_active, validators_balances_recorded) VALUES (%d, (SELECT validators_active_at_epoch(%d)) , (SELECT validators_active_at_epoch(%d)))`, epoch, epoch, epoch)
	}
	r, err := postgres.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
		return err
	}
	return err
}

func (e *ValidatorsEpochCheckpoint) GetEpochCheckpoint(ctx context.Context, epoch int) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`SELECT validator_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM validators_epoch_checkpoint WHERE validator_balance_epoch = %d`, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRecorded, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetEpochCheckpoint")
		return err
	}
	return err
}

func SelectCountValidatorActive(ctx context.Context, epoch int) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM validators WHERE validators.activation_epoch <= %d AND %d < validators.exit_epoch`, epoch, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&count)
	log.Err(err).Msg("SelectCountValidatorEntries")
	if err != nil {
		return 0, err
	}
	return count, nil
}
