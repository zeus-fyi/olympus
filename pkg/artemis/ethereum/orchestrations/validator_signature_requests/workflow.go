package eth_validator_signature_requests

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisArtemisEthereumValidatorSignatureRequestWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumValidatorSignatureRequestActivities
}

const defaultTimeout = 6 * time.Second

func NewArtemisEthereumValidatorSignatureRequestWorkflow() ArtemisArtemisEthereumValidatorSignatureRequestWorkflow {
	deployWf := ArtemisArtemisEthereumValidatorSignatureRequestWorkflow{
		temporal_base.Workflow{},
		ArtemisEthereumValidatorSignatureRequestActivities{},
	}
	return deployWf
}

func (t *ArtemisArtemisEthereumValidatorSignatureRequestWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.ArtemisSendValidatorSignatureRequestWorkflow}
}

func (t *ArtemisArtemisEthereumValidatorSignatureRequestWorkflow) ArtemisSendValidatorSignatureRequestWorkflow(ctx workflow.Context, params interface{}) error {
	//log := workflow.GetLogger(ctx)
	//ao := workflow.ActivityOptions{
	//	StartToCloseTimeout: defaultTimeout,
	//}

	// TODO, send serverless request, wait for reply for only as much time as makes since. eg 12s slot, 6.3min epoch
	// don't block other validators

	// TODO, send request back to origin, then done

	return nil
}
