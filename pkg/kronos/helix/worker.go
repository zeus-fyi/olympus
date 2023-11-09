package kronos_helix

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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

func (k *KronosWorker) ExecuteAiTaskWorkflow(ctx context.Context, ou org_users.OrgUser, content string) error {
	tc := k.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: k.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewKronosWorkflow()
	wf := txWf.AiWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, ou, content)
	if err != nil {
		log.Err(err).Msg("ExecuteAiTaskWorkflow")
		return err
	}
	return nil
}

//// ExecuteKronosWorkflow2 starts the endless cycles of the Kronos workflow
//func (k *KronosWorker) ExecuteKronosWorkflow2(ctx context.Context, inst Instructions) error {
//	tc := k.ConnectTemporalClient()
//	defer tc.Close()
//	workflowOptions := client.StartWorkflowOptions{
//		TaskQueue: k.TaskQueueName,
//	}
//	txWf := NewKronosWorkflow()
//	wf := txWf.CronJob
//	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, inst, 2)
//	if err != nil {
//		log.Err(err).Msg("ExecuteKronosWorkflow")
//		return err
//	}
//	return nil
//}
