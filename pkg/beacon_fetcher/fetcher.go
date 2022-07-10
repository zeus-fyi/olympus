package beacon_fetcher

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	beacon_models "github.com/zeus-fyi/olympus/databases/postgres/beacon-indexer/beacon-models"
)

// temp value 406,180 number of validataors
var numValidators int64 = 406180

func FetchBeaconState() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	sleepBetweenFetches := time.Second * 10

	for {
		ctx := context.Background()
		vbalIndexes, err := beacon_models.SelectValidatorIndexesInStrArrayForQueryURL(ctx, 1000)
		log.Error().Err(err).Msg("")
		log.Info().Interface("URL encoded Validator Indexes", vbalIndexes)

		time.Sleep(sleepBetweenFetches)
	}
}
