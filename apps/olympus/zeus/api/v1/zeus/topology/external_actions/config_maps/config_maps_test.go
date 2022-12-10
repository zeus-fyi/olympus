package config_maps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	beacon_cookbooks "github.com/zeus-fyi/olympus/cookbooks/ethereum/beacons"
	zeus_configmap_reqs "github.com/zeus-fyi/olympus/pkg/zeus/client/zeus_req_types/config_maps"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type ConfigMapsActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *ConfigMapsActionRequestTestSuite) TestUpload() {
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

	cmr := zeus_configmap_reqs.ConfigMapActionRequest{
		TopologyDeployRequest: beacon_cookbooks.DeployConsensusClientKnsReq,
		Action:                zeus_configmap_reqs.KeySwapAction,
		ConfigMapName:         "cm-lighthouse",
		Keys: zeus_configmap_reqs.KeySwap{
			KeyOne: "start.sh",
			KeyTwo: "pause.sh",
		},
		FilterOpts: nil,
	}

	r, err := t.ZeusClient.SwapConfigMapKeys(ctx, cmr)
	t.Require().Nil(err)
	t.Assert().NotEmpty(r)
}

func TestConfigMapsActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigMapsActionRequestTestSuite))
}
