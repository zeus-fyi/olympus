package deploy_topology_activities

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

const deployRoute = "/v1/internal/deploy"

type DeployTopologyActivity struct {
	activities.TopologyActivity
}
type ActivityDefinition func(ctx context.Context) error
type ActivitiesSlice []ActivityDefinition

func (d *DeployTopologyActivity) GetActivities() ActivitiesSlice {
	return []ActivityDefinition{
		d.CreateNamespace,
		d.DeployDeployment, d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService, d.DeployIngress,
	}
}

func (d *DeployTopologyActivity) postDeployTarget(target string) error {
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

func (d *DeployTopologyActivity) GetDeployURL(target string) url.URL {
	return d.GetURL(deployRoute, target)
}
