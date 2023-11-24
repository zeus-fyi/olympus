package read_topology

import (
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologyTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *TopologyTestSuite) TestSelectTopology() {
	s.InitLocalConfigs()
	//apps.Pg.InitPG(context.Background(), s.Tc.ProdLocalDbPgconn)
	tr := NewInfraTopologyReader()

	tr.TopologyID = 1700785711692338585
	tr.OrgID = s.Tc.ProductionLocalTemporalOrgID
	tr.UserID = s.Tc.ProductionLocalTemporalUserID
	err := tr.SelectTopologyForOrg(ctx)
	s.Require().Nil(err)
}

func TestTopologyTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyTestSuite))
}
