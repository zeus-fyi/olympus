package temporal_base

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (t *TemporalClient) ConnectTemporalClient() error {
	dial, err := client.Dial(t.Options)
	if err != nil {
		log.Err(err).Msg("ConnectTemporalClient: dial failed")
		return err
	}
	t.Client = dial
	return nil
}
