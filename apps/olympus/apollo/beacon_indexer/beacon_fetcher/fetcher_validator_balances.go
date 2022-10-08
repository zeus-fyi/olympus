package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

var NewValidatorBalancesBatchSize = 1000
var NewValidatorBalancesTimeout = time.Second * 180

// FetchFindAndQueryAndUpdateValidatorBalances Routine SIX
func FetchFindAndQueryAndUpdateValidatorBalances() {
	log.Info().Msg("FetchFindAndQueryAndUpdateValidatorBalances")

	for {
		timeBegin := time.Now()
		err := fetchAndUpdateValidatorBalances(context.Background(), NewValidatorBalancesBatchSize, NewValidatorBalancesTimeout)
		log.Err(err)
		log.Info().Interface("FetchFindAndQueryAndUpdateValidatorBalances took this many seconds to complete: ", time.Now().Sub(timeBegin))
	}
}

func fetchAndUpdateValidatorBalances(ctx context.Context, batchSize int, contextTimeout time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	err := Fetcher.FindAndQueryAndUpdateValidatorBalances(ctxTimeout, batchSize)
	log.Info().Err(err).Msg("FetchNewOrMissingValidators")
	return err
}
