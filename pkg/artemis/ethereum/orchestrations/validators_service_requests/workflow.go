package eth_validators_service_requests

import (
	"time"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisNewEthereumValidatorsServiceRequestWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumValidatorsServiceRequestActivities
}

// TODO revise this timeout
const defaultTimeout = 6 * time.Second

func NewArtemisEthereumValidatorServiceRequestWorkflow() ArtemisNewEthereumValidatorsServiceRequestWorkflow {
	deployWf := ArtemisNewEthereumValidatorsServiceRequestWorkflow{
		temporal_base.Workflow{},
		ArtemisEthereumValidatorsServiceRequestActivities{},
	}
	return deployWf
}

func (t *ArtemisNewEthereumValidatorsServiceRequestWorkflow) GetWorkflows() []interface{} {
	return []interface{}{t.ServiceNewValidatorsToCloudCtxNsWorkflow}
}

const (
	ephemeryCloudCtxNs = 1671248907408699000
	mainnsetCloudCtxNs = 1
)

type ArtemisCreateAndDistributeValidatorsToCloudCtxNsPayload struct {
	ArtemisEthereumValidatorsServiceRequestPayload
}

func (t *ArtemisNewEthereumValidatorsServiceRequestWorkflow) ServiceNewValidatorsToCloudCtxNsWorkflow(ctx workflow.Context, params ArtemisCreateAndDistributeValidatorsToCloudCtxNsPayload) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	assignValidatorsStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(assignValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.AssignValidatorsToCloudCtxNs, params.ArtemisEthereumValidatorsServiceRequestPayload).Get(assignValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to assign validators to cloud ctx ns", "Error", err)
		return err
	}

	updateClusterValidatorsStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(updateClusterValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.SendValidatorsToCloudCtxNs, params.ArtemisEthereumValidatorsServiceRequestPayload).Get(assignValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to assign validators to cloud ctx ns", "Error", err)
		return err
	}
	//
	return nil
}
