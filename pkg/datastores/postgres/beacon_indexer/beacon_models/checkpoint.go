package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
)

type EpochCheckpoint struct {
	Epoch                       int
	ValidatorsActive            int
	ValidatorsBalancesRecorded  int
	ValidatorsBalancesRemaining int
}

func InsertEpochCheckpoint(ctx context.Context, epoch int) error {
	log.Info().Msg("EpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`INSERT INTO epoch_checkpoint (validator_balance_epoch, validators_active) VALUES (%d, (SELECT validators_active_at_epoch(%d)) )`, epoch, epoch)
	r, err := postgres.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Err(err).Msg("EpochCheckpoint: InsertEpochCheckpoint")
		return err
	}
	return err
}

func (e *EpochCheckpoint) GetEpochCheckpoint(ctx context.Context, epoch int) error {
	log.Info().Msg("EpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`SELECT validator_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM epoch_checkpoint WHERE validator_balance_epoch = %d`, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRemaining, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("EpochCheckpoint: GetEpochCheckpoint")
		return err
	}
	return err
}
