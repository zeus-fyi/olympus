package create

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

type ChartPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartPackagesTestSuite) TestConvertDeploymentAndInsert() {
	ns := sql.NullString{}
	c := Chart{autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err := c.InsertChart(ctx, q)
	s.Require().Nil(err)
	s.Assert().NotZero(c.ChartPackageID)
}

func (s *ChartPackagesTestSuite) TestInsert() {
	ns := sql.NullString{}
	c := Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	q := sql_query_templates.NewQueryParam("InsertMockChartForTest", "table", "where", 1000, []string{})
	ctx := context.Background()
	err := c.InsertChart(ctx, q)
	s.Require().Nil(err)

}
func TestChartPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartPackagesTestSuite))
}
