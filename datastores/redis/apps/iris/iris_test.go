package iris_redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type IrisRedisTestSuite struct {
	test_suites_base.TestSuite
}

func (r *IrisRedisTestSuite) SetupTest() {
	ctx := context.Background()
	r.InitLocalConfigs()
	apps.Pg.InitPG(ctx, r.Tc.ProdLocalDbPgconn)
	//InitLocalTestProductionRedisIrisCache(ctx)

	InitLocalTestRedisIrisCache(ctx)
}

func TestIrisRedisTestSuite(t *testing.T) {
	suite.Run(t, new(IrisRedisTestSuite))
}
