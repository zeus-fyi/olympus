package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type ConfigMapTestSuite struct {
	K8TestSuite
}

func (c *ConfigMapTestSuite) TestGetConfigMap() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "do", Region: "sfo3", Context: "do-nyc1-do-nyc1-zeus-demo", Namespace: "ephemeral"}
	cm, err := c.K.GetConfigMapWithKns(ctx, kns, "cm-choreography", nil)
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

func (c *ConfigMapTestSuite) TestSwitchCmKeys() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "", Namespace: "ethereum"}
	cm, err := c.K.ConfigMapKeySwap(ctx, kns, "cm-lighthouse", "start.sh", "pause.sh", nil)
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

func (c *ConfigMapTestSuite) TestConfigMapOverwriteOrCreateFromKey() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "", Namespace: "ethereum"}
	cm, err := c.K.ConfigMapOverwriteOrCreateFromKey(ctx, kns, "cm-lighthouse", "start.sh", "pause.sh", nil)
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

const lighthouseStart = `#!/bin/sh
    exec lighthouse beacon_node \
              --log-format=JSON \
              --datadir=/data \
              --enr-tcp-port=9000 \
              --enr-udp-port=9000 \
              --listen-address=0.0.0.0 \
              --port=9000 \
              --discovery-port=9000 \
              --http \
              --http-address=0.0.0.0 \
              --http-port=5052 \
              --execution-jwt=/data/jwt.hex \
              --execution-endpoint="http://zeus-geth:8551"`

const pause = `
	#!/bin/sh
    exec sleep 100000000000000000`

func (c *ConfigMapTestSuite) TestConfigMapCreateNewKeyOrOverwrite() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{Env: "", CloudProvider: "", Region: "", Context: "", Namespace: "ethereum"}

	m := make(map[string]string)
	m["lighthouse.sh"] = lighthouseStart

	cm, err := c.K.ConfigMapOverwriteOrCreateNewKeys(ctx, kns, "cm-lighthouse", m, nil)
	c.Require().Nil(err)
	c.Require().NotEmpty(cm)
}

func TestConfigMapTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapTestSuite))
}
