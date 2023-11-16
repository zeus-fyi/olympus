package eth_validators_service_requests

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_orchestration_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/orchestration_auth"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	zeus_client "github.com/zeus-fyi/zeus/zeus/z_client"
)

var Zeus zeus_client.ZeusClient

func InitZeusClientValidatorServiceGroup(ctx context.Context) {
	log.Info().Msg("Artemis: InitZeusClientValidatorServiceGroup")
	Zeus = zeus_client.NewDefaultZeusClient(artemis_orchestration_auth.Bearer)
	return
}

func InitArtemisEthereumEphemeryValidatorsRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Artemis: ArtemisEthereumEphemeryValidatorsRequestsWorker")
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
	ArtemisEthereumEphemeryValidatorsRequestsWorker = ArtemisEthereumValidatorsRequestsWorker{Worker: w}
	ArtemisEthereumEphemeryValidatorsRequestsWorker.TemporalClient = tc
	return
}

func InitArtemisEthereumMainnetValidatorsRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Artemis: ArtemisEthereumMainnetValidatorsRequestsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("ArtemisEthereumMainnetValidatorsRequestsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumMainnetValidatorsRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorServiceRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumMainnetValidatorsRequestsWorker = ArtemisEthereumValidatorsRequestsWorker{Worker: w}
	ArtemisEthereumMainnetValidatorsRequestsWorker.TemporalClient = tc
	return
}

func InitArtemisEthereumGoerliValidatorsRequestsWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Info().Msg("Artemis: InitArtemisEthereumGoerliValidatorsRequestsWorker")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitArtemisEthereumGoerliValidatorsRequestsWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumGoerliValidatorsRequestsTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorServiceRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumGoerliValidatorsRequestsWorker = ArtemisEthereumValidatorsRequestsWorker{Worker: w}
	ArtemisEthereumGoerliValidatorsRequestsWorker.TemporalClient = tc
	return
}
