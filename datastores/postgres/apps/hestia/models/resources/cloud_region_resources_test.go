package hestia_compute_resources

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

type CloudResourcesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CloudResourcesTestSuite) TestSetup() {
	s.InitLocalConfigs()
}

func (s *CloudResourcesTestSuite) TestCloudResourcesQuery() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	nf := NodeFilter{
		ResourceSums: zeus_core.ResourceSums{
			MemRequests: "1Gi",
			CpuRequests: "1",
		},
		Ou: s.Ou,
	}
	nf.Ou.OrgID = 1699642242976434000
	cm, err := SelectNodesV2(ctx, nf)
	s.Require().NoError(err)
	s.Require().NotEmpty(cm)
}

func (s *CloudResourcesTestSuite) TestCloudResourcesStruct() {
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
