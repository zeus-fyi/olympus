package artemis_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func InitArtemisApiRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitArtemisApiRequestsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitArtemisApiRequestsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := ApiRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisApiRequestsActivities()
	wf := NewArtemisApiRequestsWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisProxyWorker.Worker = w
	ArtemisProxyWorker.TemporalClient = tc
	return
}
