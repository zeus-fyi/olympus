package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DeploymentsTestSuite struct {
	K8TestSuite
}

func (ing *DeploymentsTestSuite) TestGetDeployment() {
	ctx := context.Background()
	var kns = KubeCtxNs{Env: "", CloudProvider: "", Region: "", CtxType: "data", Namespace: "eth-indexer"}
	pods, err := ing.K.GetDeployment(ctx, kns, "eth-indexer")
	ing.Require().Nil(err)
	ing.Require().NotEmpty(pods)
}

func TestDeploymentsTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentsTestSuite))
}
