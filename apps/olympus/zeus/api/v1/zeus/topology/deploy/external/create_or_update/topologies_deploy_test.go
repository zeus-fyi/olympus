package create_or_update_deploy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	beacon_cookbooks "github.com/zeus-fyi/zeus/cookbooks/ethereum/beacons"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

type TopologyDeployActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

//
//func (t *TopologyDeployActionRequestTestSuite) TestDeployCluster() {
//	t.InitLocalConfigs()
//	t.Eg.POST("/deploy/cluster", ClusterTopologyDeploymentHandler)
//
//	start := make(chan struct{}, 1)
//	go func() {
//		close(start)
//		_ = t.E.Start(":9010")
//	}()
//
//	<-start
//	ctx := context.Background()
//	defer t.E.Shutdown(ctx)
//
//	cd := zeus_req_types.ClusterTopologyDeployRequest{
//		ClusterClassName:    "ethereumEphemeralValidatorCluster",
//		SkeletonBaseOptions: []string{"gethHercules", "lighthouseHercules", "lighthouseHerculesValidatorClient"},
//		CloudCtxNs:          beacon_cookbooks.BeaconCloudCtxNs,
//	}
//	resp, err := t.ZeusClient.DeployCluster(ctx, cd)
//	t.Require().Nil(err)
//	t.Assert().NotEmpty(resp)
//}

func (t *TopologyDeployActionRequestTestSuite) TestReadCluster() {
	t.InitLocalConfigs()
	ctx := context.Background()
	//apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)

	cd := zeus_req_types.ClusterTopologyDeployRequest{
		ClusterClassName:    "ethereumEphemeralValidatorCluster",
		SkeletonBaseOptions: []string{"gethHercules", "lighthouseHercules", "lighthouseHerculesValidatorClient"},
		CloudCtxNs:          beacon_cookbooks.BeaconCloudCtxNs,
	}

	cl, err := read_topology.SelectClusterTopology(ctx, t.Tc.ProductionLocalTemporalOrgID, cd.ClusterClassName, cd.SkeletonBaseOptions)
	t.Require().Nil(err)
	t.Assert().NotEmpty(cl)
}

func (t *TopologyDeployActionRequestTestSuite) TestReadAppByID() {
	t.InitLocalConfigs()
	ctx := context.Background()
	//apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	cl, err := read_topology.SelectAppTopologyByID(ctx, t.Tc.ProductionLocalTemporalOrgID, 1677816792038031000)
	t.Require().Nil(err)
	t.Assert().NotEmpty(cl)
}
func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
