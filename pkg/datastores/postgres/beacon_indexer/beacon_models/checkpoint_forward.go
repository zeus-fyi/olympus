package beacon_models

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
)

func (e *ValidatorsEpochCheckpoint) GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch(ctx context.Context, epoch int) error {
	log.Info().Msg("ValidatorsEpochCheckpoint: GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch")
	query := fmt.Sprintf(`SELECT validators_balance_epoch WHERE validators_balances_remaining <> 0 AND validators_balance_epoch > %d ORDER BY validators_balance_epoch LIMIT 1`, epoch)
	err := postgres.Pg.QueryRow(ctx, query).Scan(&e.Epoch)
	if err != nil {
		log.Err(err).Msg("ValidatorsEpochCheckpoint: GetAnyEpochCheckpointWithBalancesRemainingAfterEpoch")
		return err
	}
	return err
}
