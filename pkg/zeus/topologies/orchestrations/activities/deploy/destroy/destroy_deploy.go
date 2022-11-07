package destroy_deploy_activities

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

const destroyDeployRoute = "/v1/internal/deploy/destroy"

type DestroyDeployTopologyActivity struct {
	activities.TopologyActivity
	BaseRoute string
}
type ActivityDefinition func(ctx context.Context) error
type ActivitiesSlice []ActivityDefinition

func (d *DestroyDeployTopologyActivity) GetActivities() ActivitiesSlice {
	return []ActivityDefinition{
		d.DestroyNamespace,
		d.DestroyDeployStatefulSet, d.DestroyDeployDeployment,
		d.DeployConfigMap,
		d.DestroyDeployService, d.DestroyDeployIngress,
	}
}

func (d *DestroyDeployTopologyActivity) postDestroyDeployTarget(target string) error {
	u := d.GetDestroyDeployURL(target)
	client := resty.New()
	_, err := client.R().
		SetBody(d.TopologyActivity).
		Post(u.Path)
	if err != nil {
		return err
	}
	return err
}

func (d *DestroyDeployTopologyActivity) GetDestroyDeployURL(target string) url.URL {
	return d.GetURL(destroyDeployRoute, target)
}
