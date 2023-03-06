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

var ctx = context.Background()

func (t *TopologyReadActionRequestTestSuite) TestReadTopologiesOrgCloudCtxNs() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	t.Eg.POST("/infra/read/org/topologies", ReadTopologiesOrgCloudCtxNsHandler)

	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9010")
	}()

	<-start
	defer t.E.Shutdown(ctx)
	resp, err := t.ZeusClient.ReadTopologiesOrgCloudCtxNs(ctx)
	t.Require().Nil(err)
	t.Require().NotEmpty(resp)
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = 1675200222894809000
	// from auth lookup
	tr.OrgID = t.Tc.ProductionLocalTemporalOrgID
	tr.UserID = t.Tc.ProductionLocalTemporalUserID
	err := tr.SelectTopology(ctx)
	t.Require().Nil(err)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
