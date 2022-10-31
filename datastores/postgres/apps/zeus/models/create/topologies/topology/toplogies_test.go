package create_topology

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologiesTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (t *TopologiesTestSuite) TestInsert() Topologies {
	top := NewCreateTopology()
	top.TopologyID = t.Ts.UnixTimeStampNow()
	top.Name = fmt.Sprintf("testTopology_%d", top.TopologyID)
	ctx := context.Background()
	err := top.InsertTopology(ctx)
	t.Require().Nil(err)
	t.Require().NotZero(top.TopologyID)
	return top
}

func (t *TopologiesTestSuite) TestInsertOrgUsersTopology() {
	top := t.TestInsert()
	top.OrgID, top.UserID = t.h.NewTestOrgAndUser()
	ctx := context.Background()

	err := top.InsertOrgUsersTopology(ctx)
	t.Require().Nil(err)

}

func TestTopologiesTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesTestSuite))
}
