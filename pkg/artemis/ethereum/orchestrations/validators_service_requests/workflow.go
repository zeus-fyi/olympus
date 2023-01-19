package eth_validators_service_requests

import (
	"time"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
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

func (t *ArtemisNewEthereumValidatorsServiceRequestWorkflow) ServiceNewValidatorsToCloudCtxNsWorkflow(ctx workflow.Context, params artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	// TODO, if this keeps failing, terminate workflow
	validateValidatorsRemoteServicesStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(validateValidatorsRemoteServicesStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.ValidateKeysToServiceURL, params).Get(validateValidatorsRemoteServicesStatusCtx, nil)
	if err != nil {
		log.Warn("Failed to validate key to service url", "ValidatorServiceRequest", params)
		log.Error("Failed to validate key to service url", "Error", err)
		return err
	}
	// If succeed, continue forever or basically forever, TODO alert if failure continues for > 1hr
	assignValidatorsStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(assignValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.AssignValidatorsToCloudCtxNs, params).Get(assignValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to assign validators to cloud ctx ns", "Error", err)
		return err
	}

	updateClusterValidatorsStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(updateClusterValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.RestartValidatorClient, params).Get(assignValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to assign validators to cloud ctx ns", "Error", err)
		return err
	}
	//
	return nil
}
