package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var UpdateAllValidatorTimeout = time.Minute * 10

// UpdateAllValidators Routine TWO
func UpdateAllValidators() {
	log.Info().Msg("UpdateAllValidators")
	for {
		timeBegin := time.Now()
		err := fetchAllValidatorsToUpdate(context.Background(), UpdateAllValidatorTimeout)
		log.Err(err)
		log.Info().Interface("UpdateAllValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(UpdateAllValidatorTimeout)
	}
}

func fetchAllValidatorsToUpdate(ctx context.Context, contextTimeout time.Duration) error {
	log.Info().Msg("fetchAllValidatorsToUpdate")
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := fetcher.BeaconUpdateAllValidatorStates(ctxTimeout)
	log.Info().Err(err).Msg("UpdateAllValidators: fetchAllValidatorsToUpdate")
	return err
}
