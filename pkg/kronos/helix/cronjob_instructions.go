package kronos_helix

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
)

type CronJobInstructions struct {
	Endpoint     string        `json:"endpoint"`
	PollInterval time.Duration `json:"pollInterval"` // can use 0 for no repeats?
}

func (k *KronosActivities) StartCronJobWorkflow(ctx context.Context, mi Instructions) error {
	rc := resty.New()
	rc.SetAuthToken(artemis_orchestration_auth.Bearer)
	resp, err := rc.R().Get(mi.CronJob.Endpoint)
	if err != nil {
		log.Err(err).Msg("KronosActivities: StartCronJobWorkflow")
		return err
	}
	if resp.StatusCode() >= 400 {
		err = fmt.Errorf("CronJob endpoint %s is not healthy", mi.CronJob.Endpoint)
		log.Err(err).Msg("KronosActivities: StartCronJobWorkflow")
		return err
	}
	return nil
}
