package ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
	read_infra "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/read"
)

func (t *AresZeusTestSuite) readUploadedConsensusChart(topID int) {
	tar := read_infra.TopologyReadRequest{TopologyID: topID}
	resp, err := t.ZeusTestClient.ReadChart(ctx, tar)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)

	ethereum.ChangeDirToAresEthereumDir()
	p := ethereum.ConsensusReadChartThenWritePath()
	err = resp.PrintWorkload(p)
	t.Require().Nil(err)
}
