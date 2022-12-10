package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
)

type ConfigMapTestSuite struct {
	K8TestSuite
}

func (c *ConfigMapTestSuite) TestGetConfigMap() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "", Namespace: "ethereum"}
	cm, err := c.K.GetConfigMapWithKns(ctx, kns, "cm-lighthouse", nil)
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

func (c *ConfigMapTestSuite) TestSwitchCmKeys() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "", Namespace: "ethereum"}
	cm, err := c.K.ConfigMapKeySwap(ctx, kns, nil, "cm-lighthouse", "start.sh", "pause.sh")
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}
