package eth_validator_signature_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

type ArtemisEthereumValidatorSignatureRequestsWorker struct {
	temporal_base.Worker
}

var ArtemisEthereumValidatorSignatureRequestsMainnetWorker ArtemisEthereumValidatorSignatureRequestsWorker

// TODO, identify how to setup task queue better
// TODO, add ephemeral testnet, maybe make that part a param?

const EthereumTxBroadcastTaskQueue = "EthereumValidatorSignatureRequestsMainnetTaskQueue"

func InitArtemisEthereumValidatorSignatureRequestsMainnetWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitArtemisEthereumValidatorSignatureRequestsMainnetWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitArtemisEthereumValidatorSignatureRequestsMainnetWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorSignatureRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumValidatorSignatureRequestsMainnetWorker.Worker = w
	ArtemisEthereumValidatorSignatureRequestsMainnetWorker.TemporalClient = tc
	return
}
