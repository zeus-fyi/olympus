package deployment_status

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
)

type TopologyActivityDeploymentStatusActivity struct {
	Host string
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate, d.PostKnsStatusUpdate}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context, status topology_deployment_status.Status) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	_, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(status.DeployStatus).
		Post(zeus_endpoints.InternalDeployStatusUpdatePath)
	if err != nil {
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployStatusUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) PostKnsStatusUpdate(ctx context.Context, status topology_deployment_status.Status) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	_, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(status.TopologyKubeCtxNs).
		Post(zeus_endpoints.InternalDeployKnsStatusUpdatePath)
	if err != nil {
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployKnsStatusUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) GetDeploymentStatusUpdateURL() url.URL {
	return d.GetURL()
}

func (d *TopologyActivityDeploymentStatusActivity) GetURL() url.URL {
	if len(d.Host) <= 0 {
		d.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: d.Host,
	}
	return u
}
