package olympus_hydra_cookbooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type HydraCookbookTestSuite struct {
	test_suites_base.TestSuite
}

func (t *HydraCookbookTestSuite) TestClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	gcd := EphemeralHydraClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := EphemeralHydraClusterDefinition.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := EphemeralHydraClusterDefinition.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *HydraCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestHydraCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(HydraCookbookTestSuite))
}
