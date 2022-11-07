package topology_worker

import (
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deploy_topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/create"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
)

var Worker TopologyWorker

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
	tw.RegisterWorkflow(deployWf)
	deployActivities := deploy_topology_activities.DeployTopologyActivity{}
	tw.RegisterActivity(deployActivities)
	return tw, err
}
