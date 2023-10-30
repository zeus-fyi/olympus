package platform_service_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
)

func (h *HestiaPlatformServicesWorker) ExecuteIrisServerlessResyncWorkflow(ctx context.Context) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisServerlessResyncWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisServerlessResyncWorkflow")
		return err
	}
	return err
}
