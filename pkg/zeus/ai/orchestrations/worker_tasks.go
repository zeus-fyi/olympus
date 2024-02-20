package ai_platform_service_orchestrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (z *ZeusAiPlatformServicesWorker) ExecuteRunAiWorkflowChildAnalysisProcess(ctx context.Context, cp *MbChildSubProcessParams) (*MbChildSubProcessParams, error) {
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("analysis-task-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RunAiChildAnalysisProcessWorkflow
	workflowRun, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, cp)
	if err != nil {
		log.Err(err).Msg("ExecuteRunAiWorkflowChildAnalysisProcess")
		return nil, err
	}
	err = workflowRun.Get(ctx, &cp)
	if err != nil {
		log.Err(err).Msg("ExecuteRunAiWorkflowChildAnalysisProcess: Get MbChildSubProcessParams")
		return nil, err
	}
	return cp, nil
}

func (z *ZeusAiPlatformServicesWorker) ExecuteRunAiChildAggAnalysisProcessWorkflow(ctx context.Context, cp *MbChildSubProcessParams) error {
	if cp == nil {
		return fmt.Errorf("cp is nil")
	}
	tc := z.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: z.TaskQueueName,
		ID:        fmt.Sprintf("agg-analysis-task-%s", uuid.New().String()),
	}
	txWf := NewZeusPlatformServiceWorkflows()
	wf := txWf.RunAiChildAggAnalysisProcessWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, cp)
	if err != nil {
		log.Err(err).Msg("ExecuteRunAiChildAggAnalysisProcessWorkflow")
		return err
	}

	return nil
}
