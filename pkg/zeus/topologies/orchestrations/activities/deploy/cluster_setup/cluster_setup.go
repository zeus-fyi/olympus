package deploy_topology_activities_create_setup

import base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"

type CreateSetupTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (c *CreateSetupTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		c.AddDiskResourcesToOrg,
		c.AddDomainRecord,
		c.SendEmailNotification,
		c.MakeNodePoolRequest,
		c.AddNodePoolToOrgResources,
		c.AddAuthCtxNsOrg,
	}
}
