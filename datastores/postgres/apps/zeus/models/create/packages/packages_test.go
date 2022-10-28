package packages

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	create_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"k8s.io/apimachinery/pkg/util/rand"
)

type PackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PackagesTestSuite) TestInsert() {
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
	p.Require().Nil(err)

	nd := deployments.NewDeployment()
	nsvc := networking.NewService()
	pkg := Packages{
		Chart:      charts.Chart{},
		Deployment: &nd,
		Service:    &nsvc,
	}
	pkg.Chart.ChartPackageID = c.GetChartPackageID()
	p.Require().NotZero(pkg.Chart.ChartPackageID)

	filepath := p.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := p.Yr.ReadYamlConfig(filepath)
	p.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &pkg.K8sDeployment)
	p.Require().Nil(err)
	err = pkg.ConvertDeploymentConfigToDB()
	p.Require().Nil(err)

	filepath = p.TestDirectory + "/mocks/test/service_eth_indexer.yaml"
	jsonBytes, err = p.Yr.ReadYamlConfig(filepath)
	err = json.Unmarshal(jsonBytes, &pkg.K8sService)
	p.Require().Nil(err)
	pkg.ConvertK8sServiceToDB()
	p.Assert().NotEmpty(pkg.Service)

	ctx = context.Background()
	q = sql_query_templates.NewQueryParam("InsertPackages", "table", "where", 1000, []string{})
	err = pkg.InsertPackages(ctx, q)
	p.Require().Nil(err)
}

func TestPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(PackagesTestSuite))
}
