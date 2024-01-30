package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/*
1. should update state to approved or should update state to rejected, then end workflow
2. should lookup trigger information
	- social media engagement (twitter, reddit, discord, telegram)
		execute API calls

*/

type ApprovalTaskGroup struct {
	WfID           string                                          `json:"wfID"`
	RequestedState string                                          `json:"requestedState"`
	Ou             org_users.OrgUser                               `json:"ou"`
	Taps           []artemis_orchestrations.TriggerActionsApproval `json:"taps"`
}

func (z *ZeusAiPlatformServiceWorkflows) TriggerActionsWorkflow(ctx workflow.Context, approvalTaskGroup ApprovalTaskGroup) error {
	if len(approvalTaskGroup.Taps) == 0 {
		return nil
	}
	logger := workflow.GetLogger(ctx)
	aoAiAct := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    5,
		},
	}

	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(approvalTaskGroup.Ou.OrgID, approvalTaskGroup.WfID, "ZeusAiPlatformServiceWorkflows", "TriggerActionsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai trigger hil action", "Error", err)
		return err
	}

	for _, v := range approvalTaskGroup.Taps {
		var ta artemis_orchestrations.TriggerAction
		switch ta.TriggerAction {
		case apiApproval:
			var apiApprovalReqs []artemis_orchestrations.ApprovalApiReqResp
			// if conditions are met, create or update the trigger action
			recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err = workflow.ExecuteActivity(recordTriggerCondCtx, z.SelectTriggerActionApiApprovalWithReqResponses, approvalTaskGroup.Ou, approvalTaskGroup.RequestedState, v.ApprovalID).Get(recordTriggerCondCtx, &apiApprovalReqs)
			if err != nil {
				logger.Error("failed to create or update trigger action", "Error", err)
				return err
			}
			if len(apiApprovalReqs) <= 0 {
				continue
			}
			ar := apiApprovalReqs[0]
			tte := TaskToExecute{
				WfID: approvalTaskGroup.WfID + "-api-approval-" + v.ApprovalStrID,
				Ou:   approvalTaskGroup.Ou,
				Ec:   artemis_orchestrations.EvalContext{},
				Tc: TaskContext{
					EvalID:                             ar.TriggerActionsApproval.EvalID,
					Retrieval:                          ar.RetrievalItem,
					AIWorkflowTriggerResultApiResponse: ar.AIWorkflowTriggerResultApiReqResponse,
				},
				Sg:          &hera_search.SearchResultGroup{},
				RetryPolicy: GetRetryPolicy(ar.RetrievalItem, 5*time.Minute),
			}

			childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
				WorkflowID: approvalTaskGroup.WfID,
				//WorkflowExecutionTimeout: tar.Mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			}
			childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
			err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, tte).Get(childAnalysisCtx, &tte.Sg)
			if err != nil {
				logger.Error("failed to execute child api retrieval workflow", "Error", err)
				return err
			}
		case socialMediaEngagementResponseFormat:
			//childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			//	//WorkflowID:               mb.Oj.OrchestrationName + "-eval-trigger-" + strconv.Itoa(mb.RunCycle) + suffix,
			//	//WorkflowExecutionTimeout: mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			//}
			//childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
			//err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunApprovedSocialMediaTriggerActionsWorkflow, tar).Get(childAnalysisCtx, nil)
			//if err != nil {
			//	logger.Error("failed to execute child run trigger actions workflow", "Error", err)
			//	return err
			//}
		}
	}
	return nil
}
