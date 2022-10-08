package postgres_apps

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/admin"
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
func (s *PostgresTestSuite) TestDumpValidatorBalancesAtEpochTable() {
	ctx := context.Background()
	var PgTestDB Db
	conn := PgTestDB.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.Assert().NotNil(conn)
	defer conn.Close()

	le, he := 134000, 135000
	_, _, err := admin.DumpValidatorBalancesAtEpochTable(ctx, le, he)
	s.Require().Nil(err)
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
