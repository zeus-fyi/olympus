package test_suites

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatastoresTestSuite struct {
	PGTestSuite
	RedisTestSuite
}

func (d *DatastoresTestSuite) SetupTest() {
	d.SetupRedisConn()
	d.SetupPGConn()
}

func (d *DatastoresTestSuite) Cleanup() {
	//d.CleanupDb()
	//d.CleanCache()
}
func TestDatastoresTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
