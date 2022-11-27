package artemis_ethereum_transcations

import (
	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var ArtemisTxBroadcastWorker temporal_base.Worker

const EthereumTxBroadcastTaskQueue = "EthereumTxBroadcastTaskQueue"

func InitTxBroadcastWorker(temporalAuthCfg temporal_auth.TemporalAuth) {
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitTxBroadcastWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)

	activityDef := ArtemisEthereumBroadcastTxActivities{}
	wf := NewArtemisBroadcastEthereumTxWorkflow()

	w.AddWorkflow(wf.ArtemisSendEthTxWorkflow)
	w.AddWorkflow(wf.ArtemisSendSignedTxWorkflow)

	w.AddActivities(activityDef.GetActivities())

	ArtemisTxBroadcastWorker = w
	ArtemisTxBroadcastWorker.TemporalClient = tc
	return
}
