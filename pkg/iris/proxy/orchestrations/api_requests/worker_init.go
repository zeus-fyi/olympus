package iris_api_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func InitIrisApiRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Iris: InitIrisApiRequestsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitIrisRequestsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := ApiRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewIrisApiRequestsActivities()
	wf := NewIrisApiRequestsWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	IrisProxyWorker.Worker = w
	IrisProxyWorker.TemporalClient = tc
	return
}

func InitIrisCacheWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Iris: InitIrisCacheWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitIrisCacheWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := CacheUpdateRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewIrisApiRequestsActivities()
	wf := NewIrisApiRequestsWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	IrisCacheWorker.Worker = w
	IrisCacheWorker.TemporalClient = tc
	return
}
