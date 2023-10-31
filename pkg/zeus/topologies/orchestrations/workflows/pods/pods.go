package pods_workflows

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	orchestrate_pods_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/pods"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type PodsWorkflows struct {
	temporal_base.Workflow
	orchestrate_pods_activities.PodsActivity
}

const defaultTimeout = 120 * time.Minute

func NewPodsWorkflows() PodsWorkflows {
	deployWf := PodsWorkflows{
		Workflow:     temporal_base.Workflow{},
		PodsActivity: orchestrate_pods_activities.PodsActivity{},
	}
	return deployWf
}

func (p *PodsWorkflows) GetWorkflows() []interface{} {
	return []interface{}{p.DeletePodWorkflow}
}

func (p *PodsWorkflows) DeletePodWorkflow(ctx workflow.Context, wfId, podName string, cctx zeus_common_types.CloudCtxNs, delay time.Duration) error {
	logger := workflow.GetLogger(ctx)
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second * 60,
		BackoffCoefficient: 1.0,
		MaximumInterval:    time.Minute * 10,
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: defaultTimeout,
		RetryPolicy:         retryPolicy,
	}
	err := workflow.Sleep(ctx, delay)
	if err != nil {
		return err
	}

	dpCtx := workflow.WithActivityOptions(ctx, ao)
	err = workflow.ExecuteActivity(dpCtx, "", podName, cctx).Get(dpCtx, nil)
	if err != nil {
		logger.Error("Failed to delete pod", "Error", err)
		return err
	}
	return nil
}

func (t *PodsWorker) ExecuteDeletePodWorkflow(ctx context.Context, cctx zeus_common_types.CloudCtxNs, podName string, delay time.Duration) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
		ID:        uuid.New().String(),
	}
	podsWfs := NewPodsWorkflows()
	wf := podsWfs.DeletePodWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, cctx, podName, delay)
	if err != nil {
		log.Err(err).Msg("DeletePodWorkflow")
		return err
	}
	return nil
}

var PodsServiceWorker PodsWorker

type PodsWorker struct {
	temporal_base.Worker
}

func InitPodsWorker(temporalAuthCfg temporal_auth.TemporalAuth) {
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitTopologyWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := "PodsTaskQueue"
	w := temporal_base.NewWorker(taskQueueName)
	podsWfs := NewPodsWorkflows()
	w.AddWorkflows(podsWfs.GetWorkflows())
	w.AddActivities(podsWfs.GetActivities())
	PodsServiceWorker = PodsWorker{w}
	PodsServiceWorker.TemporalClient = tc

	return
}

func ExecuteDeletePodWorkflow(c echo.Context, ctx context.Context, cctx zeus_common_types.CloudCtxNs, podName string, delay time.Duration) error {
	err := PodsServiceWorker.ExecuteDeletePodWorkflow(ctx, cctx, podName, delay)
	if err != nil {
		log.Err(err).Msg("ExecuteDeletePodWorkflow, ExecuteWorkflow error")
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusAccepted, nil)
}
