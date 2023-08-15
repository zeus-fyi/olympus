package platform_service_orchestrations

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"go.temporal.io/sdk/client"
)

type IrisPlatformServiceRequest struct {
	Ou           org_users.OrgUser
	OrgGroupName string
	Routes       []string
}

func (h *HestiaPlatformServicesWorker) ExecuteIrisPlatformSetupRequestWorkflow(ctx context.Context, pr IrisPlatformServiceRequest) error {
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisRoutingServiceRequestWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisRoutingServiceRequestWorkflow")
		return err
	}
	return err
}

func (h *HestiaPlatformServicesWorker) ExecuteIrisDeleteOrgGroupRoutingTableWorkflow(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if pr.OrgGroupName == "" {
		return errors.New("org group name is empty")
	}
	if pr.Ou.OrgID == 0 {
		return errors.New("org id is empty")
	}
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisDeleteOrgGroupRoutingTableWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisDeleteOrgGroupRoutingTableWorkflow")
		return err
	}
	return err
}

func (h *HestiaPlatformServicesWorker) ExecuteIrisDeleteOrgRoutesWorkflow(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if pr.Ou.OrgID == 0 {
		return errors.New("org id is empty")
	}
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisDeleteOrgRoutesWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisDeleteOrgRoutesWorkflow")
		return err
	}
	return err
}

func (h *HestiaPlatformServicesWorker) ExecuteIrisRemoveAllOrgRoutesFromCacheWorkflow(ctx context.Context, pr IrisPlatformServiceRequest) error {
	if pr.Ou.OrgID == 0 {
		return errors.New("org id is empty")
	}
	tc := h.ConnectTemporalClient()
	defer tc.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: h.TaskQueueName,
		ID:        uuid.New().String(),
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisRemoveAllOrgRoutesFromCacheWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, workflowOptions.ID, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisRemoveAllOrgRoutesFromCacheWorkflow")
		return err
	}
	return err
}
