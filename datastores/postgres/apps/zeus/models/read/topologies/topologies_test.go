package read_topologies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologiesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

var ctx = context.Background()

func (t *TopologiesTestSuite) TestSelectTopologiesMetadata() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	orgID := 7138983863666903883
	tps, err := SelectTopologiesMetadata(ctx, orgID)
	t.Require().Nil(err)
	t.Assert().NotEmpty(tps)
}

func (t *TopologiesTestSuite) TestRead() {
	dr := NewReadTopologiesMetadataGroup()
	orgID := 1667452524363177528
	userID := 1667452524356256466
	ou := org_users.NewOrgUserWithID(orgID, userID)
	err := dr.SelectTopologiesMetadata(ctx, ou)
	t.Require().Nil(err)
	t.Assert().NotEmpty(dr.Slice)
}

func TestTopologiesTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesTestSuite))
}
