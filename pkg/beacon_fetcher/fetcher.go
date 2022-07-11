package beacon_fetcher

import (
	"time"

	"github.com/rs/zerolog"
)

func FetchBeaconState() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	sleepBetweenFetches := time.Second * 10

	for {
		time.Sleep(sleepBetweenFetches)
	}
}
