package config_maps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
	zeus_config_map_reqs "github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_req_types/config_maps"
)

type ConfigMapsActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *ConfigMapsActionRequestTestSuite) TestCmKeySwap() {
	t.InitLocalConfigs()
	t.Eg.POST("/configmaps", ConfigMapActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	cmr := zeus_config_map_reqs.ConfigMapActionRequest{
		Action:        zeus_config_map_reqs.KeySwapAction,
		ConfigMapName: "cm-lighthouse",
		Keys: zeus_config_map_reqs.KeySwap{
			KeyOne: "start.sh",
			KeyTwo: "pause.sh",
		},
		FilterOpts: nil,
	}

	r, err := t.ZeusClient.SwapConfigMapKeys(ctx, cmr)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

func (t *ConfigMapsActionRequestTestSuite) TestSetOrCreateKeyFromExisting() {
	t.InitLocalConfigs()
	t.Eg.POST("/configmaps", ConfigMapActionRequestHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)
	cmr := zeus_config_map_reqs.ConfigMapActionRequest{
		Action:        zeus_config_map_reqs.KeySwapAction,
		ConfigMapName: "cm-lighthouse",
		Keys: zeus_config_map_reqs.KeySwap{
			KeyOne: "pause.sh",
			KeyTwo: "start.sh",
		},
		FilterOpts: nil,
	}

	r, err := t.ZeusClient.SetOrCreateKeyFromConfigMapKey(ctx, cmr)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

func TestConfigMapsActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapsActionRequestTestSuite))
}
