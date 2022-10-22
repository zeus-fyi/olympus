package hera

import (
	"testing"

	"github.com/fraenky8/tables-to-go/pkg/database"
	"github.com/fraenky8/tables-to-go/pkg/settings"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type HeraTestSuite struct {
	test_suites.PGTestSuite
}

func (s *HeraTestSuite) TestTableSchemaRead() {
	s.InitLocalConfigs()
	pgSettings := settings.New()
	pgSettings.User = "postgres"
	pgSettings.Pswd = "postgres"
	pgSettings.Host = "localhost"
	pgSettings.Port = "5432"
	pgSettings.DbName = "postgres"

	pg := database.NewPostgresql(pgSettings)
	err := pg.Connect()
	s.Require().Nil(err)
	tables, err := pg.GetTables()
	s.Require().Nil(err)

	s.Assert().NotEmpty(tables)
}

func TestHeraTestSuite(t *testing.T) {
	suite.Run(t, new(HeraTestSuite))
}
