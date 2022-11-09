package topology_worker

import (
	"github.com/rs/zerolog/log"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
)

func InitTopologyWorker(authCfg temporal_base.TemporalAuth) (TopologyWorker, error) {
	w, err := NewTopologyWorker(authCfg)
	if err != nil {
		log.Err(err).Msg("InitTopologyWorker failed")
	}
	Worker = w
	return w, err
}

func NewTopologyWorker(authCfg temporal_base.TemporalAuth) (TopologyWorker, error) {
	tc, err := temporal_base.NewTemporalClient(authCfg)
	if err != nil {
		return TopologyWorker{}, err
	}
	err = tc.Connect()
	if err != nil {
		return TopologyWorker{}, err
	}
	defer tc.Close()

	taskQueueName := "TopologyTaskQueue"
	w := temporal_base.NewWorker(tc, "TopologyTaskQueue")

	// status
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{}
	w.AddActivities(statusActivity.GetActivities())

	// deploy create
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	w.AddWorkflow(deployWf.GetWorkflow())
	w.AddActivities(deployWf.GetActivities())

	// deploy destroy
	deployDestroyWf := destroy_deployed_workflow.NewDestroyDeployTopologyWorkflow()
	w.AddWorkflow(deployDestroyWf.GetWorkflow())
	w.AddActivities(deployDestroyWf.GetActivities())
	err = w.RegisterWorker()
	if err != nil {
		return TopologyWorker{}, err
	}
	tw := TopologyWorker{w}
	tw.TaskQueueName = taskQueueName
	return tw, err
}
