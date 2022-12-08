package read_topologies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type TopologiesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (t *TopologiesTestSuite) TestRead() {
	dr := NewReadTopologiesMetadataGroup()
	ctx := context.Background()
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
