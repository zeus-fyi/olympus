package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/databases/postgres"
)

type General struct {
	Index  string
	Pubkey int
}

type Wrapper struct {
	Generals []General
}

func (g *General) GetRowValues() postgres.RowValues {
	pgValues := postgres.RowValues{g.Index, fmt.Sprintf("%d", g.Pubkey)}
	return pgValues
}

func (w *Wrapper) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, gen := range w.Generals {
		pgRows.Rows = append(pgRows.Rows, gen.GetRowValues())
	}
	return pgRows
}

func (w *Wrapper) GetManyRowValuesFlattened() postgres.RowValues {
	var pgRows postgres.RowValues
	for _, gen := range w.Generals {
		pgRows = append(pgRows, gen.GetRowValues()...)
	}
	return pgRows
}

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) SetupTest() {
}

func (s *UtilTestSuite) TestSQLDelimitedStrBuilder() {
	genSlice := makeGeneralSlice(2)
	sql := "INSERT INTO table (id, column) VALUES "
	rowValues := genSlice.GetManyRowValues()
	query := DelimitedSliceStrBuilderSQLRows(sql, rowValues)
	sqlExpected := "INSERT INTO table (id, column) VALUES ('0','1'),('1','2')"
	s.Assert().Equal(sqlExpected, query)
}

func (s *UtilTestSuite) TestSQLArrayListStrBuilder() {
	genSlice := makeGeneralSlice(2)
	rowValues := genSlice.GetManyRowValuesFlattened()
	query := ArraySliceStrBuilderSQL(rowValues)
	sqlStrExpected := "ANY(ARRAY['0','1','1','2'])"
	s.Assert().Equal(sqlStrExpected, query)
}

func makeGeneralSlice(len int) Wrapper {
	var w Wrapper
	genSlice := make([]General, 2)
	for i := 0; i < len; i++ {
		genSlice[i] = General{
			Index:  fmt.Sprintf("%d", i),
			Pubkey: i + 1,
		}
	}
	w.Generals = genSlice
	return w
}
func TestUtilTestSuite(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}
