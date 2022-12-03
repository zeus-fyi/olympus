package env

import (
	"github.com/go-redis/redis/v9"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type ProductionPrototypeTest struct {
	test_suites_base.TestSuite

	Redis     *redis.Client
	RedisTest test_suites.RedisTestSuite
}
