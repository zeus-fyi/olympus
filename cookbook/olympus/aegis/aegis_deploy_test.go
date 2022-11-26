package aegis_olympus_cookbook

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

func (t *AegisCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	t.TestAegisSecretsCopy()
}

func (t *AegisCookbookTestSuite) TestChartUpload() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, aegisChartPath, uploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	deployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: deployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(aegisChartPath)
	t.Require().Nil(err)
}
