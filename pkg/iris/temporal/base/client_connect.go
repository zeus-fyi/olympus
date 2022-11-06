package temporal_base

import (
	"go.temporal.io/sdk/client"
)

func (t *TemporalClient) Connect() error {
	temporalClient, err := client.Dial(t.Options)
	if err != nil {
		return err
	}
	t.Client = temporalClient
	return err
}
