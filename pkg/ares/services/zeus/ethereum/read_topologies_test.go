package zeus_ethereum

import "github.com/zeus-fyi/olympus/pkg/ares/ethereum"

func (t *AresZeusEthereumTestSuite) TestReadTopologies() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.ReadTopologies(ctx)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
