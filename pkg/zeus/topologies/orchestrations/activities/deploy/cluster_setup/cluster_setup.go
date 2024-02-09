package deploy_topology_activities_create_setup

import kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"

type CreateSetupTopologyActivities struct {
	Host string
	kronos_helix.KronosActivities
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (c *CreateSetupTopologyActivities) GetActivities() ActivitiesSlice {
	actSlice := []interface{}{
		c.AddDiskResourcesToOrg,
		c.AddDomainRecord,
		c.SendEmailNotification,
		c.MakeNodePoolRequest,
		c.AddNodePoolToOrgResources,
		c.AddAuthCtxNsOrg,
		c.RemoveAuthCtxNsOrg,
		c.RemoveNodePoolRequest,
		c.RemoveFreeTrialOrgResources,
		c.UpdateFreeTrialOrgResourcesToPaid,
		c.RemoveDomainRecord,
		c.SelectFreeTrialNodes,
		c.SelectNodeResources,
		c.EndResourceService,
		c.SelectDiskResourcesAtCloudCtxNs,
		c.GkeMakeNodePoolRequest,
		c.GkeRemoveNodePoolRequest,
		c.SelectGkeNodeResources,
		c.GkeSelectFreeTrialNodes,
		c.GkeAddNodePoolToOrgResources,
		c.EksMakeNodePoolRequest,
		c.PrivateEksMakeNodePoolRequest,
		c.EksAddNodePoolToOrgResources,
		c.SelectEksNodeResources,
		c.EksRemoveNodePoolRequest,
		c.PrivateEksRemoveNodePoolRequest,
		c.EksSelectFreeTrialNodes,
		c.OvhAddNodePoolToOrgResources,
		c.OvhMakeNodePoolRequest,
		c.SelectOvhNodeResources,
		c.OvhRemoveNodePoolRequest,
		c.OvhSelectFreeTrialNodes,
		c.GetClusterTopologyIds,
		c.GetClusterAuthCtx,
	}
	return actSlice
}
