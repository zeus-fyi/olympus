package temporal_base

import (
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"go.temporal.io/sdk/client"
)

func (t *TemporalClient) ConnectTemporalClient() client.Client {
	c, err := client.Dial(t.Options)
	if err != nil {
		log.Fatal().Err(err).Msg("ConnectTemporalClient: dial failed")
		misc.DelayedPanic(err)
	}
	return c
}
