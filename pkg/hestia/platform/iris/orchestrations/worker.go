package platform_service_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type HestiaPlatformServicesWorker struct {
	temporal_base.Worker
}

var (
	HestiaPlatformServiceWorker HestiaPlatformServicesWorker
)

const (
	HestiaIrisPlatformServicesTaskQueue = "HestiaIrisPlatformServicesTaskQueue"
)

func InitHestiaIrisPlatformServicesWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Hestia: InitHestiaIrisPlatformServicesWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitHestiaIrisPlatformServicesWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := HestiaIrisPlatformServicesTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewHestiaPlatformActivities()
	wf := NewHestiaPlatformServiceWorkflows()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	HestiaPlatformServiceWorker.Worker = w
	HestiaPlatformServiceWorker.TemporalClient = tc
	return
}
