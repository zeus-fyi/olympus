package containers

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/workloads"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/deployments"
	autogen_structs "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/structs/autogen"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/rand"
)

type PodContainersGroupTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PodContainersGroupTestSuite) TestConvertDeploymentAndInsert() {
	filepath := p.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := p.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	p.Require().Nil(err)
	p.Require().NotEmpty(d)

	dbDeploymentConfig := workloads.ConvertDeploymentConfigToDB(d)
	p.Require().NotEmpty(dbDeploymentConfig)

	mockC, err := mockChart()
	p.Require().Nil(err)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertPodResourceContainers", "table", "where", 1000, []string{})
	dbDeploy := deployments.NewDeploymentConfigForDB(dbDeploymentConfig)
	err = dbDeploy.InsertDeployment(ctx, q, mockC)
	p.Require().Nil(err)

	// TODO, extract the pod resource group then use for test and remove excess copied here from deployment test
}

func mockChart() (create.Chart, error) {
	ns := sql.NullString{}
	c := create.Chart{autogen_structs.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertPodResourceContainers", "table", "where", 1000, []string{})
	err := c.InsertChart(ctx, q, c)
	return c, err
}

func TestPodContainersGroupTestSuite(t *testing.T) {
	suite.Run(t, new(PodContainersGroupTestSuite))
}
