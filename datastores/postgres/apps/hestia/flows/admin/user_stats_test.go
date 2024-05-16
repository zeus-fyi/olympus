package flows_admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type UserStatsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *UserStatsTestSuite) TestInsertUser() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	us, err := SelectUserFlowStats(ctx)
	s.Require().Nil(err)
	s.Require().NotNil(us)
}

func TestUserStatsTestSuite(t *testing.T) {
	suite.Run(t, new(UserStatsTestSuite))
}
