package beacon_cookbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbook"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type BeaconCookbookTestSuite struct {
	base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *BeaconCookbookTestSuite) TestUploadBeaconCharts() {
	ctx := context.Background()
	// Consensus
	resp, err := t.ZeusTestClient.UploadChart(ctx, beaconConsensusClientChartPath, consensusClientChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	deployConsensusClientKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: deployConsensusClientKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(beaconConsensusClientChartPath)
	t.Require().Nil(err)

	// Exec
	resp, err = t.ZeusTestClient.UploadChart(ctx, beaconExecClientChartPath, execClientChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	deployExecClientKnsReq.TopologyID = resp.TopologyID
	tar = zeus_req_types.TopologyRequest{TopologyID: deployExecClientKnsReq.TopologyID}
	chartResp, err = t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(beaconExecClientChartPath)
	t.Require().Nil(err)

}

func (t *BeaconCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbook.ChangeToCookbookDir()
}

func TestBeaconCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconCookbookTestSuite))
}
