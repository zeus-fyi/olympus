package kronos_helix

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

func InitKronosHelixWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Kronos: InitKronosHelixWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitKronosHelixWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := KronosHelixTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewKronosActivities()
	wf := NewKronosWorkflow()
	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	KronosServiceWorker.Worker = w
	KronosServiceWorker.TemporalClient = tc
	return
}
