package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type ZeusAiPlatformServicesWorker struct {
	temporal_base.Worker
}

var (
	ZeusAiPlatformWorker ZeusAiPlatformServicesWorker
)

const (
	ZeusAiPlatformServicesTaskQueue = "ZeusAiPlatformServicesTaskQueue"
)

func InitZeusAiServicesWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Hestia: InitZeusAiServicesWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitZeusAiServicesWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := ZeusAiPlatformServicesTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewZeusAiPlatformActivities()
	wf := NewZeusPlatformServiceWorkflows()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ZeusAiPlatformWorker.Worker = w
	ZeusAiPlatformWorker.TemporalClient = tc
	return
}
