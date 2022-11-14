package zeus_ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
)

func (t *AresZeusEthereumTestSuite) TestDestroyDeployConsensusClient() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyConsensusClientKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}

func (t *AresZeusEthereumTestSuite) TestDestroyDeployExecClient() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyExecClientKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
