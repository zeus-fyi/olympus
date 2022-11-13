package deploy_topology_activities

import (
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/actions/base_request"
)

type DeployTopologyActivities struct {
	base_deploy_params.TopologyWorkflowRequest
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *DeployTopologyActivities) GetActivities() ActivitiesSlice {
	return []interface{}{
		d.CreateNamespace,
		d.DeployDeployment, d.DeployStatefulSet,
		d.DeployConfigMap,
		d.DeployService, d.DeployIngress,
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

	if err != nil || resp.StatusCode() != http.StatusOK {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployTarget failed")
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDeployPath, target)
}
