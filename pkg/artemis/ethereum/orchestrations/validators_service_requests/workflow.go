package eth_validators_service_requests

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisEthereumValidatorsServiceCreateRequestWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumValidatorsServiceRequestActivities
}

// TODO revise this timeout
const defaultTimeout = 6 * time.Second

func NewArtemisEthereumValidatorSignatureRequestWorkflow() ArtemisEthereumValidatorsServiceCreateRequestWorkflow {
	deployWf := ArtemisEthereumValidatorsServiceCreateRequestWorkflow{
		temporal_base.Workflow{},
		ArtemisEthereumValidatorsServiceRequestActivities{},
	}
	return deployWf
}

func (t *ArtemisEthereumValidatorsServiceCreateRequestWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.CreateAndDistributeEphemeryValidatorsToCloudCtxNs}
}

const ephemeryCloudCtxNs = 1671248907408699000

func (t *ArtemisEthereumValidatorsServiceCreateRequestWorkflow) CreateAndDistributeEphemeryValidatorsToCloudCtxNs(ctx workflow.Context, params interface{}) error {

	// TODO, write cloud location to database, if non-existent

	return nil
}
