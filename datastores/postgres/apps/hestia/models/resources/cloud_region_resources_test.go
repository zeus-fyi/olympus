package hestia_compute_resources

import (
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type CloudResourcesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CloudResourcesTestSuite) TestSetup() {
	s.InitLocalConfigs()
}

func (s *CloudResourcesTestSuite) TestCloudResourcesQuery() {
	cloudResources := CloudProviderRegionsResourcesMap{
		"aws": RegionResourcesMap{
			"us-east-1": Resources{
				Nodes: []hestia_autogen_bases.Nodes{
					{
						Region:        "us-east-1",
						CloudProvider: "aws",
					},
				},
			},
		},
	}
	s.Require().NotNil(cloudResources)
}
func TestCloudResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(CloudResourcesTestSuite))
}
