package deployment_status

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/auth"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/deploy/workload_state"
)

const updateDeployStatusRoute = "/v1/internal/deploy/status"

type TopologyActivityDeploymentStatusActivity struct {
	Host string
	workload_state.InternalWorkloadStatusUpdateRequest
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context, status workload_state.InternalWorkloadStatusUpdateRequest) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	_, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(status).
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
		Path: target,
	}
	return u
}
