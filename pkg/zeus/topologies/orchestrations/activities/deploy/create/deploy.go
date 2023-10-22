package deploy_topology_activities

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

type DeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
	kronos_helix.KronosActivities
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DeployTopologyActivities) GetActivities() ActivitiesSlice {
	actSlice := []interface{}{
		d.CreateNamespace,
		d.DeployDeployment,
		d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService,
		d.DeployIngress,
		d.DeployClusterTopology,
		d.DeployServiceMonitor,
		d.CreateChoreographySecret,
		d.CreateSecret,
		d.CreateJob,
		d.CreateCronJob,
	}
	return actSlice
}

func (d *DeployTopologyActivities) postDeployTarget(target string, params base_request.InternalDeploymentActionRequest) error {
	u := d.GetDeployURL(target)
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(params).
		Post(u.Path)

	if err != nil {
		log.Err(err).Interface("path", u.Path).Interface("err", err).Msg(fmt.Sprintf("DeployTopologyActivities: postDeployTarget failed status code %d err %s", resp.StatusCode(), err.Error()))
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		err = errors.New(fmt.Sprintf("DeployTopologyActivities: postDeployTarget failed bad status code %d err %s", resp.StatusCode(), resp.String()))
		log.Err(err).Interface("path", u.Path).Interface("err", err).Interface("statusCode", resp.StatusCode()).Msg("DeployTopologyActivities: postDeployTarget failed with bad status code")
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDeployPath, target)
}

func (d *DeployTopologyActivities) postDeployClusterTopology(params zeus_req_types.TopologyDeployRequest, ou org_users.OrgUser) error {
	if len(d.Host) <= 0 {
		d.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: d.Host,
	}

	token, err := auth.FetchUserAuthToken(context.Background(), ou)
	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: FetchUserAuthToken failed")
		return err
	}
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(token.PublicKey).
		SetBody(params).
		Post(zeus_endpoints.DeployTopologyV1Path)

	if err != nil {
		log.Err(err).Interface("path", u.Path).Interface("err", err).Msg("DeployTopologyActivities: postDeployClusterTopology failed")
		return err
	}
	if resp.StatusCode() != http.StatusAccepted {
		err = errors.New("DeployTopologyActivities: postDeployClusterTopology failed bad status code")
		log.Err(err).Interface("path", u.Path).Interface("statusCode", resp.StatusCode()).Msg("DeployTopologyActivities: postDeployClusterTopology failed")
		return err
	}
	return err
}
