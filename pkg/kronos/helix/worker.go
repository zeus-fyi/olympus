package kronos_helix

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"go.temporal.io/sdk/client"
)

type KronosWorker struct {
	temporal_base.Worker
}

var (
	KronosServiceWorker KronosWorker
)

const (
	KronosHelixTaskQueue = "KronosHelixTaskQueue"
)

// ExecuteKronosWorkflow starts the endless cycles of the Kronos workflow
func (k *KronosWorker) ExecuteKronosWorkflow(ctx context.Context) error {
	tc := k.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: k.TaskQueueName,
	}
	txWf := NewKronosWorkflow()
	wf := txWf.Yin
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteKronosWorkflow")
		return err
	}
	return nil
}

func (k *KronosWorker) ExecuteMockingbirdWorkflow(ctx context.Context) error {
	tc := k.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: k.TaskQueueName,
	}
	txWf := NewKronosWorkflow()
	wf := txWf.Mockingbird
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf)
	if err != nil {
		log.Err(err).Msg("ExecuteMockingbirdWorkflow")
		return err
	}
	return nil
}
