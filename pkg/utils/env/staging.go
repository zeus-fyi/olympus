package env

import (
	"github.com/go-redis/redis/v9"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type StagingPrototypeTest struct {
	test_suites_base.TestSuite

	Redis     *redis.Client
	PG        apps.Db
	PGTest    test_suites.PGTestSuite
	RedisTest test_suites.RedisTestSuite
}

func (d *StagingPrototypeTest) SetupTest() {
	d.InitStagingConfigs()
	d.RedisTest.SetupRedisConnStaging()

	d.Redis = d.RedisTest.Redis
	d.PGTest.SetupStagingPGConn()
	d.PG = d.PGTest.Pg
}
