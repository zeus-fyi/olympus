package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var NewValidatorBatchSize = 10
var NewValidatorTimeout = 10 * time.Minute

// FetchNewOrMissingValidators Routine ONE
func FetchNewOrMissingValidators() {
	log.Info().Msg("FetchNewOrMissingValidators")

	for {
		timeBegin := time.Now()
		err := fetchValidatorsToInsert(context.Background(), NewValidatorTimeout)
		log.Err(err)
		log.Info().Interface("FetchNewOrMissingValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(NewValidatorTimeout)
	}
}

func fetchValidatorsToInsert(ctx context.Context, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	err := Fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, NewValidatorBatchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}
