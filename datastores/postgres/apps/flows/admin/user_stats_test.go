package flows_admin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type UserStatsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *UserStatsTestSuite) TestInsertUser() {
	us, err := SelectUserFlowStats(ctx, s.Ou)
	s.Require().Nil(err)
	s.Require().NotNil(us)
}

func TestUserStatsTestSuite(t *testing.T) {
	suite.Run(t, new(UserStatsTestSuite))
}
