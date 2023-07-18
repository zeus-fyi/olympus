package artemis_mev_transcations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type ArtemisMevWorker struct {
	temporal_base.Worker
}

var (
	ArtemisMevWorkerMainnet              ArtemisMevWorker
	ArtemisActiveMevWorkerMainnet        ArtemisMevWorker
	ArtemisMevWorkerMainnetHistoricalTxs ArtemisMevWorker
)

const (
	EthereumMainnetTaskQueue                = "EthereumMainnetTaskQueue"
	EthereumMainnetMevHistoricalTxTaskQueue = "EthereumMainnetMevHistoricalTxTaskQueue"
	ActiveMainnetMEVTaskQueue               = "ActiveMainnetMEVTaskQueue"
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

func InitMainnetEthereumActiveMEVWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMainnetEthereumActiveMEVWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitMainnetEthereumActiveMEVWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := ActiveMainnetMEVTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientMainnet)
	activityDef.Network = hestia_req_types.Mainnet
	wf := NewArtemisMevWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisActiveMevWorkerMainnet.Worker = w
	ArtemisActiveMevWorkerMainnet.TemporalClient = tc
	return
}

func InitMainnetEthereumMevHistoricalTxsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitMainnetEthereumMevHistoricalTxsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitMainnetEthereumMevHistoricalTxsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumMainnetMevHistoricalTxTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisMevActivities(ArtemisMevClientMainnet)
	activityDef.Network = "mainnet"
	wf := NewArtemisMevWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisMevWorkerMainnetHistoricalTxs.Worker = w
	ArtemisMevWorkerMainnetHistoricalTxs.TemporalClient = tc
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
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitMainnetEthereumMevHistoricalTxsWorker")
	InitMainnetEthereumMevHistoricalTxsWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitMainnetEthereumActiveMEVWorker")
	InitMainnetEthereumActiveMEVWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers: InitGoerliEthereumMevWorker")
	InitGoerliEthereumMevWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitMevWorkers succeeded")
}
