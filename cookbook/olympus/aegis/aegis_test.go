package aegis_olympus_cookbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbook"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type AegisCookbookTestSuite struct {
	base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

var ctx = context.Background()

func (t *AegisCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	cookbook.ChangeToCookbookDir()
}

func TestAegisCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(AegisCookbookTestSuite))
}
