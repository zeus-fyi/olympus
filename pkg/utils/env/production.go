package env

import (
	"github.com/go-redis/redis/v9"
	"github.com/zeus-fyi/olympus/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type ProductionPrototypeTest struct {
	base.TestSuite

	Redis     *redis.Client
	PG        postgres.Db
	PGTest    test_suites.PGTestSuite
	RedisTest test_suites.RedisTestSuite
}

func (d *ProductionPrototypeTest) SetupTest() {
	d.InitProductionConfig()
	d.RedisTest.SetupRedisConnProduction()
	d.Redis = d.RedisTest.Redis
	d.PGTest.SetupProductionPGConn()
	d.PG = d.PGTest.Pg
}
