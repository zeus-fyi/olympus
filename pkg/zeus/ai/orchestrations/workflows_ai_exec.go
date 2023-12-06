package ai_platform_service_orchestrations

import (
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (z *ZeusAiPlatformServiceWorkflows) RunAiWorkflowProcess(ctx workflow.Context, wfID string, ou org_users.OrgUser, wfExecParams artemis_orchestrations.WorkflowExecParams) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
	}
	timer := UpdatableTimer{}
	err := workflow.SetQueryHandler(ctx, QueryType, func() (time.Time, error) {
		return timer.GetWakeUpTime(), nil
	})
	if err != nil {
		logger.Error("failed to set query handler", "Error", err)
		return err
	}
	ojCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(ojCtx, z.UpsertAiOrchestration, ou, wfID, wfExecParams).Get(ojCtx, nil)
	if err != nil {
		logger.Error("failed to UpsertAiOrchestration", "Error", err)
		return err
	}
	err = timer.SleepUntil(ctx, wfExecParams.RunWindow.Start, workflow.GetSignalChannel(ctx, SignalType))
	if err != nil {
		logger.Error("failed to sleep until", "Error", err)
		return err
	}
	oj := artemis_orchestrations.NewActiveTemporalOrchestrationJobTemplate(ou.OrgID, wfID, wfExecParams.WorkflowTemplate.WorkflowGroup, wfExecParams.WorkflowTemplate.WorkflowName)
	startCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(startCtx, "UpdateAndMarkOrchestrationActive", oj).Get(startCtx, nil)
	if err != nil {
		logger.Error("failed to UpdateAndMarkOrchestrationActive", "Error", err)
		return err
	}
	aoAiAct := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 15, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    5,
		},
	}
	for i := 1; i < wfExecParams.RunCycles+1; i++ {
		startTime := wfExecParams.RunWindow.Start.Add(time.Duration(i) * wfExecParams.TimeStepSize)
		if time.Now().Before(startTime) {
			err = workflow.Sleep(ctx, startTime.Sub(time.Now()))
			if err != nil {
				logger.Error("failed to sleep", "Error", err)
				return err
			}
		}

		for _, analysisInst := range wfExecParams.WorkflowTasks {
			if i%analysisInst.AnalysisTaskID == 0 {
				// run retrieval
				retrievalCtx := workflow.WithActivityOptions(ctx, ao)
				window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.RunWindow.UnixStartTime, i-analysisInst.AnalysisCycleCount, i, wfExecParams.TimeStepSize)

				var content string
				err = workflow.ExecuteActivity(retrievalCtx, z.AiRetrievalTask, ou, analysisInst, window).Get(retrievalCtx, &content)
				if err != nil {
					logger.Error("failed to run retrieval", "Error", err)
					return err
				}
				if len(content) == 0 {
					continue
				}
				analysisCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				var aiResp openai.ChatCompletionResponse
				err = workflow.ExecuteActivity(analysisCtx, z.AiAnalysisTask, ou, analysisInst, content).Get(analysisCtx, &aiResp)
				if err != nil {
					logger.Error("failed to run analysis", "Error", err)
					return err
				}
				if len(aiResp.Choices) == 0 {
					continue
				}
				recordAnalysisCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordAnalysisCtx, z.SaveTaskOutput, ou, analysisInst, window, aiResp).Get(recordAnalysisCtx, nil)
				if err != nil {
					logger.Error("failed to save analysis", "Error", err)
					return err
				}
			}
		}

		for _, aggInst := range wfExecParams.WorkflowTasks {
			if aggInst.AggTaskID == nil || aggInst.AggCycleCount == nil || aggInst.AggPrompt == nil || aggInst.AggModel == nil {
				continue
			}
			aggCycle := wfExecParams.AggNormalizedCycleCounts[*aggInst.AggTaskID]
			if i%aggCycle == 0 {
				retrievalCtx := workflow.WithActivityOptions(ctx, ao)
				window := artemis_orchestrations.CalculateTimeWindowFromCycles(wfExecParams.RunWindow.UnixStartTime, i-aggCycle, i, wfExecParams.TimeStepSize)
				var content string
				err = workflow.ExecuteActivity(retrievalCtx, z.AiAggregateAnalysisRetrievalTask, ou, aggInst, window).Get(retrievalCtx, &content)
				if err != nil {
					logger.Error("failed to run retrieval", "Error", err)
					return err
				}
				if len(content) == 0 {
					continue
				}
				aggCtx := workflow.WithActivityOptions(ctx, aoAiAct)
				var aiAggResp openai.ChatCompletionResponse
				err = workflow.ExecuteActivity(aggCtx, z.AiAggregateTask, ou, aggInst, content).Get(aggCtx, &aiAggResp)
				if err != nil {
					logger.Error("failed to run aggregation", "Error", err)
					return err
				}
				if len(aiAggResp.Choices) == 0 {
					continue
				}
				recordAggCtx := workflow.WithActivityOptions(ctx, ao)
				err = workflow.ExecuteActivity(recordAggCtx, z.SaveTaskOutput, ou, aggInst, window, aiAggResp).Get(recordAggCtx, nil)
				if err != nil {
					logger.Error("failed to save aggregation resp", "Error", err)
					return err
				}
			}
		}
	}

	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
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
