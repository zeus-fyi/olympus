package hestia_compute_resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type OrgResourcesListTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *DisksTestSuite) TestSelectFreeTrialDoNodes() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	orgID := 1679515557647002001

	orSlice, err := SelectFreeTrialDigitalOceanNodes(ctx, orgID)
	s.Require().NoError(err)
	s.Require().NotEmpty(orSlice)
	fmt.Println(orSlice)
}

func (s *DisksTestSuite) TestOrgResourcesListNodesTestSuite() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	orgID := 1679515557647002001

	orSlice, err := SelectOrgResourcesNodes(ctx, orgID)
	s.Require().NoError(err)
	s.Require().NotEmpty(orSlice)
	fmt.Println(orSlice)
}

func (s *DisksTestSuite) TestOrgResourcesListDisksTestSuite() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	orgID := 1679515557647002001

	orSlice, err := SelectOrgResourcesDisks(ctx, orgID)
	s.Require().NoError(err)
	s.Require().NotEmpty(orSlice)
	fmt.Println(orSlice)
}

func TestOrgResourcesListTestSuite(t *testing.T) {
	suite.Run(t, new(OrgResourcesListTestSuite))
}
