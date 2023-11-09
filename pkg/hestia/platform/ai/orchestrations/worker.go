package ai_platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type HestiaAiPlatformServicesWorker struct {
	temporal_base.Worker
}

var (
	HestiaAiPlatformWorker HestiaAiPlatformServicesWorker
)

const (
	HestiaAiPlatformServicesTaskQueue = "HestiaAiPlatformServicesTaskQueue"
)

func InitHestiaIrisPlatformServicesWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Hestia: InitHestiaIrisPlatformServicesWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitHestiaIrisPlatformServicesWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := HestiaAiPlatformServicesTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewHestiaAiPlatformActivities()
	wf := NewHestiaPlatformServiceWorkflows()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	HestiaAiPlatformWorker.Worker = w
	HestiaAiPlatformWorker.TemporalClient = tc
	return
}
