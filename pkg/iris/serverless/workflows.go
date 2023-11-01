package iris_serverless

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/sdk/workflow"
)

type IrisPlatformServiceWorkflows struct {
	temporal_base.Workflow
	IrisPlatformActivities
}

const defaultTimeout = 72 * time.Hour

func NewIrisPlatformServiceWorkflows() IrisPlatformServiceWorkflows {
	deployWf := IrisPlatformServiceWorkflows{
		Workflow: temporal_base.Workflow{},
	}
	return deployWf
}

func (i *IrisPlatformServiceWorkflows) GetWorkflows() []interface{} {
	return []interface{}{i.IrisServerlessResyncWorkflow, i.IrisServerlessPodRestartWorkflow}
}

func (i *IrisPlatformServiceWorkflows) IrisServerlessResyncWorkflow(ctx workflow.Context, wfID string) error {
	logger := workflow.GetLogger(ctx)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10, // Setting a valid non-zero timeout
	}
	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "IrisPlatformServiceWorkflows", "IrisServerlessResyncWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err := workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, i.ResyncServerlessRoutes, nil).Get(pCtx, nil)
	if err != nil {
		logger.Error("IrisPlatformServiceWorkflows: failed to ResyncServerlessRoutes", "Error", err)
		return err
	}
	finishedCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(finishedCtx, "UpdateAndMarkOrchestrationInactive", oj).Get(finishedCtx, nil)
	if err != nil {
		logger.Error("failed to update cache for qn services", "Error", err)
		return err
	}
	return nil
}

func (i *IrisPlatformServiceWorkflows) IrisServerlessPodRestartWorkflow(ctx workflow.Context, wfID string, orgID int, cctx zeus_common_types.CloudCtxNs, podName, serverlessTable, sessionID string, waitTime time.Time) error {
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
	err = timer.SleepUntil(ctx, waitTime, workflow.GetSignalChannel(ctx, SignalType))
	if err != nil {
		logger.Error("failed to sleep until", "Error", err)
		return err
	}

	oj := artemis_orchestrations.NewInternalActiveTemporalOrchestrationJobTemplate(wfID, "IrisPlatformServiceWorkflows", "IrisServerlessPodRestartWorkflow")
	alertCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(alertCtx, "UpsertAssignment", oj).Get(alertCtx, nil)
	if err != nil {
		logger.Error("failed to update QuickNode services", "Error", err)
		return err
	}
	cCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(cCtx, i.ClearServerlessSessionRouteCache, orgID, serverlessTable, sessionID).Get(cCtx, nil)
	if err != nil {
		logger.Error("IrisPlatformServiceWorkflows: failed to RestartServerlessPod", "Error", err)
		return err
	}
	pCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(pCtx, i.RestartServerlessPod, cctx, podName, 0).Get(pCtx, nil)
	if err != nil {
		logger.Error("IrisPlatformServiceWorkflows: failed to RestartServerlessPod", "Error", err)
		return err
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
