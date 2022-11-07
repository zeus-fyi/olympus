package deploy_topology_activities

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

const deployRoute = "/v1/internal/deploy"

type DeployTopologyActivities struct {
	topology_activities.TopologyActivity
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.CreateNamespace,
		d.DeployDeployment, d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService, d.DeployIngress,
	}
}

func (d *DeployTopologyActivities) postDeployTarget(target string) error {
	u := d.GetDeployURL(target)
	client := resty.New()
	_, err := client.R().
		SetAuthToken(d.Bearer).
		SetBody(d.TopologyActivity).
		Post(u.Path)
	if err != nil {
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(deployRoute, target)
}
