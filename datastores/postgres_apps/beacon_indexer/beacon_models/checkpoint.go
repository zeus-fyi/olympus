package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps"
)

type ValidatorsEpochCheckpoint struct {
	Epoch                       int
	ValidatorsActive            int
	ValidatorsBalancesRecorded  int
	ValidatorsBalancesRemaining int
}

func InsertEpochCheckpoint(ctx context.Context, epoch int) (int64, error) {
	log.Info().Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validators_balance_epoch, validators_active) VALUES (%d, (SELECT validators_active_at_epoch(%d)) )`, epoch, epoch)
	if epoch == 0 {
		query = fmt.Sprintf(`INSERT INTO validators_epoch_checkpoint (validators_balance_epoch, validators_active, validators_balances_recorded) VALUES (%d, (SELECT validators_active_at_epoch(%d)) , (SELECT validators_active_at_epoch(%d)))`, epoch, epoch, epoch)
	}
	r, err := postgres_apps.Pg.Exec(ctx, query)
	rowsAffected := r.RowsAffected()
	log.Info().Int64("rows affected: ", rowsAffected)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
		return rowsAffected, err
	}
	return rowsAffected, err
}

func (e *ValidatorsEpochCheckpoint) GetEpochCheckpoint(ctx context.Context, epoch int) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: InsertEpochCheckpoint")
	query := fmt.Sprintf(`SELECT validators_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM validators_epoch_checkpoint WHERE validators_balance_epoch = %d`, epoch)
	err := postgres_apps.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRecorded, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetEpochCheckpoint")
		return err
	}
	return err
}

func (e *ValidatorsEpochCheckpoint) GetFirstEpochCheckpointWithBalancesRemaining(ctx context.Context) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
	query := `SELECT validators_balance_epoch, validators_active, validators_balances_recorded, validators_balances_remaining FROM validators_epoch_checkpoint WHERE validators_balances_remaining <> 0 ORDER BY validators_balance_epoch LIMIT 1`
	err := postgres_apps.Pg.QueryRow(ctx, query).Scan(&e.Epoch, &e.ValidatorsActive, &e.ValidatorsBalancesRecorded, &e.ValidatorsBalancesRemaining)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
		return err
	}
	return err
}

func (e *ValidatorsEpochCheckpoint) GetCurrentFinalizedEpoch(ctx context.Context) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: GetCurrentFinalizedEpoch")
	query := "SELECT mainnet_finalized_epoch()"
	err := postgres_apps.Pg.QueryRow(ctx, query).Scan(&e.Epoch)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetCurrentFinalizedEpoch")
		return err
	}
	return err
}

func (e *ValidatorsEpochCheckpoint) GetNextEpochCheckpoint(ctx context.Context) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
	query := `SELECT COUNT(*) FROM validators_epoch_checkpoint`
	err := postgres_apps.Pg.QueryRow(ctx, query).Scan(&e.Epoch)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetFirstEpochCheckpointWithBalancesRemaining")
		return err
	}
	return err
}

func UpdateEpochCheckpointBalancesRecordedAtEpoch(ctx context.Context, epoch int) error {
	log.Info().Msgf("ValidatorsEpochCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch at Epoch %d", epoch)
	query := fmt.Sprintf(`SELECT update_checkpoint_at_epoch(%d)`, epoch)
	_, err := postgres_apps.Pg.Exec(ctx, query)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: UpdateEpochCheckpointBalancesRecordedAtEpoch")
		return err
	}
	return err
}

func SelectCountValidatorActive(ctx context.Context, epoch int) (int, error) {
	log.Info().Msg("SelectCountValidatorActive")
	var count int
	query := fmt.Sprintf(`SELECT validators_active_at_epoch(%d)`, epoch)
	err := postgres_apps.Pg.QueryRow(ctx, query).Scan(&count)
	log.Err(err).Msg("SelectCountValidatorEntries")
	if err != nil {
		return 0, err
	}
	return count, nil
}
