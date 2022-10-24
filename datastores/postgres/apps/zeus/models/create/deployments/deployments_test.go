package deployments

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/workloads"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type ConvertDeploymentPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertDeploymentPackagesTestSuite) TestConvertDeploymentAndInsert() {
	filepath := s.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig, err := workloads.ConvertDeploymentConfigToDB(d)
	s.Require().Nil(err)
	s.Require().NotEmpty(dbDeploymentConfig)

	mockC, err := mockChart()
	s.Require().Nil(err)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertDeployment", "table", "where", 1000, []string{})
	dbDeploy := NewDeploymentConfigForDB(dbDeploymentConfig)
	err = dbDeploy.InsertDeployment(ctx, q, mockC)
	s.Require().Nil(err)
}

func mockChart() (create.Chart, error) {
	ns := sql.NullString{}
	c := create.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err := c.InsertChart(ctx, q, c)
	return c, err
}

func TestConvertDeploymentPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertDeploymentPackagesTestSuite))
}
