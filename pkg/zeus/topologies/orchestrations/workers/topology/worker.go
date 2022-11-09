package topology_worker

import (
	"context"

	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
	"go.temporal.io/sdk/client"
)

var Worker TopologyWorker

type TopologyWorker struct {
	temporal_base.Worker
}

func (t *TopologyWorker) ExecuteDeploy(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	err := t.Connect()
	if err != nil {
		return err
	}
	defer t.Close()

	workflowOptions := client.StartWorkflowOptions{}
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	wf := deployWf.GetWorkflow()
	_, err = t.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDestroyDeploy(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	err := t.Connect()
	if err != nil {
		return err
	}
	defer t.Close()

	workflowOptions := client.StartWorkflowOptions{}
	deployDestroyWf := destroy_deployed_workflow.NewDestroyDeployTopologyWorkflow()
	wf := deployDestroyWf.GetWorkflow()
	_, err = t.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		return err
	}
	return err
}
