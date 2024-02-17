package ai_platform_service_orchestrations

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
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
			MaximumAttempts:    25,
		},
	}
	tq := artemis_orchestrations.TriggersWorkflowQueryParams{
		Ou:                 tar.Mb.Ou,
		EvalID:             tar.Cpe.TaskToExecute.Tc.EvalID,
		TaskID:             tar.Cpe.TaskToExecute.Tc.TaskID,
		WorkflowTemplateID: tar.Mb.WfExecParams.WorkflowTemplate.WorkflowTemplateID,
	}
	if !tq.ValidateEvalTaskQp() {
		return nil
	}

	var triggerActions []artemis_orchestrations.TriggerAction
	triggerEvalsLookupCtx := workflow.WithActivityOptions(ctx, aoAiAct)
	err := workflow.ExecuteActivity(triggerEvalsLookupCtx, z.LookupEvalTriggerConditions, tq).Get(triggerEvalsLookupCtx, &triggerActions)
	if err != nil {
		logger.Error("failed to get eval trigger info", "Error", err)
		return err
	}

	// if there are no trigger actions to execute, check if conditions are met for execution for filter
	if len(triggerActions) == 0 && tar.Emr.EvaluatedJsonResponses != nil {
		// just filter passing to next stage then if no trigger action with specific pass/fail conditions
		jro := FilterPassingEvalPassingResponses(tar.Emr.EvaluatedJsonResponses)
		var payloadJsonSlice []artemis_orchestrations.JsonSchemaDefinition
		// update to pass sg
		sgIn := tar.Cpe.SearchResultGroup
		if sgIn == nil {
			sgIn = &hera_search.SearchResultGroup{}
		}
		wfr := tar.Mb.WorkflowResult
		if tar.Cpe.ParentOutputToEval != nil {
			wfr.WorkflowResultID = tar.Cpe.ParentOutputToEval.WorkflowResultID
			wfr.ResponseID = tar.Cpe.ParentOutputToEval.ResponseID
		}
		updateTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(updateTaskCtx, z.UpdateTaskOutput, &wfr, jro, sgIn).Get(updateTaskCtx, &payloadJsonSlice)
		if err != nil {
			logger.Error("failed to update task", "Error", err)
			return err
		}
	}
	for _, ta := range triggerActions {
		var jro JsonResponseGroupsByOutcomeMap
		filterJsonEvalCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(filterJsonEvalCtx, z.FilterEvalJsonResponses, &ta, tar.Emr).Get(filterJsonEvalCtx, &jro)
		if err != nil {
			logger.Error("failed to check eval trigger condition", "Error", err)
			return err
		}
		var payloadJsonSlice []artemis_orchestrations.JsonSchemaDefinition
		sgIn := tar.Cpe.SearchResultGroup
		if sgIn == nil {
			sgIn = &hera_search.SearchResultGroup{}
		}
		wfr := tar.Mb.WorkflowResult
		if tar.Cpe.ParentOutputToEval != nil {
			wfr.WorkflowResultID = tar.Cpe.ParentOutputToEval.WorkflowResultID
			wfr.ResponseID = tar.Cpe.ParentOutputToEval.ResponseID
		}
		updateTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(updateTaskCtx, z.UpdateTaskOutput, &wfr, jro, sgIn).Get(updateTaskCtx, &payloadJsonSlice)
		if err != nil {
			logger.Error("failed to update task", "Error", err)
			return err
		}
		switch ta.TriggerAction {
		case apiApproval:
			/*
					EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"` // eg. info, filter, etc
					EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"` // all-pass, any-fail, etc

					this contains the result of all-pass, any-fail, etc per each element json item vs all elements
				    so we want to only use the passing elements regardless of the trigger on, since that is already accounted for
					on a per element level
			*/
			var echoReqs []echo.Map
			payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(payloadJsonSlice)
			for _, m := range payloadMaps {
				echoMap := echo.Map{}
				for k, v := range m {
					echoMap[k] = v
				}
				echoReqs = append(echoReqs, echoMap)
			}
			for _, ret := range ta.TriggerRetrievals {
				trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
					TriggerID:   ta.TriggerID,
					RetrievalID: aws.ToInt(ret.RetrievalID),
					ReqPayloads: echoReqs,
				}
				tap := artemis_orchestrations.TriggerActionsApproval{
					TriggerAction:    apiApproval,
					EvalID:           tq.EvalID,
					TriggerID:        ta.TriggerID,
					WorkflowResultID: tar.Emr.EvalContext.AIWorkflowAnalysisResult.WorkflowResultID,
					ApprovalState:    pendingStatus,
					RequestSummary:   "Requesting approval for trigger action",
				}
				recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(recordTriggerCondCtx, z.CreateOrUpdateTriggerActionApprovalWithApiReq, tar.Mb.Ou, tap, trrr).Get(recordTriggerCondCtx, nil)
				if err != nil {
					logger.Error("failed to create or update trigger action approval for api", "Error", err)
					return err
				}
			}

		}
	}
	return nil
}

func FilterPassingEvalPassingResponses(jres []artemis_orchestrations.JsonSchemaDefinition) JsonResponseGroupsByOutcomeMap {
	jro := make(map[string]JsonResponseGroupsByOutcome)
	jro["filter"] = JsonResponseGroupsByOutcome{
		Passed: []artemis_orchestrations.JsonSchemaDefinition{},
		Failed: []artemis_orchestrations.JsonSchemaDefinition{},
	}
	for _, jr := range jres {
		if jr.ScoredEvalMetrics == nil {
			continue
		}
		count := 0
		for _, er := range jr.ScoredEvalMetrics {
			if er.EvalExpectedResultState == "ignore" {
				count += 1
				continue
			}
			if er.EvalState != "filter" {
				continue
			}
			if er.EvalExpectedResultState == "pass" && er.EvalMetricResult != nil && er.EvalMetricResult.EvalResultOutcomeBool != nil && *er.EvalMetricResult.EvalResultOutcomeBool {
				count += 1
			} else if er.EvalExpectedResultState == "fail" && er.EvalMetricResult != nil && er.EvalMetricResult.EvalResultOutcomeBool != nil && !*er.EvalMetricResult.EvalResultOutcomeBool {
				count += 1
			}
		}
		if count == len(jr.ScoredEvalMetrics) && len(jr.ScoredEvalMetrics) > 0 {
			tmp := jro["filter"]
			tmp.Passed = append(tmp.Passed, jr)
			jro["filter"] = tmp
		} else {
			tmp := jro["filter"]
			tmp.Failed = append(tmp.Failed, jr)
			jro["filter"] = tmp
		}
	}
	return jro
}
