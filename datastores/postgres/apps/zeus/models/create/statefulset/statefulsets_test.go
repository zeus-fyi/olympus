package statefulset

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	create_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type StatefulSetTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *StatefulSetTestSuite) TestConvertStatefulSetAndInsert() {
	filepath := s.TestDirectory + "/mocks/test/statefulset.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var sts *v1.StatefulSet
	err = json.Unmarshal(jsonBytes, &sts)

	s.Require().Nil(err)
	s.Require().NotEmpty(sts)

	dbStatefulSetConfig, err := ConvertStatefulSetSpecConfigToDB(sts)
	s.Require().Nil(err)
	s.Require().NotEmpty(dbStatefulSetConfig)

	mockC := s.mockChart()
	s.Require().Nil(err)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertStatefulSet", "table", "where", 1000, []string{})
	stsInsert := StatefulSet{dbStatefulSetConfig}
	err = stsInsert.InsertStatefulSet(ctx, q, &mockC)
	s.Require().Nil(err)
}

func (s *StatefulSetTestSuite) mockChart() charts.Chart {
	ns := sql.NullString{}
	c := create_charts.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err := c.InsertChart(ctx, q)
	s.Require().Nil(err)

	mockC := charts.Chart{}
	mockC.ChartPackageID = c.GetChartPackageID()
	s.Require().Nil(err)
	return mockC
}

func (s *StatefulSetTestSuite) TestSeedChartComponents() {
	// only used to bootstrap for the main test
	s.SkipTest(true)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartComponentResources", "table", "where", 1000, []string{})

	cr := seedStatefulSet()
	err := s.InsertChartResource(ctx, q, cr)
	s.Require().Nil(err)
}

func (s *StatefulSetTestSuite) insertChartResource(c autogen_bases.ChartComponentResources) string {
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO chart_component_resources(chart_component_resource_id, chart_component_kind_name, chart_component_api_version)
 				 VALUES ('%d', '%s', '%s')`,
		c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion)
	return sqlInsertStatement
}

func (s *StatefulSetTestSuite) InsertChartResource(ctx context.Context, q sql_query_templates.QueryParams, c autogen_bases.ChartComponentResources) error {
	query := s.insertChartResource(c)
	_, err := apps.Pg.Exec(ctx, query)
	return err
}

func seedStatefulSet() autogen_bases.ChartComponentResources {
	cr := autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 0,
		ChartComponentKindName:   "StatefulSet",
		ChartComponentApiVersion: "apps/v1",
	}
	return cr
}

func TestStatefulsetTestSuite(t *testing.T) {
	suite.Run(t, new(StatefulSetTestSuite))
}
