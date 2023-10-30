package iris_serverless

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type IrisServicesWorker struct {
	temporal_base.Worker
}

var (
	IrisPlatformServicesWorker IrisServicesWorker
)

const (
	IrisPlatformServicesTaskQueue = "IrisPlatformServicesTaskQueue"
)

func InitIrisPlatformServicesWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Iris: InitIrisPlatformServicesWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitIrisPlatformServicesWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := IrisPlatformServicesTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewIrisPlatformActivities()
	wf := NewIrisPlatformServiceWorkflows()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())

	IrisPlatformServicesWorker.Worker = w
	IrisPlatformServicesWorker.TemporalClient = tc
	return
}
