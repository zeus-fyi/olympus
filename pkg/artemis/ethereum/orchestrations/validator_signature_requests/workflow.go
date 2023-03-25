package eth_validator_signature_requests

import (
	"time"

	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisEthereumValidatorSignatureRequestWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumValidatorSignatureRequestActivities
}

const defaultTimeout = 6 * time.Second

func NewArtemisEthereumValidatorSignatureRequestWorkflow() ArtemisEthereumValidatorSignatureRequestWorkflow {
	deployWf := ArtemisEthereumValidatorSignatureRequestWorkflow{
		temporal_base.Workflow{},
		ArtemisEthereumValidatorSignatureRequestActivities{},
	}
	return deployWf
}

func (t *ArtemisEthereumValidatorSignatureRequestWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.ArtemisSendValidatorSignatureRequestsWorkflow, t.ValidatorsHeartbeatWorkflow}
}

func (t *ArtemisEthereumValidatorSignatureRequestWorkflow) ArtemisSendValidatorSignatureRequestsWorkflow(ctx workflow.Context, sigRequests aegis_inmemdbs.EthereumBLSKeySignatureRequests) (aegis_inmemdbs.EthereumBLSKeySignatureResponses, error) {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	// TODO group pubkeys by serverless function then send requests
	var sigResponses aegis_inmemdbs.EthereumBLSKeySignatureResponses
	sigRespCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(sigRespCtx, t.RequestValidatorSignatures, sigRequests).Get(sigRespCtx, &sigResponses)
	if err != nil {
		log.Error("Failed to get signatures", "Error", err)
		return sigResponses, err
	}
	return sigResponses, nil
}
