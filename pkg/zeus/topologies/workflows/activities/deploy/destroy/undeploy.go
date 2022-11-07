package destroy_deploy

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/zeus-fyi/olympus/pkg/zeus/topologies/workflows/activities"
)

const destroyDeployRoute = "/v1/internal/deploy/destroy"

type UndeployTopologyActivity struct {
	activities.TopologyActivity
}

func (d *UndeployTopologyActivity) postDestroyDeployTarget(target string) error {
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

func (d *UndeployTopologyActivity) GetDestroyDeployURL(target string) url.URL {
	return d.GetURL(destroyDeployRoute, target)
}
