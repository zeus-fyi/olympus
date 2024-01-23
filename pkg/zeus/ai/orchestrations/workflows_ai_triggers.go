package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TriggerActionsWorkflowParams struct {
	Emr *artemis_orchestrations.EvalMetricsResults `json:"emr,omitempty"`
	Mb  *MbChildSubProcessParams                   `json:"mb,omitempty"`
}

func (z *ZeusAiPlatformServiceWorkflows) CreateTriggerActionsWorkflow(ctx workflow.Context, tar TriggerActionsWorkflowParams) error {
	if tar.Emr == nil || tar.Mb == nil {
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

	ou := tar.Mb.Ou
	triggerEvalsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	var triggerActions []artemis_orchestrations.TriggerAction
	// looks up if there are any trigger actions to execute by eval id
	err := workflow.ExecuteActivity(triggerEvalsLookupCtx, z.LookupEvalTriggerConditions, ou, tar.Emr.EvalContext.EvalID).Get(triggerEvalsLookupCtx, &triggerActions)
	if err != nil {
		logger.Error("failed to get eval info", "Error", err)
		return err
	}

	// if there are no trigger actions to execute, check if conditions are met for execution
	for _, triggerAction := range triggerActions {
		var ta *artemis_orchestrations.TriggerAction
		checkTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(checkTriggerCondCtx, z.CheckEvalTriggerCondition, &triggerAction, tar.Emr).Get(checkTriggerCondCtx, &ta)
		if err != nil {
			logger.Error("failed to check eval trigger condition", "Error", err)
			return err
		}

		if ta == nil {
			continue
		}
		var cr *ChatCompletionQueryResponse
		if tar.Mb.AnalysisEvalActionParams != nil && tar.Mb.AnalysisEvalActionParams.SearchResultGroup != nil {
			ta.TriggerPlatformReference.PlatformReferenceName = tar.Mb.AnalysisEvalActionParams.SearchResultGroup.PlatformName
			switch ta.TriggerAction {
			case socialMediaEngagementResponseFormat:
				smApiEvalFormatCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(smApiEvalFormatCtx, z.EvalFormatForApi, ou, ta).Get(smApiEvalFormatCtx, &cr)
				if err != nil {
					logger.Error("failed to check eval trigger condition", "Error", err)
					return err
				}
			}
		}
		if cr != nil {
			var triggerCompletionID int
			triggerCompletionResponseCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err = workflow.ExecuteActivity(triggerCompletionResponseCtx, z.RecordCompletionResponse, ou, cr).Get(triggerCompletionResponseCtx, &triggerCompletionID)
			if err != nil {
				logger.Error("failed to get completion response", "Error", err)
				return err
			}
			trr := artemis_orchestrations.AIWorkflowTriggerResultResponse{
				WorkflowResultID: tar.Mb.WorkflowResult.WorkflowResultID,
				TriggerID:        ta.TriggerID,
				ResponseID:       triggerCompletionID,
			}
			saveCompletionResponseCtx := workflow.WithActivityOptions(ctx, aoAiAct)
			err = workflow.ExecuteActivity(saveCompletionResponseCtx, z.SaveTriggerResponseOutput, trr).Get(saveCompletionResponseCtx, nil)
			if err != nil {
				logger.Error("failed to save trigger response", "Error", err)
				return err
			}
		}
		// if conditions are met, create or update the trigger action
		recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(recordTriggerCondCtx, z.CreateOrUpdateTriggerActionToExec, tar.Mb, ta).Get(recordTriggerCondCtx, nil)
		if err != nil {
			logger.Error("failed to create or update trigger action", "Error", err)
			return err
		}
	}
	return nil
}
