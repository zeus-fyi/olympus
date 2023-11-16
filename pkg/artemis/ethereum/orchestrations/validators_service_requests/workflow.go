package eth_validators_service_requests

import (
	"time"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/workflow"
)

type ArtemisNewEthereumValidatorsServiceRequestWorkflow struct {
	temporal_base.Workflow
	ArtemisEthereumValidatorsServiceRequestActivities
}

// TODO revise this timeout
const defaultTimeout = 300 * time.Second

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

func (t *ArtemisNewEthereumValidatorsServiceRequestWorkflow) ServiceNewValidatorsToCloudCtxNsWorkflow(ctx workflow.Context, params ValidatorServiceGroupWorkflowRequest) error {
	log := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
	}
	// TODO, if this continues failing, should terminate workflow and send email
	validateValidatorsRemoteServicesStatusCtx := workflow.WithActivityOptions(ctx, ao)
	var verifiedPubkeys hestia_req_types.ValidatorServiceOrgGroupSlice
	err := workflow.ExecuteActivity(validateValidatorsRemoteServicesStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.VerifyValidatorKeyOwnershipAndSigning, params).Get(validateValidatorsRemoteServicesStatusCtx, &verifiedPubkeys)
	if err != nil {
		log.Warn("Failed to validate key to service url", "ValidatorServiceRequest", params)
		log.Error("Failed to validate key to service url", "Error", err)
		workflow.Sleep(ctx, time.Second*30)
		return err
	}
	insertParams := artemis_validator_service_groups_models.OrgValidatorService{
		GroupName:         params.GroupName,
		ProtocolNetworkID: params.ProtocolNetworkID,
		ServiceURL:        "https://deprecated.com",
		OrgID:             params.OrgID,
		Enabled:           params.Enabled,
		MevEnabled:        params.MevEnabled,
	}
	if len(verifiedPubkeys) <= 0 {
		log.Info("No new validators to insert", "ValidatorServiceRequest", params)
		return nil
	}
	insertVerifiedValidatorsStatusCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(insertVerifiedValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.InsertVerifiedValidatorsWithFeeRecipient, insertParams, verifiedPubkeys).Get(insertVerifiedValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to insert new validators to service", "Error", err)
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
	err = workflow.ExecuteActivity(updateClusterValidatorsStatusCtx, t.ArtemisEthereumValidatorsServiceRequestActivities.RestartValidatorClient, params).Get(updateClusterValidatorsStatusCtx, nil)
	if err != nil {
		log.Error("Failed to restart validators client", "Error", err)
		return err
	}
	return nil
}
