package deploy_topology_activities_create_setup

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	zeus_endpoints "github.com/zeus-fyi/olympus/pkg/zeus/client/endpoints"
	api_auth_temporal "github.com/zeus-fyi/olympus/pkg/zeus/topologies/orchestrations/orchestration_auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func (c *CreateSetupTopologyActivities) DeployClusterTopologyFromUI(ctx context.Context, clusterName string, sbBases []string, cloudCtxNs zeus_common_types.CloudCtxNs, ou org_users.OrgUser) error {
	cdRequest := zeus_req_types.ClusterTopologyDeployRequest{
		ClusterClassName:    clusterName,
		SkeletonBaseOptions: sbBases,
		CloudCtxNs:          cloudCtxNs,
		AppTaint:            true,
	}
	return c.postDeployClusterTopology(ctx, cdRequest, ou)
}

func (c *CreateSetupTopologyActivities) DestroyCluster(ctx context.Context, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	return c.destroyClusterTopology(cloudCtxNs)
}
func (c *CreateSetupTopologyActivities) destroyClusterTopology(cloudCtxNs zeus_common_types.CloudCtxNs) error {
	if len(c.Host) <= 0 {
		c.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: c.Host,
	}

	params := zeus_req_types.TopologyDeployRequest{
		TopologyID: 0,
		CloudCtxNs: cloudCtxNs,
	}
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(api_auth_temporal.Bearer).
		SetBody(params).
		Post(zeus_endpoints.DestroyDeployInfraV1Path)

	if err != nil || resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: destroyClusterTopology failed")
		return err
	}
	if resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: destroyClusterTopology failed")
		return errors.New("CreateSetupTopologyActivities: destroyClusterTopology failed")
	}
	return err
}

func (c *CreateSetupTopologyActivities) postDeployClusterTopology(ctx context.Context, params zeus_req_types.ClusterTopologyDeployRequest, ou org_users.OrgUser) error {
	//ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	//defer cancel()
	if len(c.Host) <= 0 {
		c.Host = "https://api.zeus.fyi"
	}
	u := url.URL{
		Host: c.Host,
	}

	token, err := auth.FetchUserAuthToken(context.Background(), ou)
	if err == pgx.ErrNoRows {
		key, err2 := create_keys.CreateUserAPIKey(ctx, ou)
		if err2 != nil {
			log.Err(err2).Msg("CreateUserAPIKey error")
			return err2
		}
		token.PublicKey = key.PublicKey
	}
	if err != nil {
		log.Err(err).Interface("params", params).Interface("path", u.Path).Interface("ou", ou).Msg("CreateSetupTopologyActivities: FetchUserAuthToken failed")
		return err
	}
	if len(token.PublicKey) <= 0 {
		log.Err(err).Interface("params", params).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: FetchUserAuthToken failed zero length")
		return err
	}
	client := resty.New()
	client.SetBaseURL(u.Host)
	resp, err := client.R().
		SetAuthToken(token.PublicKey).
		SetBody(params).
		Post(zeus_endpoints.DeployClusterTopologyV1Path)

	if err != nil {
		log.Err(err).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: postDeployClusterTopology failed")
		return err
	}
	if resp.StatusCode() != http.StatusAccepted {
		log.Err(err).Interface("path", u.Path).Interface("statusCode", resp.StatusCode()).Msg("CreateSetupTopologyActivities: postDeployClusterTopology failed with bad status code")
		return errors.New("CreateSetupTopologyActivities: postDeployClusterTopology failed")
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
