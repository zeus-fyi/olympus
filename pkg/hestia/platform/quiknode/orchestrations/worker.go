package quicknode_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type HestiaQuicknodeWorker struct {
	temporal_base.Worker
}

var (
	HestiaQnWorker HestiaQuicknodeWorker
)

const (
	HestiaQuicknodeTaskQueue = "HestiaQuicknodeTaskQueue"
)

func InitHestiaQuicknodeWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Hestia: InitHestiaQuicknodeWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitHestiaQuicknodeWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := HestiaQuicknodeTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewHestiaQuicknodeActivities()
	wf := NewHestiaQuickNodeWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	HestiaQnWorker.Worker = w
	HestiaQnWorker.TemporalClient = tc
	return
}
