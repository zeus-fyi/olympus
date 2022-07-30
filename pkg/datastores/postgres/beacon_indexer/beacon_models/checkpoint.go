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
	query := fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validators_balance_epoch, validators_active) VALUES (%d, (SELECT validators_active_at_epoch(%d)) )`, epoch, epoch)
	if epoch == 0 {
		query = fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validators_balance_epoch, validators_active, validators_balances_recorded) VALUES (%d, (SELECT validators_active_at_epoch(%d)) , (SELECT validators_active_at_epoch(%d)))`, epoch, epoch, epoch)
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
	query := fmt.Sprintf(`SELECT validators_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM validators_epoch_checkpoint WHERE validators_balance_epoch = %d`, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRecorded, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetEpochCheckpoint")
		return err
	}
	return err
}

func (e *ValidatorsEpochCheckpoint) GetFirstEpochCheckpointWithBalancesRemaining(ctx context.Context) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
	query := `SELECT validators_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM validators_epoch_checkpoint WHERE validators_balances_remaining <> 0 ORDER BY validators_balance_epoch LIMIT 1`
	err := postgres.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRecorded, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
		return err
	}
	return err
}

func UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx context.Context, epoch int) error {
	log.Info().Msgf("ValidatorsEpochCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch at Epoch %d", epoch)
	query := fmt.Sprintf(`SELECT update_checkpoint_at_epoch(%d)`, epoch)
	_, err := postgres.Pg.Exec(ctx, query)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch")
		return err
	}
	return err
}

func SelectCountValidatorActive(ctx context.Context, epoch int) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT validators_active_at_epoch(%d)`, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&count)
	log.Err(err).Msg("SelectCountValidatorEntries")
	if err != nil {
		return 0, err
	}
	return count, nil
}
