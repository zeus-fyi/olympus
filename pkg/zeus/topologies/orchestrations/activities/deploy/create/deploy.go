package deploy_topology_activities

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
)

const deployRoute = "/v1/internal/deploy"

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

func (d *DeployTopologyActivities) postDeployTarget() error {
	u := d.GetDeployURL()
	client := resty.New()
	_, err := client.R().
		SetAuthToken(d.Bearer).
		SetBody(d.TopologyWorkflowRequest).
		Post(u.Path)
	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployTarget failed")
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL() url.URL {
	return d.GetURL(deployRoute)
}
