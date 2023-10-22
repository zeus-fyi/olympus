package topology_worker

import (
	"github.com/rs/zerolog/log"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	temporal_base "github.com/zeus-fyi/olympus/pkg/iris/temporal/base"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	deployment_status "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities/deploy/status"
	clean_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/clean"
	deploy_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create"
	deploy_workflow_cluster_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/create_setup"
	destroy_deployed_workflow "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy"
	deploy_workflow_destroy_setup "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/destroy_setup"
	deploy_workflow_cluster_updates "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/update"
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
	// teardown setup
	deployDestroyWfClusterSetup := deploy_workflow_destroy_setup.NewDeployDestroyClusterSetupWorkflow()
	// resource removal
	resourcesDestroyWf := deploy_workflow_destroy_setup.NewDestroyResourcesWorkflow()
	// namespace destroy from ui
	deployDestroyNamespaceWf := deploy_workflow_destroy_setup.NewDestroyNamespaceSetupWorkflow()
	// fleet upgrades
	upgradeFleetWf := deploy_workflow_cluster_updates.NewDeployFleetUpgradeWorkflow()

	// workflows added
	w.AddWorkflows(deployWf.GetWorkflows())
	w.AddWorkflows(deployDestroyWf.GetWorkflows())
	w.AddWorkflow(cleanDeployWf.GetWorkflow())
	w.AddWorkflows(deployWfClusterSetup.GetWorkflows())
	w.AddWorkflow(deployDestroyWfClusterSetup.GetWorkflow())
	w.AddWorkflow(resourcesDestroyWf.GetWorkflow())
	w.AddWorkflow(deployDestroyNamespaceWf.GetWorkflow())
	w.AddWorkflows(upgradeFleetWf.GetWorkflows())
	// activities added
	kr := kronos_helix.NewKronosActivities()
	w.AddActivities(kr.GetActivities())
	w.AddActivities(statusActivity.GetActivities())
	w.AddActivities(deployWf.GetActivities())
	w.AddActivities(deployDestroyWf.GetActivities())
	w.AddActivities(deployWfClusterSetup.GetActivities())
	w.AddActivities(upgradeFleetWf.GetActivities())

	Worker = TopologyWorker{w}
	Worker.TemporalClient = tc
	return
}
