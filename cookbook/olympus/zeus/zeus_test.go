package zeus_cookbook

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbook"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type ZeusCookbookTestSuite struct {
	base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *ZeusCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbook.ChangeToCookbookDir()
}

func TestZeusCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusCookbookTestSuite))
}
