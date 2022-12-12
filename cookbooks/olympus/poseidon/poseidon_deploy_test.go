package poseidon_olympus_cookbook

import "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types"

func (t *PoseidonCookbookTestSuite) TestDeploy() {
	resp, err := t.ZeusTestClient.Deploy(ctx, PoseidonDeployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
	//t.TestPoseidonSecretsCopy()
}

func (t *PoseidonCookbookTestSuite) TestChartUpload() {
	resp, err := t.ZeusTestClient.UploadChart(ctx, PoseidonChartPath, PoseidonUploadChart)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.TopologyID)

	PoseidonDeployKnsReq.TopologyID = resp.TopologyID
	tar := zeus_req_types.TopologyRequest{TopologyID: PoseidonDeployKnsReq.TopologyID}
	chartResp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(chartResp)

	err = chartResp.PrintWorkload(PoseidonChartPath)
	t.Require().Nil(err)
}
