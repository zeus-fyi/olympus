package ingresses

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type IngressTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *IngressTestSuite) TestK8sIngressYamlReaderAndK8sToDBCte() {
	ing := NewIngress()
	filepath := s.TestDirectory + "/mocks/test/ingress.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &ing.K8sIngress)

	s.Require().Nil(err)
	s.Require().NotEmpty(ing.K8sIngress)

	err = ing.ParseK8sConfigToDB()
	s.Require().Nil(err)

	s.Require().NotEmpty(ing.Metadata)
	s.Require().NotEmpty(ing.Spec)
	s.Require().NotEmpty(ing.TLS)
	s.Require().NotEmpty(ing.Rules)

	c := charts.Chart{}
	c.ChartPackageID = 100
	cte := ing.GetIngressSpecCTE(&c)
	s.Require().NotEmpty(cte)
}

func (s *IngressTestSuite) TestSeedChartComponents() {
	// only used to bootstrap for the main test
	s.SkipTest(true)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("ChartComponentResources", "table", "where", 1000, []string{})

	cr := seedService()
	err := s.InsertChartResource(ctx, q, cr)
	s.Require().Nil(err)
}

func (s *IngressTestSuite) insertChartResource(c autogen_bases.ChartComponentResources) string {
	sqlInsertStatement := fmt.Sprintf(
		`INSERT INTO chart_component_resources(chart_component_resource_id, chart_component_kind_name, chart_component_api_version)
 				 VALUES ('%d', '%s', '%s')`,
		c.ChartComponentResourceID, c.ChartComponentKindName, c.ChartComponentApiVersion)
	return sqlInsertStatement
}

func (s *IngressTestSuite) InsertChartResource(ctx context.Context, q sql_query_templates.QueryParams, c autogen_bases.ChartComponentResources) error {
	query := s.insertChartResource(c)
	_, err := apps.Pg.Exec(ctx, query)
	return err
}

func seedService() autogen_bases.ChartComponentResources {
	cr := autogen_bases.ChartComponentResources{
		ChartComponentResourceID: 14,
		ChartComponentKindName:   "Ingress",
		ChartComponentApiVersion: "networking.k8s.io/v1",
	}
	return cr
}

func TestIngressTestSuite(t *testing.T) {
	suite.Run(t, new(IngressTestSuite))
}
