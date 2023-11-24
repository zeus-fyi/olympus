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

//func (t *TopologyReadActionRequestTestSuite) TestReadTopologiesOrgCloudCtxNs() {
//	t.InitLocalConfigs()
//	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
//	t.Eg.POST("/infra/read/org/topologies", ReadTopologiesOrgCloudCtxNsHandler)
//
//	start := make(chan struct{}, 1)
//	go func() {
//		close(start)
//		_ = t.E.Start(":9010")
//	}()
//
//	<-start
//	defer t.E.Shutdown(ctx)
//	resp, err := t.ZeusClient.ReadTopologiesOrgCloudCtxNs(ctx)
//	t.Require().Nil(err)
//	t.Require().NotEmpty(resp)
//}

func (t *TopologyReadActionRequestTestSuite) TestReadClusterDefinition() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	cl, err := read_topology.SelectClusterTopology(ctx, t.Tc.ProductionLocalTemporalOrgID, "ethereumEphemeralBeacons", []string{"lighthouseHercules", "gethHercules"})
	t.Require().Nil(err)
	t.Require().NotNil(cl)

	AppID := 1670997020811171000
	selectedApp, err := read_topology.SelectAppTopologyByID(ctx, AppsOrgID, AppID)
	t.Require().Nil(err)
	t.Require().NotNil(selectedApp)
}

func (t *TopologyReadActionRequestTestSuite) TestReadChart() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	tr := read_topology.NewInfraTopologyReader()
	tr.TopologyID = 1700788779135945000
	// from auth lookup
	tr.OrgID = t.Tc.ProductionLocalTemporalOrgID
	tr.UserID = t.Tc.ProductionLocalTemporalUserID
	err := tr.SelectTopologyForOrg(ctx)
	t.Require().Nil(err)
}

func TestTopologyReadActionRequestTestSuite(t *testing.T) {
	suite.Run(t, new(TopologyReadActionRequestTestSuite))
}
