package ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
	create_or_update_deploy "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/create_or_update"
)

var deployKnsReq = create_or_update_deploy.TopologyDeployRequest{
	TopologyID:    1667958167340986000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "demo",
	Env:           "dev",
}

func (t *AresZeusTestSuite) TestDeploy() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.Deploy(ctx, deployKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
