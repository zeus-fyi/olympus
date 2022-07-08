package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/test_utils/test_suites"
)

type PostgresTestSuite struct {
	test_suites.BaseTestSuite
}

func (s *PostgresTestSuite) TestConnPG() {
	conn, err := InitPG(context.Background(), s.Tc.TEST_DB_PGCONN)
	s.Require().Nil(err)
	s.Assert().NotNil(conn)
	defer conn.Close()
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
