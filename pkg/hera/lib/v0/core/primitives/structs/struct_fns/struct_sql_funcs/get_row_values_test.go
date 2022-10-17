package struct_sql_funcs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/hera/lib/v0/test"
)

type StructSQLTestSuite struct {
	test.AutoGenBaseTestSuiteBase
}

func (s *StructSQLTestSuite) TestGetRowValues() {

}

func TestStructSQLTestSuite(t *testing.T) {
	suite.Run(t, new(StructSQLTestSuite))
}
