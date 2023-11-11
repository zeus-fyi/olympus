package deployment_status

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

type TopologyActivityDeploymentStatusActivity struct {
	Host string
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *TopologyActivityDeploymentStatusActivity) GetActivities() ActivitiesSlice {
	return []interface{}{d.PostStatusUpdate, d.CreateOrUpdateKubeCtxNsStatus, d.DeleteKubeCtxNsStatus}
}

func (d *TopologyActivityDeploymentStatusActivity) PostStatusUpdate(ctx context.Context, status topology_deployment_status.DeployStatus) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(status).
		Post(zeus_endpoints.InternalDeployStatusUpdatePath)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployStatusUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) CreateOrUpdateKubeCtxNsStatus(ctx context.Context, topDepReq zeus_req_types.TopologyDeployRequest) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(topDepReq).
		Post(zeus_endpoints.InternalDeployKnsCreateOrUpdatePath)
	if err != nil {
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployKnsCreateOrUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	if resp != nil && resp.StatusCode() >= 400 {
		err = fmt.Errorf("non-OK status code: %d", resp.StatusCode())
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployKnsCreateOrUpdatePath).Msg("TopologyActivityDeploymentStatusActivity")
		return err
	}
	return err
}

func (d *TopologyActivityDeploymentStatusActivity) DeleteKubeCtxNsStatus(ctx context.Context, topDepReq zeus_req_types.TopologyDeployRequest) error {
	u := d.GetDeploymentStatusUpdateURL()
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(topDepReq).
		Post(zeus_endpoints.InternalDeployKnsDestroyPath)
	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Err(err).Interface("path", zeus_endpoints.InternalDeployKnsDestroyPath).Msg("TopologyActivityDeploymentStatusActivity")
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
