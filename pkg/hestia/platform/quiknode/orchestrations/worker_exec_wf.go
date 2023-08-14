package quicknode_orchestrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	"go.temporal.io/sdk/client"
)

func (h *HestiaQuickNodeWorker) ExecuteQnProvisionWorkflow(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest, user hestia_quicknode.QuickNodeUserInfo) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuickNodeWorkflow()
	wf := txWf.ProvisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, ou, pr, user)
	if err != nil {
		log.Err(err).Msg("ExecuteQnProvisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuickNodeWorker) ExecuteQnUpdateProvisionWorkflow(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.ProvisionRequest) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuickNodeWorkflow()
	wf := txWf.UpdateProvisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, ou, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteQnUpdateProvisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuickNodeWorker) ExecuteQnDeprovisionWorkflow(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.DeprovisionRequest) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuickNodeWorkflow()
	wf := txWf.DeprovisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, ou, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteQnDeprovisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuickNodeWorker) ExecuteQnDeactivateWorkflow(ctx context.Context, ou org_users.OrgUser, pr hestia_quicknode.DeactivateRequest) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuickNodeWorkflow()
	wf := txWf.DeactivateWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, ou, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteQnDeactivateWorkflow")
		return err
	}
	return err
}
