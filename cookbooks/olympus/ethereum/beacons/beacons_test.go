package olympus_beacon_cookbooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	olympus_cookbooks "github.com/zeus-fyi/olympus/cookbooks"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type BeaconsTestCookbookTestSuite struct {
	test_suites_base.TestSuite
}

func (t *BeaconsTestCookbookTestSuite) TestClusterSetup() {
	olympus_cookbooks.ChangeToCookbookDir()

	gcd := EphemeralBeaconBaseClusterDefinition.BuildClusterDefinitions()
	t.Assert().NotEmpty(gcd)
	fmt.Println(gcd)

	gdr := EphemeralBeaconBaseClusterDefinition.GenerateDeploymentRequest()
	t.Assert().NotEmpty(gdr)
	fmt.Println(gdr)

	sbDefs, err := EphemeralBeaconBaseClusterDefinition.GenerateSkeletonBaseCharts()
	t.Require().Nil(err)
	t.Assert().NotEmpty(sbDefs)
}

func (t *BeaconsTestCookbookTestSuite) SetupTest() {
	olympus_cookbooks.ChangeToCookbookDir()
}

func TestBeaconsTestCookbookTestSuite(t *testing.T) {
	suite.Run(t, new(BeaconsTestCookbookTestSuite))
}
