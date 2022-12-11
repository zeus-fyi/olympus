package beacon_cookbooks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type BeaconCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *BeaconCookbookTestSuite) TestUploadBeaconCharts() {
	ctx := context.Background()
	// Consensus
	resp, err := t.ZeusTestClient.UploadChart(ctx, BeaconConsensusClientChartPath, ConsensusClientChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	DeployConsensusClientKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: DeployConsensusClientKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(BeaconConsensusClientChartPath)
	t.Require().Nil(err)

	// Exec
	resp, err = t.ZeusTestClient.UploadChart(ctx, BeaconExecClientChartPath, ExecClientChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	DeployExecClientKnsReq.TopologyID = resp.TopologyID
	tar = zeus_req_types.TopologyRequest{TopologyID: DeployExecClientKnsReq.TopologyID}
	chartResp, err = t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(BeaconExecClientChartPath)
	t.Require().Nil(err)

}

func (t *BeaconCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbooks.ChangeToCookbookDir()
}

func TestBeaconCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconCookbookTestSuite))
}
