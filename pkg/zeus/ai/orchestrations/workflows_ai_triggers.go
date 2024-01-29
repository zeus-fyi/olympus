package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TriggerActionsWorkflowParams struct {
	Emr *artemis_orchestrations.EvalMetricsResults `json:"emr,omitempty"`
	Mb  *MbChildSubProcessParams                   `json:"mb,omitempty"`
	Cpe *EvalActionParams                          `json:"cpe,omitempty"`
}

const (
	apiApproval = "api"
)

func (z *ZeusAiPlatformServiceWorkflows) CreateTriggerActionsWorkflow(ctx workflow.Context, tar TriggerActionsWorkflowParams) error {
	if tar.Emr == nil || tar.Mb == nil || tar.Cpe == nil {
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

	triggerEvalsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	var triggerActions []artemis_orchestrations.TriggerAction
	tq := artemis_orchestrations.TriggersWorkflowQueryParams{
		Ou:                 tar.Mb.Ou,
		EvalID:             tar.Emr.EvalContext.EvalID,
		TaskID:             tar.Cpe.TaskToExecute.Tc.TaskID,
		WorkflowTemplateID: tar.Mb.WfExecParams.WorkflowTemplate.WorkflowTemplateID,
	}
	if !tq.ValidateEvalTaskQp() {
		return nil
	}
	err := workflow.ExecuteActivity(triggerEvalsLookupCtx, z.LookupEvalTriggerConditions, tq).Get(triggerEvalsLookupCtx, &triggerActions)
	if err != nil {
		logger.Error("failed to get eval trigger info", "Error", err)
		return err
	}

	// if there are no trigger actions to execute, check if conditions are met for execution
	for _, ta := range triggerActions {
		var taps []artemis_orchestrations.TriggerActionsApproval
		checkTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(checkTriggerCondCtx, z.CheckEvalTriggerCondition, &ta, tar.Emr).Get(checkTriggerCondCtx, &taps)
		if err != nil {
			logger.Error("failed to check eval trigger condition", "Error", err)
			return err
		}
		if taps == nil {
			continue
		}

		if tar.Mb.AnalysisEvalActionParams == nil && ta.TriggerAction == apiApproval {
			return nil
		}

		var cr *ChatCompletionQueryResponse
		if tar.Mb.AnalysisEvalActionParams != nil && tar.Mb.AnalysisEvalActionParams.SearchResultGroup != nil {
			switch ta.TriggerAction {
			case apiApproval:
				for _, ret := range ta.TriggerRetrievals {
					if ret.RetrievalID == nil || aws.ToInt(ret.RetrievalID) == 0 {
						continue
					}
					retWfID := tar.Mb.Oj.OrchestrationName + "-trigger-eval-api-ret-" + strconv.Itoa(tar.Mb.RunCycle) + "-chunk-" +
						strconv.Itoa(tar.Mb.WorkflowResult.ChunkOffset) + "-iter-" + strconv.Itoa(tar.Mb.WorkflowResult.IterationCount)
					for _, v := range taps {
						tte := TaskToExecute{
							WfID: retWfID,
							Ou:   tar.Mb.Ou,
							Sg:   tar.Mb.AnalysisEvalActionParams.SearchResultGroup,
							Wft:  tar.Mb.AnalysisEvalActionParams.WorkflowTemplateData,
							Tc: TaskContext{
								AIWorkflowTriggerResultApiResponse: artemis_orchestrations.AIWorkflowTriggerResultApiResponse{
									ApprovalID:  v.ApprovalID,
									TriggerID:   ta.TriggerID,
									RetrievalID: aws.ToInt(ret.RetrievalID),
								},
							},
						}
						tte.Wft.RetrievalPlatform = apiApproval
						tte.Wft.RetrievalID = ret.RetrievalID
						childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
							WorkflowID:               retWfID,
							WorkflowExecutionTimeout: tar.Mb.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
						}
						childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, tte).Get(childAnalysisCtx, &tar.Mb.AnalysisEvalActionParams.SearchResultGroup)
						if err != nil {
							logger.Error("failed to execute child api retrieval workflow", "Error", err)
							return err
						}
					}
				}
			case socialMediaEngagementResponseFormat:
				smApiEvalFormatCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(smApiEvalFormatCtx, z.EvalFormatForApi, tar.Mb.Ou, ta).Get(smApiEvalFormatCtx, &cr)
				if err != nil {
					logger.Error("failed to check eval trigger condition", "Error", err)
					return err
				}
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
