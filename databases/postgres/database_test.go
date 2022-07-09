package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/internal/test_utils/test_suites"
)

type PostgresTestSuite struct {
	test_suites.BaseTestSuite
}

func (s *PostgresTestSuite) TestConnPG() {
	var PgTestDB Db
	conn := PgTestDB.InitPG(context.Background(), s.Tc.TEST_DB_PGCONN)
	s.Assert().NotNil(conn)
	defer conn.Close()
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
