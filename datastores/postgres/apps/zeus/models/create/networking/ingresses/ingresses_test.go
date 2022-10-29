package create_ingresses

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	create_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

type IngressTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *IngressTestSuite) TestConvertIngressAndInsert() {

	ing := NewCreateIngress()
	filepath := s.TestDirectory + "/mocks/test/ingress.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	err = json.Unmarshal(jsonBytes, &ing.K8sIngress)

	s.Require().Nil(err)
	s.Require().NotEmpty(ing.K8sIngress)
	ing.ConvertK8sIngressToDB()
	s.Require().NotEmpty(ing.Metadata)
	s.Require().NotEmpty(ing.Spec)
	s.Require().NotEmpty(ing.TLS)
	s.Require().NotEmpty(ing.Rules)

	ns := sql.NullString{}
	c := create_charts.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err = c.InsertChart(ctx, q)
	s.Require().Nil(err)

	mockC := charts.Chart{}
	mockC.ChartPackageID = c.GetChartPackageID()
	// TODO insert ingress to chart
	s.Require().Nil(err)
}

func TestIngressTestSuite(t *testing.T) {
	suite.Run(t, new(IngressTestSuite))
}
