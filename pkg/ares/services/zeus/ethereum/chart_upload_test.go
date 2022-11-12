package zeus_ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
	create_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create"
)

func (t *AresZeusEthereumTestSuite) TestCreateAndUploadConsensusClientChart() create_infra.TopologyCreateResponse {
	ethereum.ChangeDirToAresEthereumDir()
	p := ethereum.ConsensusClientPath()
	chartInfo := ethereum.ConsensusClientChartUploadRequest()
	resp, err := t.ZeusTestClient.UploadChart(ctx, p, chartInfo)
	t.Require().Nil(err)
	t.Assert().NotZero(resp.ID)
	return resp
}
