package quicknode_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type HestiaQuickNodeWorker struct {
	temporal_base.Worker
}

var (
	HestiaQnWorker HestiaQuickNodeWorker
)

const (
	HestiaQuickNodeTaskQueue = "HestiaQuickNodeTaskQueue"
)

func InitHestiaQuickNodeWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Hestia: InitHestiaQuickNodeWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitHestiaQuickNodeWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := HestiaQuickNodeTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewHestiaQuickNodeActivities()
	wf := NewHestiaQuickNodeWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	HestiaQnWorker.Worker = w
	HestiaQnWorker.TemporalClient = tc
	return
}
