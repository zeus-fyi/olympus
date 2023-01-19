package eth_validators_service_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

var Zeus zeus_client.ZeusClient

func InitArtemisEthereumEphemeryValidatorsRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: ArtemisEthereumEphemeryValidatorsRequestsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("ArtemisEthereumEphemeryValidatorsRequestsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumEphemeryValidatorsRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorServiceRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumEphemeryValidatorsRequestsWorker.Worker = w
	ArtemisEthereumEphemeryValidatorsRequestsWorker.TemporalClient = tc
	return
}
