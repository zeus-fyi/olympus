package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
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

const (
	requestApprovedState = "approved"
	requestRejectedState = "rejected"
)

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
			MaximumAttempts:    25,
		},
	}

	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(approvalTaskGroup.Ou.OrgID, approvalTaskGroup.WfID, "ZeusAiPlatformServiceWorkflows", "TriggerActionsWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update ai trigger hil action", "Error", err)
		return err
	}
	if approvalTaskGroup.RequestedState == requestRejectedState {
		for _, v := range approvalTaskGroup.Taps {
			recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			v.ApprovalState = requestRejectedState
			err = workflow.ExecuteActivity(recordTriggerCondCtx, z.UpdateTriggerActionApproval, approvalTaskGroup.Ou, v).Get(recordTriggerCondCtx, nil)
			if err != nil {
				logger.Error("failed to create or update trigger action", "Error", err)
				return err
			}
			finishedCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
			if err != nil {
				logger.Error("failed to update cache for qn services", "Error", err)
				return err
			}
		}
		return nil
	}
	for _, v := range approvalTaskGroup.Taps {
		switch v.TriggerAction {
		case apiApproval:
			var apiApprovalReqs []artemis_orchestrations.ApprovalApiReqResp
			// if conditions are met, create or update the trigger action
			recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			v.ApprovalState = approvalTaskGroup.RequestedState
			err = workflow.ExecuteActivity(recordTriggerCondCtx, z.SelectTriggerActionApiApprovalWithReqResponses, approvalTaskGroup.Ou, "pending", v.ApprovalID, v.WorkflowResultID).Get(recordTriggerCondCtx, &apiApprovalReqs)
			if err != nil {
				logger.Error("failed to create or update trigger action", "Error", err)
				return err
			}
			if approvalTaskGroup.RequestedState != requestApprovedState {
				continue
			}
			if len(apiApprovalReqs) <= 0 {
				continue
			}
			for i, ar := range apiApprovalReqs {
				if GetRetryPolicy(ar.RetrievalItem, 5*time.Minute) != nil {
					aoAiAct.RetryPolicy = GetRetryPolicy(ar.RetrievalItem, 5*time.Minute)
				}
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               CreateExecAiWfId(approvalTaskGroup.WfID + "-api-approval-" + v.ApprovalStrID + "-" + strconv.Itoa(i)),
					RetryPolicy:              aoAiAct.RetryPolicy,
					WorkflowExecutionTimeout: 10 * time.Minute,
				}
				ar.RetrievalItem.RetrievalPlatform = apiApproval
				cp := &MbChildSubProcessParams{
					WfID: approvalTaskGroup.WfID,
					Ou:   approvalTaskGroup.Ou,
					Oj:   oj,
					Wsr: artemis_orchestrations.WorkflowStageReference{
						ChildWfID: childAnalysisWorkflowOptions.WorkflowID,
					},
					Tc: TaskContext{
						TriggerActionsApproval:             ar.TriggerActionsApproval,
						EvalID:                             ar.TriggerActionsApproval.EvalID,
						Retrieval:                          ar.RetrievalItem,
						AIWorkflowTriggerResultApiResponse: ar.AIWorkflowTriggerResultApiReqResponse,
					},
				}
				childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
				err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, cp).Get(childAnalysisCtx, &cp)
				if err != nil {
					logger.Error("failed to execute child api retrieval workflow", "Error", err)
					return err
				}
			}
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}
