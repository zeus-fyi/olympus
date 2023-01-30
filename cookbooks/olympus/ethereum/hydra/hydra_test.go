package olympus_hydra_cookbooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_client "github.com/zeus-fyi/olympus/pkg/zeus/client"
	api_configs "github.com/zeus-fyi/olympus/test/configs"
)

type HydraCookbookTestSuite struct {
	test_suites_base.TestSuite
	ZeusTestClient zeus_client.ZeusClient
}

func (t *HydraCookbookTestSuite) TestClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	cd := HydraClusterConfig("ephemery")
	gcd := cd.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := cd.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := cd.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *HydraCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()

	tc := api_configs.InitLocalTestConfigs()
	t.ZeusTestClient = zeus_client.NewDefaultZeusClient(tc.Bearer)
}

func TestHydraCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(HydraCookbookTestSuite))
}
