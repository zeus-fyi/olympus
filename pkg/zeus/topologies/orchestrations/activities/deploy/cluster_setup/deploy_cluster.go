package deploy_topology_activities_create_setup

import (
	"context"
	"net/http"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types"
)

func (c *CreateSetupTopologyActivities) DeployClusterTopology(ctx context.Context, params zeus_req_types.ClusterTopologyDeployRequest, ou org_users.OrgUser) error {
	return c.postDeployClusterTopology(params, ou)
}

func (c *CreateSetupTopologyActivities) GetDeployURL(target string) url.URL {
	return c.GetURL(zeus_endpoints.InternalDeployPath, target)
}

func (c *CreateSetupTopologyActivities) postDeployClusterTopology(params zeus_req_types.ClusterTopologyDeployRequest, ou org_users.OrgUser) error {
	if len(c.Host) <= 0 {
		c.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: c.Host,
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
		Post(zeus_endpoints.DeployClusterTopologyV1Path)

	if err != nil || resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Interface("path", u.Path).Msg("DeployTopologyActivities: postDeployClusterTopology failed")
		return err
	}
	return err
}

func (c *CreateSetupTopologyActivities) GetURL(prefix, target string) url.URL {
	if len(c.Host) <= 0 {
		c.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: c.Host,
		Path: path.Join(prefix, target),
	}
	return u
}
