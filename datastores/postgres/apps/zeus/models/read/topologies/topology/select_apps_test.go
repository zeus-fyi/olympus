package read_topology

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type SelectOrgAppsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

var ctx = context.Background()

func (s *TopologyTestSuite) TestSelectAppsByOrg() {
	s.InitLocalConfigs()
	//apps.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)
	orgID := s.Tc.ProductionLocalTemporalOrgID

	apps, err := SelectOrgApps(ctx, orgID)
	s.Require().Nil(err)
	s.Require().NotNil(apps)
}

func TestSelectOrgAppsTestSuite(t *testing.T) {
	suite.Run(t, new(SelectOrgAppsTestSuite))
}
