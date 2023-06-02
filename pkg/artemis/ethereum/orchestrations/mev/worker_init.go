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

const (
	EthereumMainnetTaskQueue                = "EthereumMainnetTaskQueue"
	EthereumMainnetMevHistoricalTxTaskQueue = "EthereumMainnetMevHistoricalTxTaskQueue"
)

func InitMainnetEthereumMevWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMainnetEthereumMevWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitMainnetEthereumMevWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumMainnetTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientMainnet)
	activityDef.Network = "mainnet"
	wf := NewArtemisMevWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisMevWorkerMainnet.Worker = w
	ArtemisMevWorkerMainnet.TemporalClient = tc
	return
}

func InitMainnetEthereumMevHistoricalWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMainnetEthereumMevHistoricalWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitMainnetEthereumMevHistoricalWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumMainnetMevHistoricalTxTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientMainnet)
	activityDef.Network = "mainnet"
	wf := NewArtemisMevWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisMevWorkerMainnet.Worker = w
	ArtemisMevWorkerMainnet.TemporalClient = tc
	return
}

var ArtemisMevWorkerGoerli ArtemisMevWorker

const EthereumGoerliTaskQueue = "EthereumGoerliTaskQueue"

func InitGoerliEthereumMevWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitGoerliEthereumMevWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitGoerliEthereumMevWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumGoerliTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientGoerli)
	activityDef.Network = "goerli"

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
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitMainnetEthereumMevWorker")
	InitMainnetEthereumMevWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: IInitMevWorkers: InitGoerliEthereumMevWorker")
	InitGoerliEthereumMevWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers succeeded")
}
