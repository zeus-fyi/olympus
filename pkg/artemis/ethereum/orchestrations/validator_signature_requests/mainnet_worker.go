package eth_validator_signature_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

var ArtemisEthereumValidatorSignatureRequestsMainnetWorker ArtemisEthereumValidatorSignatureRequestsWorker

const EthereumTxBroadcastMainnetTaskQueue = "EthereumValidatorSignatureRequestsMainnetTaskQueue"

func InitArtemisEthereumValidatorSignatureRequestsMainnetWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitArtemisEthereumValidatorSignatureRequestsMainnetWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitArtemisEthereumValidatorSignatureRequestsMainnetWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumTxBroadcastMainnetTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorSignatureRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumValidatorSignatureRequestsMainnetWorker.Worker = w
	ArtemisEthereumValidatorSignatureRequestsMainnetWorker.TemporalClient = tc
	return
}
