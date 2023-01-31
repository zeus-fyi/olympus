package read_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	read_topology "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/read/topologies/topology"
	"github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/test"
)

type TopologyReadActionRequestTestSuite struct {
	test.TopologyActionRequestTestSuite
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	ctx := context.Background()

	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = 1675132721360688000
	// from auth lookup
	tr.OrgID = t.Tc.ProductionLocalTemporalOrgID
	tr.UserID = t.Tc.ProductionLocalTemporalUserID
	err := tr.SelectTopology(ctx)
	t.Require().Nil(err)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
