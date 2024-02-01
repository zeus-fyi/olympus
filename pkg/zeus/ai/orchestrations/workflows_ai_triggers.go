package ai_platform_service_orchestrations

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"
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
		for _, tap := range taps {
			if tap.ApprovalID == 0 {
				continue
			}
			switch ta.TriggerAction {
			case apiApproval:
				var echoReqs []echo.Map
				payloadMaps := artemis_orchestrations.CreateMapInterfaceFromAssignedSchemaFields(tar.Cpe.TaskToExecute.Ec.JsonResponseResults)
				for _, m := range payloadMaps {
					echoMap := echo.Map{}
					for k, v := range m {
						echoMap[k] = v
					}
					echoReqs = append(echoReqs, echoMap)
				}
				for _, ret := range ta.TriggerRetrievals {
					trrr := &artemis_orchestrations.AIWorkflowTriggerResultApiReqResponse{
						ApprovalID:  tap.ApprovalID,
						TriggerID:   ta.TriggerID,
						RetrievalID: aws.ToInt(ret.RetrievalID),
						ReqPayloads: echoReqs,
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
	}
	return nil
}
