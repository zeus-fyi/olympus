package artemis_ethereum_transcations

import (
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
)

var ArtemisTxBroadcastWorker temporal_base.Worker

const EthereumTxBroadcastTaskQueue = "EthereumTxBroadcastTaskQueue"

func InitTxBroadcastWorker(temporalAuthCfg temporal_auth.TemporalAuth) {
	//tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	//if err != nil {
	//	log.Err(err).Msg("InitTopologyWorker: NewTemporalClient failed")
	//	misc.DelayedPanic(err)
	//}
	taskQueueName := EthereumTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	// status
	ArtemisTxBroadcastWorker = w
	return
}
