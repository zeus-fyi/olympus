package deploy_topology_activities_create_setup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

func (s *DeployTestSuite) TestDeployCluster() {
	ctx := context.Background()
	act := CreateSetupTopologyActivities{}
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.NewOrgUserWithID(1699610911300434000, 1699610911367404000)
	dr := zeus_req_types.ClusterTopologyDeployRequest{
		ClusterClassName:    "microservice",
		SkeletonBaseOptions: []string{"api"},
		CloudCtxNs: zeus_common_types.CloudCtxNs{
			CloudProvider: "gcp",
			Region:        "us-central1",
			Context:       "gke_zeusfyi_us-central1-a_zeus-gcp-pilot-0",
			Namespace:     "microservice-6bad551f",
			Alias:         "microservice-6bad551f",
		},
	}
	err := act.postDeployClusterTopology(ctx, dr, ou)
	s.Require().Nil(err)
	//authed, err := read_topology.IsOrgCloudCtxNsAuthorized(ctx, ou.OrgID, dr.CloudCtxNs)
	//s.Require().Nil(err)
	//if authed != true || err != nil {
	//
	//}

	cl, err := read_topology.SelectClusterTopology(ctx, ou.OrgID, dr.ClusterClassName, dr.SkeletonBaseOptions)
	s.Require().Nil(err)

	log.Info().Interface("cl", cl).Msg("DeployClusterTopology: SelectClusterTopology")
	//clDeploy := base_deploy_params.ClusterTopologyWorkflowRequest{
	//	ClusterClassName:          t.ClusterClassName,
	//	TopologyIDs:               cl.GetTopologyIDs(),
	//	CloudCtxNS:                t.CloudCtxNs,
	//	OrgUser:                   ou,
	//	RequestChoreographySecret: cl.CheckForChoreographyOption(),
	//	AppTaint:                  t.AppTaint,
	//}
}

//func (c *CreateSetupTopologyActivities) postDeployClusterTopology(ctx context.Context, params zeus_req_types.ClusterTopologyDeployRequest, ou org_users.OrgUser) error {
//	//ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
//	//defer cancel()
//	if len(c.Host) <= 0 {
//		c.Host = "https://api.zeus.fyi"
//	}
//	u := url.URL{
//		Host: c.Host,
//	}
//	token, err := auth.FetchUserAuthToken(context.Background(), ou)
//	if err != nil {
//		log.Err(err).Interface("path", u.Path).Interface("ou", ou).Msg("CreateSetupTopologyActivities: FetchUserAuthToken failed")
//		return err
//	}
//	if len(token.PublicKey) <= 0 {
//		log.Err(err).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: FetchUserAuthToken failed")
//		return err
//	}
//	client := resty.New()
//	client.SetBaseURL(u.Host)
//	resp, err := client.R().
//		SetAuthToken(token.PublicKey).
//		SetBody(params).
//		Post(zeus_endpoints.DeployClusterTopologyV1Path)
//
//	if err != nil {
//		log.Err(err).Interface("path", u.Path).Msg("CreateSetupTopologyActivities: postDeployClusterTopology failed")
//		return err
//	}
//	if resp.StatusCode() != http.StatusAccepted {
//		log.Err(err).Interface("path", u.Path).Interface("statusCode", resp.StatusCode()).Msg("CreateSetupTopologyActivities: postDeployClusterTopology failed with bad status code")
//		return errors.New("CreateSetupTopologyActivities: postDeployClusterTopology failed")
//	}
//	return err
//}
