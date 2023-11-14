package deploy_topology_activities

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	kronos_helix "github.com/zeus-fyi/olympus/pkg/kronos/helix"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	base_deploy_params "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/workflows/deploy/base"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/temporal_actions/base_request"
	"github.com/zeus-fyi/zeus/zeus/workload_config_drivers/topology_workloads"
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
		d.DeployServiceMonitor,
		d.CreateChoreographySecret,
		d.CreateSecret,
		d.CreateJob,
		d.CreateCronJob,
		d.GetTopologyInfraConfig,
	}
	return actSlice
}

func (d *DeployTopologyActivities) GetTopologyInfraConfig(ctx context.Context, ou org_users.OrgUser, topID int) (*topology_workloads.TopologyBaseInfraWorkload, error) {
	tr := read_topology.NewInfraTopologyReaderWithOrgUser(ou)
	tr.TopologyID = topID
	err := tr.SelectTopology(ctx)
	if err != nil {
		log.Err(err).Interface("orgUser", ou).Int("topID", topID).Msg("DeployTopology, ReadUserTopologyConfig error")
		return nil, err
	}
	chartWkLoad := tr.GetTopologyBaseInfraWorkload()
	return &chartWkLoad, nil
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
	if resp != nil && resp.StatusCode() >= 400 {
		err = errors.New(fmt.Sprintf("DeployTopologyActivities: postDeployTarget failed bad status code %d err %s", resp.StatusCode(), resp.String()))
		log.Err(err).Interface("path", u.Path).Interface("err", err).Interface("statusCode", resp.StatusCode()).Msg("DeployTopologyActivities: postDeployTarget failed with bad status code")
		return err
	}
	return err
}

func (d *DeployTopologyActivities) GetDeployURL(target string) url.URL {
	return d.GetURL(zeus_endpoints.InternalDeployPath, target)
}

const (
	internalOrgID = 7138983863666903883
)
