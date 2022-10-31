package create_topology

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
	ctx := context.Background()
	ou := org_users.OrgUser{}
	ou.OrgID, ou.UserID = t.h.NewTestOrgAndUser()
	top := NewCreateOrgUsersInfraTopology(ou)
	top.Name = fmt.Sprintf("top_name_%d", top.TopologyID)
	err := top.InsertOrgUsersTopology(ctx)
	t.Require().Nil(err)
	t.Require().NotZero(top.TopologyID)
}

func TestTopologiesTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesTestSuite))
}
