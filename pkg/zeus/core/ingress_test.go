package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type IngressTestSuite struct {
	K8TestSuite
}

func (ing *IngressTestSuite) TestGetIngress() {
	ctx := context.Background()
	var kns = KubeCtxNs{Env: "", CloudProvider: "", Region: "", CtxType: "data", Namespace: "eth-indexer"}

	pods, err := ing.K.GetPodsUsingCtxNs(ctx, kns, nil, nil)
	ing.Require().Nil(err)
	ing.Require().NotEmpty(pods)
}

func TestIngressTestSuite(t *testing.T) {
	suite.Run(t, new(IngressTestSuite))
}
