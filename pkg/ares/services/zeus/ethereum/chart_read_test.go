package ethereum

import "github.com/zeus-fyi/olympus/pkg/ares/ethereum"

func (t *AresZeusTestSuite) TestReadUploadedConsensusChart() {
	ethereum.ChangeDirToAresEthereumDir()
	p := ethereum.ConsensusClientPath()
	chartInfo := ethereum.ConsensusClientChartUploadRequest()
	resp, err := t.ProdZeusClient.UploadChart(ctx, p, chartInfo)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.ID)
}
