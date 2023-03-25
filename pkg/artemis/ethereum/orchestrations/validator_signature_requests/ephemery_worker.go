package eth_validator_signature_requests

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
)

const (
	EthereumEphemeryTxBroadcastTaskQueue = "EthereumValidatorSignatureRequestsEphemeryTaskQueue"
)

var (
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorker          ArtemisEthereumValidatorSignatureRequestsWorker
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary ArtemisEthereumValidatorSignatureRequestsWorker
)

func InitArtemisEthereumValidatorSignatureRequestsEphemeryWorker(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: EthereumValidatorSignatureRequestsEphemeryTaskQueue")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("EthereumValidatorSignatureRequestsEphemeryTaskQueue: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := EthereumEphemeryTxBroadcastTaskQueue
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorSignatureRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.Worker = w
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorker.TemporalClient = tc
	return
}

func InitArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary(ctx context.Context, temporalAuthCfg temporal_auth.TemporalAuth) {
	log.Ctx(ctx).Info().Msg("Artemis: InitArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary")
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := fmt.Sprintf("%sSecondary", EthereumEphemeryTxBroadcastTaskQueue)
	w := temporal_base.NewWorker(taskQueueName)
	activityDef := NewArtemisEthereumValidatorSignatureRequestActivities()
	wf := NewArtemisEthereumValidatorSignatureRequestWorkflow()

	w.AddWorkflows(wf.GetWorkflows())
	w.AddActivities(activityDef.GetActivities())
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary.Worker = w
	ArtemisEthereumValidatorSignatureRequestsEphemeryWorkerSecondary.TemporalClient = tc
	return
}
