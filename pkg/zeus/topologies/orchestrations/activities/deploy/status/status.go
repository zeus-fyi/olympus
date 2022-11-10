package deployment_status

import (
	"context"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	create_topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/state"
)

const updateDeployStatusRoute = "/v1/internal/deploy/status"

type TopologyActivityDeploymentStatusActivity struct {
	Host   string
	Bearer string
	create_topology_deployment_status.DeploymentStatus
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	_, err := client.R().
		SetAuthToken(d.Bearer).
		SetBody(d.DeploymentStatus).
		Post(u.Path)
	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) GetDeploymentStatusUpdateURL() url.URL {
	return d.GetURL(updateDeployStatusRoute)
}

func (d *TopologyActivityDeploymentStatusActivity) GetURL(target string) url.URL {
	if len(d.Host) <= 0 {
		d.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: d.Host,
		Path: path.Join(d.Host, target),
	}
	return u
}
