package iris_api_requests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_redis "github.com/zeus-fyi/olympus/datastores/redis/apps/iris"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type IrisActivitiesTestSuite struct {
	test_suites_base.TestSuite
}

func (s *IrisActivitiesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	iris_redis.InitLocalTestProductionRedisIrisCache(ctx)
}

func TestIrisActivitiesTestSuite(t *testing.T) {
	suite.Run(t, new(IrisActivitiesTestSuite))
}
