package ethereum_beacon_cookbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

var ctx = context.Background()

var topCloudCtxNs = zeus_common_types.CloudCtxNs{
	CloudProvider: "do",
	Region:        "sfo3",
	Context:       "do-sfo3-dev-do-sfo3-zeus",
	Namespace:     "beacon", // set with your own namespace
	Env:           "dev",
}

type ZeusEthereumBeaconTestSuite struct {
	base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *ZeusEthereumBeaconTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	// points working dir to inside /test
	cookbooks.ChangeToCookbookDir()

	// generates outputs to /test/outputs dir
	// uses inputs from /test/mocks dir
}

func TestZeusEthereumBeaconTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusEthereumBeaconTestSuite))
}

func CreateBeacon() {

	cookbooks.ChangeToCookbookDir()

}
