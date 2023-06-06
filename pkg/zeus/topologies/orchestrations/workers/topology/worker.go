package topology_worker

import (
	"context"

	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	clean_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/clean"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	deploy_workflow_cluster_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create_setup"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
	deploy_workflow_destroy_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy_setup"
	deploy_workflow_cluster_updates "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/update"
	"go.temporal.io/sdk/client"
)

var Worker TopologyWorker

type TopologyWorker struct {
	temporal_base.Worker
}

func (t *TopologyWorker) ExecuteDeployFleetUpgrade(ctx context.Context, params base_deploy_params.FleetUpgradeWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	deployWf := deploy_workflow_cluster_updates.NewDeployFleetUpgradeWorkflow()
	wf := deployWf.UpgradeFleetWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDeployFleetUpgrade")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDeployCluster(ctx context.Context, params base_deploy_params.ClusterTopologyWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	wf := deployWf.DeployClusterTopologyWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDeployCluster")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDeploy(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	wf := deployWf.DeployTopologyWorkflow
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDeploy")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDestroyDeploy(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	deployDestroyWf := destroy_deployed_workflow.NewDestroyDeployTopologyWorkflow()
	wf := deployDestroyWf.DestroyDeployDeployment
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDestroyDeploy")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteCleanDeploy(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	cleanDeployWf := clean_deployed_workflow.NewCleanDeployTopologyWorkflow()
	wf := cleanDeployWf.GetWorkflow()
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteCleanDeploy")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteCreateSetupCluster(ctx context.Context, params base_deploy_params.ClusterSetupRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	clusterSetupWf := deploy_workflow_cluster_setup.NewDeployCreateSetupTopologyWorkflow()
	wf := clusterSetupWf.GetDeployClusterSetupWorkflow()
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteCreateSetupCluster")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDestroyResources(ctx context.Context, params base_deploy_params.DestroyResourcesRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	destroyResourcesWf := deploy_workflow_destroy_setup.NewDestroyResourcesWorkflow()
	wf := destroyResourcesWf.GetWorkflow()
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDestroyResources")
		return err
	}
	return err
}

func (t *TopologyWorker) ExecuteDestroyNamespace(ctx context.Context, params base_deploy_params.TopologyWorkflowRequest) error {
	c := t.ConnectTemporalClient()
	defer c.Close()
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: t.TaskQueueName,
	}
	deployDestroyWf := deploy_workflow_destroy_setup.NewDestroyNamespaceSetupWorkflow()
	wf := deployDestroyWf.GetWorkflow()
	_, err := c.ExecuteWorkflow(ctx, workflowOptions, wf, params)
	if err != nil {
		log.Err(err).Msg("ExecuteDestroyNamespace")
		return err
	}
	return err
}
