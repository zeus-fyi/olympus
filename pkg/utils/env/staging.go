package env

import (
	"github.com/go-redis/redis/v9"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type StagingPrototypeTest struct {
	base.TestSuite

	Redis     *redis.Client
	PG        postgres.Db
	PGTest    test_suites.PGTestSuite
	RedisTest test_suites.RedisTestSuite
}

func (d *StagingPrototypeTest) SetupTest() {
	d.InitStagingConfigs()
	d.RedisTest.SetupRedisConn()
	d.Redis = d.RedisTest.Redis
	d.PGTest.SetupPGConn()
	d.PG = d.PGTest.Pg
}
