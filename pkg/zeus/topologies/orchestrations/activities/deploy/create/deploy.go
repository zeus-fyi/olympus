package deploy_topology_activities

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
)

type DeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.CreateNamespace,
		d.DeployDeployment,
		d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService,
		d.DeployIngress,
		d.DeployClusterTopology,
		d.DeployServiceMonitor,
		d.CreateChoreographySecret,
	}
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
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployTarget failed")
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployTarget failed")
		return errors.New("DeployTopologyActivities: postDeployTarget failed")
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
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployClusterTopology failed")
		return err
	}
	if resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployClusterTopology failed")
		return errors.New("DeployTopologyActivities: postDeployClusterTopology failed")
	}
	return err
}
