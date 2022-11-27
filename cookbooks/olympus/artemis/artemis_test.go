package artemis_cookbook

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type ArtemisCookbookTestSuite struct {
	base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *ArtemisCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbooks.ChangeToCookbookDir()
}

func (t *ArtemisCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.Deploy(ctx, ArtemisDeployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *ArtemisCookbookTestSuite) TestUploadCharts() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, artemisChartPath, uploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)
	ArtemisDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: ArtemisDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)
	err = chartResp.PrintWorkload(artemisChartPath)
	t.Require().Nil(err)
}

func TestArtemisCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisCookbookTestSuite))
}
