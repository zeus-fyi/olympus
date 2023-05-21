package artemis_mev_transcations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type ArtemisMevWorker struct {
	temporal_base.Worker
}

var ArtemisMevWorkerMainnet ArtemisMevWorker

const EthereumMainnetTaskQueue = "EthereumMainnetTaskQueue"

func InitMainnetEthereumTxBroadcastWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMainnetEthereumTxBroadcastWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitMainnetEthereumTxBroadcastWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumMainnetTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientMainnet)
	wf := NewArtemisMevWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisMevWorkerMainnet.Worker = w
	ArtemisMevWorkerMainnet.TemporalClient = tc
	return
}

var ArtemisMevWorkerGoerli ArtemisMevWorker

const EthereumGoerliTaskQueue = "EthereumGoerliTaskQueue"

func InitGoerliEthereumTxBroadcastWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitGoerliEthereumTxBroadcastWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitGoerliEthereumTxBroadcastWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumGoerliTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientGoerli)
	wf := NewArtemisMevWorkflow()
	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisMevWorkerGoerli.Worker = w
	ArtemisMevWorkerGoerli.TemporalClient = tc
	return
}

func InitMevWorkers(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitWeb3Clients")
	InitWeb3Clients(ctx)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitMainnetEthereumTxBroadcastWorker")
	InitMainnetEthereumTxBroadcastWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: IInitMevWorkers: nitGoerliEthereumTxBroadcastWorker")
	InitGoerliEthereumTxBroadcastWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers succeeded")
}
