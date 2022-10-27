package networking

import (
	"context"
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type NetworkingTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *NetworkingTestSuite) TestSeedChartComponents() {
	// only used to bootstrap for the main test
	s.SkipTest(true)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartComponentResources", "table", "where", 1000, []string{})

	cr := seedService()
	err := s.InsertChartResource(ctx, q, cr)
	s.Require().Nil(err)
}

func (s *NetworkingTestSuite) insertChartResource(c autogen_bases.ChartComponentResources) string {
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO chart_component_resources(chart_component_resource_id, chart_component_kind_name, chart_component_api_version)
 				 VALUES ('%d', '%s', '%s')`,
		c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion)
	return sqlInsertStatement
}

func (s *NetworkingTestSuite) InsertChartResource(ctx context.Context, q sql_query_templates.QueryParams, c autogen_bases.ChartComponentResources) error {
	query := s.insertChartResource(c)
	_, err := apps.Pg.Exec(ctx, query)
	return err
}

func seedService() autogen_bases.ChartComponentResources {
	cr := autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 0,
		ChartComponentKindName:   "Service",
		ChartComponentApiVersion: "apps/v1",
	}
	return cr
}
