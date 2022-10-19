package create

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

type ChartPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartPackagesTestSuite) TestConvertDeploymentAndInsert() {
	ns := sql.NullString{}
	c := Chart{autogen_structs.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err := c.InsertChart(ctx, q, c)
	s.Require().Nil(err)
	s.Assert().NotZero(c.ChartPackageID)
}

func TestChartPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartPackagesTestSuite))
}
