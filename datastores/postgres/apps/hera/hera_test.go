package hera

import (
	"testing"

	"github.com/fraenky8/tables-to-go/pkg/database"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type HeraTestSuite struct {
	test_suites.PGTestSuite
}

func (s *HeraTestSuite) TestTablesSchemaRead() {
	s.InitLocalConfigs()
	pgConf, err := PgxConfigToSqlX(s.Tc.LocalDbPgconn)
	s.Require().Nil(err)
	pg := database.NewPostgresql(pgConf)
	err = pg.Connect()
	s.Require().Nil(err)
	tables, err := pg.GetTables()
	s.Require().Nil(err)
	s.Assert().NotEmpty(tables)
}

func TestHeraTestSuite(t *testing.T) {
	suite.Run(t, new(HeraTestSuite))
}
