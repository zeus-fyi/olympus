package ai_platform_service_orchestrations

import (
	"strconv"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowProcess(ctx workflow.Context, wfID string, ou org_users.OrgUser, wfExecParams artemis_orchestrations.WorkflowExecParams) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 15,
			MaximumAttempts:    25,
		},
	}

	ojCtx := workflow.WithActivityOptions(ctx, ao)
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, wfExecParams.WorkflowTemplate.WorkflowGroup, wfExecParams.WorkflowTemplate.WorkflowName)
	err := workflow.ExecuteActivity(ojCtx, z.UpsertAiOrchestration, ou, wfID, wfExecParams).Get(ojCtx, &oj.OrchestrationID)
	if err != nil {
		logger.Error("failed to UpsertAiOrchestration", "Error", err)
		return err
	}
	if oj.OrchestrationID == 0 {
		logger.Error("failed to UpsertAiOrchestration", "Error", err)
		return err
	}
	timer := UpdatableTimer{}
	err = workflow.SetQueryHandler(ctx, QueryType, func() (time.Time, error) {
		return timer.GetWakeUpTime(), nil
	})
	if err != nil {
		logger.Error("failed to set query handler", "Error", err)
		return err
	}
	if !wfExecParams.WorkflowExecTimekeepingParams.IsCycleStepped {
		err = timer.SleepUntil(ctx, wfExecParams.WorkflowExecTimekeepingParams.RunWindow.Start, workflow.GetSignalChannel(ctx, SignalType))
		if err != nil {
			logger.Error("failed to sleep until", "Error", err)
			return err
		}
	}
	startCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(startCtx, "UpdateAndMarkOrchestrationActive", oj).Get(startCtx, nil)
	if err != nil {
		logger.Error("failed to UpdateAndMarkOrchestrationActive", "Error", err)
		return err
	}
	for i := 1; i < wfExecParams.WorkflowExecTimekeepingParams.RunCycles+1; i++ {
		startTime := wfExecParams.WorkflowExecTimekeepingParams.RunWindow.Start.Add(time.Duration(i) * wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize)
		if time.Now().Before(startTime) && !wfExecParams.WorkflowExecTimekeepingParams.IsCycleStepped {
			err = workflow.Sleep(ctx, startTime.Sub(time.Now()))
			if err != nil {
				logger.Error("failed to sleep", "Error", err)
				return err
			}
		}
		wfChildID := CreateExecAiWfId(oj.OrchestrationName + "-analysis-" + strconv.Itoa(i))
		childParams := &MbChildSubProcessParams{WfID: wfChildID, Ou: ou, WfExecParams: wfExecParams, Oj: oj,
			Wsr: artemis_orchestrations.WorkflowStageReference{
				WorkflowRunID: oj.OrchestrationID,
				ChildWfID:     CreateExecAiWfId(oj.OrchestrationName + "-analysis-" + strconv.Itoa(i)),
				RunCycle:      i,
				ChunkOffset:   0,
			}}
		childAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{WorkflowID: childParams.WfID, WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize, RetryPolicy: ao.RetryPolicy}
		childAnalysisCtx := workflow.WithChildOptions(ctx, childAnalysisWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childAnalysisCtx, z.RunAiChildAnalysisProcessWorkflow, childParams).Get(childAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child analysis workflow", "Error", err)
			return err
		}
		logger.Info("RunAiWorkflowProcess: all child analysis workflow executed", "RunCycle", i)
		// Execute child workflow for aggregation
		wfAggChildID := CreateExecAiWfId(oj.OrchestrationName + "-agg-analysis-" + strconv.Itoa(i))
		childAggAnalysisWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID:               wfAggChildID,
			WorkflowExecutionTimeout: wfExecParams.WorkflowExecTimekeepingParams.TimeStepSize,
			ParentClosePolicy:        enums.PARENT_CLOSE_POLICY_ABANDON,
			RetryPolicy:              ao.RetryPolicy,
		}
		aggChildParams := &MbChildSubProcessParams{WfID: wfAggChildID, Ou: ou, WfExecParams: wfExecParams, Oj: oj,
			Wsr: artemis_orchestrations.WorkflowStageReference{
				WorkflowRunID: oj.OrchestrationID,
				ChildWfID:     childAggAnalysisWorkflowOptions.WorkflowID,
				RunCycle:      i,
				ChunkOffset:   0,
			}}
		logger.Info("RunAiWorkflowProcess: child aggregation workflow starting", "RunCycle", i)
		childAggAnalysisCtx := workflow.WithChildOptions(ctx, childAggAnalysisWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childAggAnalysisCtx, z.RunAiChildAggAnalysisProcessWorkflow, aggChildParams).Get(childAggAnalysisCtx, nil)
		if err != nil {
			logger.Error("failed to execute child aggregation workflow", "Error", err)
			return err
		}
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to mark inactive", "Error", err)
		return err
	}
	return nil
}

const (
	QueryType  = "GetWakeUpTime"
	SignalType = "UpdateWakeUpTime"
)

type UpdatableTimer struct {
	wakeUpTime time.Time
}

// SleepUntil sleeps until the provided wake-up time.
// The wake-up time can be updated at any time via a signal.
func (u *UpdatableTimer) SleepUntil(ctx workflow.Context, wakeUpTime time.Time, updateWakeUpTimeCh workflow.ReceiveChannel) (err error) {
	logger := workflow.GetLogger(ctx)
	u.wakeUpTime = wakeUpTime
	timerFired := false
	for !timerFired && ctx.Err() == nil {
		timerCtx, timerCancel := workflow.WithCancel(ctx)
		duration := u.wakeUpTime.Sub(workflow.Now(timerCtx))
		timer := workflow.NewTimer(timerCtx, duration)
		logger.Info("SleepUntil", "WakeUpTime", u.wakeUpTime)
		workflow.NewSelector(timerCtx).
			AddFuture(timer, func(f workflow.Future) {
				err = f.Get(timerCtx, nil)
				// if a timer returned an error then it was canceled
				if err == nil {
					logger.Info("Timer fired")
					timerFired = true
				} else if ctx.Err() != nil { // Only log on root ctx cancellation, not on timerCancel function call.
					logger.Info("SleepUntil canceled")
				}
			}).
			AddReceive(updateWakeUpTimeCh, func(c workflow.ReceiveChannel, more bool) {
				timerCancel()                      // cancel outstanding timer
				c.Receive(timerCtx, &u.wakeUpTime) // update wake-up time
				logger.Info("Wake up time update requested")
			}).
			Select(timerCtx)
	}
	return ctx.Err()
}

func (u *UpdatableTimer) GetWakeUpTime() time.Time {
	return u.wakeUpTime
}
