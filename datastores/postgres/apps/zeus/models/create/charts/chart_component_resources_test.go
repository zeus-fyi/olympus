package create_charts

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type ChartComponentResourcesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ChartComponentResourcesTestSuite) TestSeedChartComponents() {
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartComponentResources", "table", "where", 1000, []string{})

	cr := seedStatefulSet()
	err := cr.InsertChartResource(ctx, q)
	s.Require().Nil(err)

	cr = seedService()
	err = cr.InsertChartResource(ctx, q)
	s.Require().Nil(err)
}

func seedStatefulSet() ChartComponentResources {
	cr := ChartComponentResources{autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 1,
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}}
	return cr
}
func seedService() ChartComponentResources {
	cr := ChartComponentResources{autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 2,
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "v1",
	}}
	return cr
}

func TestChartComponentResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartComponentResourcesTestSuite))
}
