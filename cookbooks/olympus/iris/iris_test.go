package iris_olympus_cookbook

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
	zeus_client "github.com/zeus-fyi/zeus/pkg/zeus/client"
)

type IrisCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

var ctx = context.Background()

func (t *IrisCookbookTestSuite) SetupTest() {
	// points dir to test/configs
	tc := api_configs.InitLocalTestConfigs()

	// uses the bearer token from test/configs/config.yaml
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
	//t.ZeusTestClient = zeus_client.NewLocalZeusClient(tc.Bearer)
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestIrisCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(IrisCookbookTestSuite))
}
