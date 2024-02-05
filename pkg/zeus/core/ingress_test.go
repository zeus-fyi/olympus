package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type IngressTestSuite struct {
	K8TestSuite
}

func (ing *IngressTestSuite) TestGetIngress() {
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "data", Namespace: "eth-indexer"}

	pods, err := ing.K.GetPodsUsingCtxNs(ctx, kns, nil, nil)
	ing.Require().Nil(err)
	ing.Require().NotEmpty(pods)
}

func TestIngressTestSuite(t *testing.T) {
	suite.Run(t, new(IngressTestSuite))
}
