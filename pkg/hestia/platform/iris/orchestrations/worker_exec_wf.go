package platform_service_orchestrations

import (
	"context"

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
		TaskQueue: h.TaskQueueName,
	}
	txWf := NewHestiaPlatformServiceWorkflows()
	wf := txWf.IrisRoutingServiceRequestWorkflow
	_, err := tc.ExecuteWorkflow(ctx, workflowOptions, wf, pr)
	if err != nil {
		log.Err(err).Msg("ExecuteIrisPlatformSetupRequestWorkflow")
		return err
	}
	return err
}
