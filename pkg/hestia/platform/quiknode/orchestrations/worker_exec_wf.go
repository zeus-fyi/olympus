package quicknode_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_quicknode "github.com/zeus-fyi/olympus/pkg/hestia/platform/quiknode"
	"go.temporal.io/sdk/client"
)

func (h *HestiaQuicknodeWorker) ExecuteQnProvisionWorkflow(ctx context.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser, user hestia_quicknode.QuickNodeUserInfo) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuicknodeWorkflow()
	wf := txWf.ProvisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr, ou, user)
	if err != nil {
		log.Err(err).Msg("ExecuteQnProvisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuicknodeWorker) ExecuteQnUpdateProvisionWorkflow(ctx context.Context, pr hestia_quicknode.ProvisionRequest, ou org_users.OrgUser) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuicknodeWorkflow()
	wf := txWf.UpdateProvisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr, ou)
	if err != nil {
		log.Err(err).Msg("ExecuteQnUpdateProvisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuicknodeWorker) ExecuteQnDeprovisionWorkflow(ctx context.Context, pr hestia_quicknode.DeprovisionRequest, ou org_users.OrgUser) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuicknodeWorkflow()
	wf := txWf.DeprovisionWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr, ou)
	if err != nil {
		log.Err(err).Msg("ExecuteQnDeprovisionWorkflow")
		return err
	}
	return err
}

func (h *HestiaQuicknodeWorker) ExecuteQnDeactivateWorkflow(ctx context.Context, pr hestia_quicknode.DeactivateRequest, ou org_users.OrgUser) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaQuicknodeWorkflow()
	wf := txWf.DeactivateWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr, ou)
	if err != nil {
		log.Err(err).Msg("ExecuteQnDeactivateWorkflow")
		return err
	}
	return err
}
