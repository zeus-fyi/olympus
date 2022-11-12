package zeus_ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
)

func (t *AresZeusEthereumTestSuite) TestDeploy() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
