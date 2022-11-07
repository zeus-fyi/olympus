package destroy_deploy_activities

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	topology_activities "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/activities"
)

const destroyDeployRoute = "/v1/internal/deploy/destroy"

type DestroyDeployTopologyActivity struct {
	topology_activities.TopologyActivity
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DestroyDeployTopologyActivity) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.DestroyNamespace,
		d.DestroyDeployStatefulSet, d.DestroyDeployDeployment,
		d.DestroyDeployConfigMap,
		d.DestroyDeployService, d.DestroyDeployIngress,
	}
}

func (d *DestroyDeployTopologyActivity) postDestroyDeployTarget(target string) error {
	u := d.GetDestroyDeployURL(target)
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

func (d *DestroyDeployTopologyActivity) GetDestroyDeployURL(target string) url.URL {
	return d.GetURL(destroyDeployRoute, target)
}
