package topology_worker

import (
	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	clean_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/clean"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	deploy_workflow_cluster_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create_setup"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
)

func InitTopologyWorker(temporalAuthCfg temporal_auth.TemporalAuth) {
	tc, err := temporal_base.NewTemporalClient(temporalAuthCfg)
	if err != nil {
		log.Err(err).Msg("InitTopologyWorker: NewTemporalClient failed")
		misc.DelayedPanic(err)
	}
	taskQueueName := "TopologyTaskQueue"
	w := temporal_base.NewWorker(taskQueueName)
	// status
	statusActivity := deployment_status.TopologyActivityDeploymentStatusActivity{}
	// deploy destroy
	deployDestroyWf := destroy_deployed_workflow.NewDestroyDeployTopologyWorkflow()
	// deploy create
	deployWf := deploy_workflow.NewDeployTopologyWorkflow()
	// deploy clean
	cleanDeployWf := clean_deployed_workflow.NewCleanDeployTopologyWorkflow()
	// deploy create setup
	deployWfClusterSetup := deploy_workflow_cluster_setup.NewDeployCreateSetupTopologyWorkflow()

	// workflows added
	w.AddWorkflows(deployWf.GetWorkflows())
	w.AddWorkflow(deployDestroyWf.GetWorkflow())
	w.AddWorkflow(cleanDeployWf.GetWorkflow())
	w.AddWorkflows(deployWfClusterSetup.GetWorkflows())

	// activities added
	w.AddActivities(statusActivity.GetActivities())
	w.AddActivities(deployWf.GetActivities())
	w.AddActivities(deployDestroyWf.GetActivities())
	w.AddActivities(deployWfClusterSetup.GetActivities())

	Worker = TopologyWorker{w}
	Worker.TemporalClient = tc
	return
}
