package zeus_ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
)

func (t *AresZeusEthereumTestSuite) TestDestroyDeploy() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
