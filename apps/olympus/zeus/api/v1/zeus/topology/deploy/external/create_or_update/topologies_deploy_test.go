package create_or_update_deploy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacon"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyDeployActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyDeployActionRequestTestSuite) TestDeployChart() {
	t.InitLocalConfigs()
	t.Eg.POST("/deploy/cluster", ClusterTopologyDeploymentHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	cd := zeus_req_types.ClusterTopologyDeployRequest{
		ClusterName: "ethereum-beacons",
		BaseOptions: []string{"geth", "lighthouse"},
		CloudCtxNs:  beacon_cookbooks.BeaconCloudCtxNs,
	}
	resp, err := t.ZeusClient.DeployCluster(ctx, cd)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

}

func TestTopologyDeployActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyDeployActionRequestTestSuite))
}
