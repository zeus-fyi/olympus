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
	err := s.GetTableData()
	s.Require().Nil(err)
	s.Assert().NotEmpty(s.TableContent)
	s.Assert().NotEmpty(s.TableMap)
}

func TestPgSchemaAutogenTestSuite(t *testing.T) {
	suite.Run(t, new(PgSchemaAutogenTestSuite))
}
