package test_suites

import (
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/datastores/postgres"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type DatastoresTestSuite struct {
	base.TestSuite

	Redis     *redis.Client
	PG        postgres.Db
	PGTest    PGTestSuite
	RedisTest RedisTestSuite
}

func (d *DatastoresTestSuite) SetupTest() {
	d.InitConfigs()
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
