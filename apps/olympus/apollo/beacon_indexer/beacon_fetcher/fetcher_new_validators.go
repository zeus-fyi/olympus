package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var NewValidatorBatchSize = 10000
var NewValidatorTimeout = 5 * time.Second

// FetchNewOrMissingValidators Routine ONE
func FetchNewOrMissingValidators() {
	log.Info().Msg("FetchNewOrMissingValidators")

	for {
		timeBegin := time.Now()
		err := fetchValidatorsToInsert(context.Background(), NewValidatorBatchSize, NewValidatorTimeout)
		log.Err(err)
		log.Info().Interface("FetchNewOrMissingValidators took this many seconds to complete: ", time.Now().Sub(timeBegin))
		time.Sleep(NewValidatorTimeout)
	}
}

func fetchValidatorsToInsert(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()
	err := Fetcher.BeaconFindNewAndMissingValidatorIndexes(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}
