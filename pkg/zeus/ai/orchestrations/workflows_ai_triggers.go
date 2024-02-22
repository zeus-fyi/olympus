package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type CreateTriggerActionsWorkflowInputs struct {
	Emr                                *artemis_orchestrations.EvalMetricsResults `json:"emr,omitempty"`
	RunAiWorkflowAutoEvalProcessInputs *RunAiWorkflowAutoEvalProcessInputs        `json:"runAiWorkflowAutoEvalProcessInputs,omitempty"`
}

const (
	apiApproval  = "api"
	apiRetrieval = "api-retrieval"
)

func (z *ZeusAiPlatformServiceWorkflows) CreateTriggerActionsWorkflow(ctx workflow.Context, cp *MbChildSubProcessParams) error {
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
		Ou:                 cp.Ou,
		EvalID:             cp.Tc.EvalID,
		TaskID:             cp.Tc.TaskID,
		WorkflowTemplateID: cp.WfExecParams.WorkflowTemplate.WorkflowTemplateID,
	}
	if !tq.ValidateEvalTaskQp() {
		log.Warn().Interface("tq", tq).Msg("CreateTriggerActionsWorkflow: invalid trigger query params")
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
	if len(triggerActions) == 0 && cp.Tc.JsonResponseResults != nil {
		// just filter passing to next stage then if no trigger action with specific pass/fail conditions
		updateTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(updateTaskCtx, z.UpdateTaskOutput, cp).Get(updateTaskCtx, nil)
		if err != nil {
			logger.Error("failed to update task", "Error", err)
			return err
		}
		return nil
	}
	log.Info().Interface("CreateTriggerActionsWorkflow: len(triggerActions)", len(triggerActions)).Msg("triggerActions")
	for _, ta := range triggerActions {
		var jro JsonResponseGroupsByOutcomeMap
		filterJsonEvalCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(filterJsonEvalCtx, z.FilterEvalJsonResponses, cp, &ta).Get(filterJsonEvalCtx, &jro)
		if err != nil {
			logger.Error("failed to check eval trigger condition", "Error", err)
			return err
		}
		var payloadJsonSlice []artemis_orchestrations.JsonSchemaDefinition
		updateTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
		err = workflow.ExecuteActivity(updateTaskCtx, z.UpdateTaskOutput, cp).Get(updateTaskCtx, &payloadJsonSlice)
		if err != nil {
			logger.Error("failed to update task", "Error", err)
			return err
		}
		if len(payloadJsonSlice) == 0 {
			logger.Warn("payload json slice is empty, skipping trigger action", "TriggerAction", ta)
			continue
		}
		switch ta.TriggerAction {
		case apiRetrieval:
			var echoReqs []echo.Map
			payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(payloadJsonSlice)
			for _, m := range payloadMaps {
				echoMap := echo.Map{}
				for k, v := range m {
					echoMap[k] = v
				}
				echoReqs = append(echoReqs, echoMap)
			}
			for ri, ret := range ta.TriggerRetrievals {
				var rets []artemis_orchestrations.RetrievalItem
				chunkedTaskCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(chunkedTaskCtx, z.SelectRetrievalTask, cp.Ou, ret.RetrievalID).Get(chunkedTaskCtx, &rets)
				if err != nil {
					logger.Error("failed to run ret task", "Error", err)
					return err
				}
				if len(rets) <= 0 {
					continue
				}

				/*
					1. get query params if/any from json payload in ext api call
					2. call ret wf to get data
					3. update task with results
				*/
				childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
					WorkflowID:               CreateExecAiWfId(cp.Oj.OrchestrationName + "-api-ret-" + strconv.Itoa(ri) + "-" + strconv.Itoa(cp.WfExecParams.WorkflowExecTimekeepingParams.CurrentCycleCount)),
					WorkflowExecutionTimeout: cp.WfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
					RetryPolicy:              aoAiAct.RetryPolicy,
				}
				cp.Tc.Retrieval = rets[0]
				cp.Wsr.ChildWfID = childAnalysisWorkflowOptions.WorkflowID
				switch *ret.WebFilters.PayloadPreProcessing {
				case "iterate":
					for _, ple := range echoReqs {
						cp.Tc.WebPayload = ple
						childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
						err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, cp).Get(childAnalysisCtx, &cp)
						if err != nil {
							logger.Error("failed to execute child retrieval workflow", "Error", err)
							return err
						}
					}
				case "bulk":
					cp.Tc.WebPayload = echoReqs
					childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
					err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RetrievalsWorkflow, cp).Get(childAnalysisCtx, &cp)
					if err != nil {
						logger.Error("failed to execute child retrieval workflow", "Error", err)
						return err
					}
				}

				updateTaskCtx = workflow.WithActivityOptions(ctx, aoAiAct)
				err = workflow.ExecuteActivity(updateTaskCtx, z.UpdateTaskOutput, cp).Get(updateTaskCtx, nil)
				if err != nil {
					logger.Error("failed to update task", "Error", err)
					return err
				}
			}
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
				var apiTrgs []*artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse
				if ret.WebFilters != nil && ret.WebFilters.PayloadPreProcessing != nil && len(payloadMaps) > 0 {
					switch *ret.WebFilters.PayloadPreProcessing {
					case "iterate":
						for _, ple := range echoReqs {
							trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
								TriggerID:   ta.TriggerID,
								RetrievalID: aws.ToInt(ret.RetrievalID),
								ReqPayloads: []echo.Map{ple},
							}
							apiTrgs = append(apiTrgs, trrr)
						}
					case "bulk":
						trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
							TriggerID:   ta.TriggerID,
							RetrievalID: aws.ToInt(ret.RetrievalID),
							ReqPayloads: echoReqs,
						}
						apiTrgs = append(apiTrgs, trrr)
					}
				} else {
					trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
						TriggerID:   ta.TriggerID,
						RetrievalID: aws.ToInt(ret.RetrievalID),
						ReqPayloads: echoReqs,
					}
					apiTrgs = append(apiTrgs, trrr)
				}

				for _, trrr := range apiTrgs {
					tap := artemis_orchestrations.TriggerActionsApproval{
						TriggerAction:    apiApproval,
						EvalID:           tq.EvalID,
						TriggerID:        ta.TriggerID,
						WorkflowResultID: cp.Tc.WorkflowResultID,
						ApprovalState:    pendingStatus,
						RequestSummary:   "Requesting approval for trigger action",
					}
					recordTriggerCondCtx := workflow.WithActivityOptions(ctx, aoAiAct)
					err = workflow.ExecuteActivity(recordTriggerCondCtx, z.CreateOrUpdateTriggerActionApprovalWithApiReq, cp.Ou, tap, trrr).Get(recordTriggerCondCtx, nil)
					if err != nil {
						logger.Error("failed to create or update trigger action approval for api", "Error", err)
						return err
					}
				}
			}
		}
	}
	return nil
}
