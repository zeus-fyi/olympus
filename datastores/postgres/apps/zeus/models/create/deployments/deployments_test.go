package deployments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type DeploymentsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *DeploymentsTestSuite) TestConvertDeploymentAndInsert() {
	filepath := s.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig, err := ConvertDeploymentConfigToDB(d)
	s.Require().Nil(err)
	s.Require().NotEmpty(dbDeploymentConfig)

	mockC := s.mockChart()
	s.Require().Nil(err)

	ctx := context.Background()

	q := sql_query_templates.NewQueryParam("InsertDeployment", "table", "where", 1000, []string{})

	deploymentInsert := Deployment{dbDeploymentConfig}
	err = deploymentInsert.InsertDeployment(ctx, q, &mockC)
	s.Require().Nil(err)
}

func (s *DeploymentsTestSuite) mockChart() create.Chart {
	ns := sql.NullString{}
	c := create.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	q := sql_query_templates.NewQueryParam("InsertDeploymentMockChart", "table", "where", 1000, []string{})
	ctx := context.Background()
	err := c.InsertChart(ctx, q)
	s.Require().Nil(err)

	return c
}

func (s *DeploymentsTestSuite) TestSeedChartComponents() {
	// only used to bootstrap for the main test
	s.SkipTest(true)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartComponentResources", "table", "where", 1000, []string{})

	cr := seedDeployment()
	err := s.InsertChartResource(ctx, q, cr)
	s.Require().Nil(err)
}

func (s *DeploymentsTestSuite) insertChartResource(c autogen_bases.ChartComponentResources) string {
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO chart_component_resources(chart_component_resource_id, chart_component_kind_name, chart_component_api_version)
 				 VALUES ('%d', '%s', '%s')`,
		c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion)
	return sqlInsertStatement
}

func (s *DeploymentsTestSuite) InsertChartResource(ctx context.Context, q sql_query_templates.QueryParams, c autogen_bases.ChartComponentResources) error {
	query := s.insertChartResource(c)
	_, err := apps.Pg.Exec(ctx, query)
	return err
}

func seedDeployment() autogen_bases.ChartComponentResources {
	cr := autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 0,
		ChartComponentKindName:   "Deployment",
		ChartComponentApiVersion: "apps/v1",
	}
	return cr
}

func TestDeploymentsTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentsTestSuite))
}
