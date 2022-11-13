package demo

import "github.com/zeus-fyi/olympus/pkg/ares/ethereum"

func (t *AresDemoTestSuite) TestDestroyDeploy() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
