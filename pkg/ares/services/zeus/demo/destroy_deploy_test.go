package zeus_demo

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
)

func (t *AresDemoTestSuite) TestDestroyDemoObol() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, obolDestroyDeployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
