package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigMapTestSuite struct {
	K8TestSuite
}

func (ing *ConfigMapTestSuite) TestGetConfigMap() {
	ctx := context.Background()
	var kns = CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "data", Namespace: "eth-indexer"}
	pods, err := ing.K.GetConfigMapWithKns(ctx, kns, "cm-eth-indexer", nil)
	ing.Require().Nil(err)
	ing.Require().NotEmpty(pods)
}

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}
