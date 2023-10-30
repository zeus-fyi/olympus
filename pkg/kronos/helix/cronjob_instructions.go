package kronos_helix

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

const (
	Cronjob = "cronjob"
)

type CronJobInstructions struct {
	Endpoint     string        `json:"endpoint"`
	PollInterval time.Duration `json:"pollInterval"` // can use 0 for no repeats?
}

func (k *KronosActivities) StartCronJobWorkflow(ctx context.Context, mi Instructions) error {
	rc := resty_base.GetBaseRestyClient(mi.CronJob.Endpoint, artemis_orchestration_auth.Bearer)
	resp, err := rc.R().Get(mi.CronJob.Endpoint)
	if err != nil {
		log.Err(err).Msg("HestiaPlatformActivities: IrisPlatformRefreshOrgGroupTableCacheRequest")
		return err
	}
	if resp.StatusCode() >= 400 {
		err = fmt.Errorf("CronJob endpoint %s is not healthy", mi.CronJob.Endpoint)
		log.Err(err).Msg("KronosActivities: StartCronJobWorkflow")
		return err
	}
	return nil
}
