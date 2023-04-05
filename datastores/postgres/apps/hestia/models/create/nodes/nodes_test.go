package hestia_nodes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type CreateNodesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateNodesTestSuite) TestInsertOrg() {
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
	suite.Run(t, new(CreateNodesTestSuite))
}
