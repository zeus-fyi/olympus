package ethereum

import (
	"github.com/zeus-fyi/olympus/pkg/ares/ethereum"
	destroy_deploy_request "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/deploy/external/destroy"
)

var deployDestroyKnsReq = destroy_deploy_request.TopologyDestroyDeployRequest{
	TopologyID:    1667958167340986000,
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "dev-sfo3-zeus",
	Namespace:     "demo",
	Env:           "dev",
}

func (t *AresZeusTestSuite) TestDestroyDeploy() {
	ethereum.ChangeDirToAresEthereumDir()
	resp, err := t.ZeusTestClient.DestroyDeploy(ctx, deployDestroyKnsReq)
	t.Require().Nil(err)
	t.Assert().NotEmpty(resp)
}
