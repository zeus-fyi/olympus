package hestia_nodes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
)

var ctx = context.Background()

type NodesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *NodesTestSuite) TestSelectNodes() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	nf := NodeFilter{
		CloudProvider: "do",
		Region:        "nyc1",
		ResourceSums: zeus_core.ResourceSums{
			MemRequests: "12Gi",
			CpuRequests: "6",
		},
	}
	nodes, err := SelectNodes(ctx, nf)
	s.Require().NoError(err)
	s.Assert().NotEmpty(nodes)
}

func (s *NodesTestSuite) TestInsertNodes() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	nodes := hestia_autogen_bases.NodesSlice{}
	nodes = append(nodes, hestia_autogen_bases.Nodes{
		Description:   "test",
		Slug:          "sg-sfo1-01",
		Disk:          25,
		PriceHourly:   1,
		CloudProvider: "do",
		Vcpus:         2,
		PriceMonthly:  1,
		Region:        "sfo1",
		Memory:        3,
	})
	err := InsertNodes(ctx, nodes)
	s.Require().NoError(err)
}

func TestCreateNodesTestSuite(t *testing.T) {
	suite.Run(t, new(NodesTestSuite))
}
