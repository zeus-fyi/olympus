package test_suites

import (
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type DatastoresTestSuite struct {
	test_suites_base.TestSuite

	Redis     *redis.Client
	PG        apps.Db
	PGTest    PGTestSuite
	RedisTest RedisTestSuite
}

func (d *DatastoresTestSuite) SetupTest() {
	d.InitLocalConfigs()
	d.RedisTest.SetupRedisConn()
	d.Redis = d.RedisTest.Redis
	d.PGTest.SetupPGConn()
	d.PG = d.PGTest.Pg
}

func (d *DatastoresTestSuite) Cleanup() {
	//d.CleanupDb()
	//d.CleanCache()
}
func TestDatastoresTestSuite(t *testing.T) {
	suite.Run(t, new(DatastoresTestSuite))
}
