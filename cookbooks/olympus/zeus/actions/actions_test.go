package zeus_actions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type ZeusCookbookActionsTestSuite struct {
	test_suites_base.TestSuite
	ZeusActionsClient ZeusActionsClient
}

var ctx = context.Background()

func (t *ZeusCookbookActionsTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusActionsClient.ZeusClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestZeusCookbookActionsTestSuite(t *testing.T) {
	suite.Run(t, new(ZeusCookbookActionsTestSuite))
}
