package zeus_ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
)

func (t *AresZeusEthereumTestSuite) TestDeployConsensusClient() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.Deploy(ctx, deployConsensusClientKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
