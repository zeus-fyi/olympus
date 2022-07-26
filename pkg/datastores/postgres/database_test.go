package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type PostgresTestSuite struct {
	base.TestSuite
}

func (s *PostgresTestSuite) TestConnPG() {
	var PgTestDB Db
	conn := PgTestDB.InitPG(context.Background(), s.Tc.LocalDbPgconn)
	s.Assert().NotNil(conn)
	defer conn.Close()
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
