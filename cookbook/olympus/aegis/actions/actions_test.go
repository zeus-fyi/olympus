package aegis_actions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbook"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type AegisCookbookActionsTestSuite struct {
	base.TestSuite
	AegisActionsClient AegisActionsClient
}

var ctx = context.Background()

func (t *AegisCookbookActionsTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.AegisActionsClient.ZeusClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbook.ChangeToCookbookDir()
}

func TestAegisCookbookActionsTestSuite(t *testing.T) {
	suite.Run(t, new(AegisCookbookActionsTestSuite))
}
