package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/databases/postgres"
)

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) SetupTest() {
}

type General struct {
	Index  string
	Pubkey int
}

func (g *General) GetRowValues() postgres.RowValues {
	return []string{g.Index, fmt.Sprintf("%d", g.Pubkey)}
}

type Wrapper struct {
	Generals []General
}

func (w *Wrapper) GetManyRowValues() postgres.RowEntries {
	var pgRows postgres.RowEntries
	for _, gen := range w.Generals {
		pgRows.Rows = append(pgRows.Rows, gen.GetRowValues())
	}
	return pgRows
}

func (s *UtilTestSuite) TestStrBuilder() {

	genSlice := makeGeneralSlice(2)

	sql := "INSERT INTO table (id, column) VALUES "

	rowValues := genSlice.GetManyRowValues()
	query := SQLDelimitedSliceStrBuilder(sql, rowValues)

	sqlExpected := "INSERT INTO table (id, column) VALUES (0,1),(1,2)"

	s.Assert().Equal(sqlExpected, query)
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
