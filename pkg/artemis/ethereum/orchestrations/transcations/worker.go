package artemis_ethereum_transcations

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var ArtemisEthereumTxBroadcastWorker temporal_base.Worker

const EthereumTxBroadcastTaskQueue = "EthereumTxBroadcastTaskQueue"

func InitEthereumTxBroadcastWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitEthereumTxBroadcastWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitEthereumTxBroadcastWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumBroadcastTxActivities(ArtemisEthereumBroadcastTxClient)
	wf := NewArtemisBroadcastEthereumTxWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumTxBroadcastWorker = w
	ArtemisEthereumTxBroadcastWorker.TemporalClient = tc
	return
}

var ArtemisEthereumGoerliTxBroadcastWorker temporal_base.Worker

const EthereumGoerliTxBroadcastTaskQueue = "EthereumGoerliTxBroadcastTaskQueue"

func InitEthereumGoerliTxBroadcastWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitEthereumGoerliTxBroadcastWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitEthereumGoerliTxBroadcastWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumGoerliTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumBroadcastTxActivities(ArtemisEthereumGoerliBroadcastTxClient)
	wf := NewArtemisBroadcastEthereumTxWorkflow()
	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumGoerliTxBroadcastWorker = w
	ArtemisEthereumGoerliTxBroadcastWorker.TemporalClient = tc
	return
}

func InitEthereumBroadcasters(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitEthereumBroadcasters")
	InitWeb3Clients(ctx)
	InitEthereumTxBroadcastWorker(ctx, temporalAuthCfg)
	InitEthereumGoerliTxBroadcastWorker(ctx, temporalAuthCfg)
	log.Ctx(ctx).Info().Msg("Artemis: InitEthereumBroadcasters succeeded")

}
