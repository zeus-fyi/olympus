package topology

import (
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
)

type TopologyWorker struct {
	temporal_base.Worker
}

func NewTopologyWorker(authCfg temporal_base.TemporalAuth) (TopologyWorker, error) {
	tc, err := temporal_base.NewTemporalClient(authCfg)
	if err != nil {
		return TopologyWorker{}, err
	}
	w := temporal_base.NewWorker(tc, "TopologyTaskQueue")
	tw := TopologyWorker{w}

	deployWf := deploy_workflow.DeployTopologyWorkflow{}
	deployWfActivities := deployWf.GetActivities()
	tw.AddActivitiesToWorker(deployWfActivities)

	destroyDeployWf := destroy_deployed_workflow.DestroyDeployTopologyWorkflow{}
	destroyDeployWfActivities := destroyDeployWf.GetActivities()
	tw.AddActivitiesToWorker(destroyDeployWfActivities)

	tw.AddWorkflowToWorker(deployWf.Workflow, destroyDeployWf.Workflow)
	return tw, err
}
