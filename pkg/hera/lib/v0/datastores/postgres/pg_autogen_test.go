package postgres

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites"
)

type PgSchemaAutogenTestSuite struct {
	test_suites.PGTestSuite
	PgSchemaAutogen
}

func (s *PgSchemaAutogenTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.PgSchemaAutogen = NewPgSchemaAutogen(s.Tc.LocalDbPgconn)
}
func (s *PgSchemaAutogenTestSuite) TestTablesSchemaRead() {
	tables, err := s.GetTables()
	s.Require().Nil(err)
	s.Assert().NotEmpty(tables)
}

func TestPgSchemaAutogenTestSuite(t *testing.T) {
	suite.Run(t, new(PgSchemaAutogenTestSuite))
}
