package string_utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type General struct {
	Index  string
	Pubkey int
}

type Wrapper struct {
	Generals []General
}

func (g *General) GetRowValues() apps.RowValues {
	pgValues := apps.RowValues{g.Index, fmt.Sprintf("%d", g.Pubkey)}
	return pgValues
}

func (w *Wrapper) GetManyRowValues() apps.RowEntries {
	var pgRows apps.RowEntries
	for _, gen := range w.Generals {
		pgRows.Rows = append(pgRows.Rows, gen.GetRowValues())
	}
	return pgRows
}

func (w *Wrapper) GetManyRowValuesFlattened() apps.RowValues {
	var pgRows apps.RowValues
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

func (s *UtilTestSuite) TestDelimitedStrBuilderSQL() {
	genSlice := makeGeneralSlice(2)
	sql := "INSERT INTO table (id, column) VALUES "
	rowValues := genSlice.GetManyRowValues()
	query := DelimitedSliceStrBuilderSQLRows(sql, rowValues)
	sqlExpected := "INSERT INTO table (id, column) VALUES ('0','1'),('1','2')"
	s.Assert().Equal(sqlExpected, query)
}

func (s *UtilTestSuite) TestArrayListStrBuilderSQL() {
	genSlice := makeGeneralSlice(2)
	rowValues := genSlice.GetManyRowValuesFlattened()
	query := AnyArraySliceStrBuilderSQL(rowValues)
	sqlStrExpected := "ANY(ARRAY['0','1','1','2'])"
	s.Assert().Equal(sqlStrExpected, query)

	onlyArrayQuery := ArraySliceStrBuilderSQL(rowValues)
	sqlArrayStrExpected := "ARRAY['0','1','1','2']"
	s.Assert().Equal(sqlArrayStrExpected, onlyArrayQuery)
}

func (s *UtilTestSuite) TestMultiArraySliceStrBuilderSQL() {
	genSlice := makeGeneralSlice(2)
	rowValues := genSlice.GetManyRowValues()
	query := MultiArraySliceStrBuilderSQL(rowValues)
	sqlStrExpected := "ARRAY['0','1'],ARRAY['1','2']"
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
