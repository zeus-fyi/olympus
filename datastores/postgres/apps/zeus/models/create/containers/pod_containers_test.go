package containers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/workloads"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/deployments"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
)

type PodContainersGroupTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PodContainersGroupTestSuite) TestContainersInsertFromParsedDeploymentFile() {
	ctx := context.Background()
	filepath := p.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := p.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	p.Require().Nil(err)
	p.Require().NotEmpty(d)

	dbDeploymentConfig := workloads.ConvertDeploymentConfigToDB(d)
	p.Require().NotEmpty(dbDeploymentConfig)

	q := sql_query_templates.NewQueryParam("InsertPodResourceContainers", "table", "where", 1000, []string{})
	dbDeploy := deployments.NewDeploymentConfigForDB(dbDeploymentConfig)

	// TODO remove dummy hardcode once better test setup exists
	setDummyPodSpecHeader(&dbDeploy)

	// specific to test, above code is just setting up
	dbDeployPodSpecContainers := NewPodContainersGroupForDB(dbDeploy.Spec.Template)
	dummyPodSpecClassTypeID := dbDeploy.Spec.DeploymentSpec.Template.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID
	err = dbDeployPodSpecContainers.InsertPodContainerGroup(ctx, q, dummyPodSpecClassTypeID)
	p.Require().Nil(err)
}

func setDummyPodSpecHeader(d *deployments.Deployment) {
	d.Spec.DeploymentSpec.Template.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentParentClassTypeID = 1666564843324726081
	d.Spec.DeploymentSpec.Template.Spec.PodTemplateSpecClassDefinition.ChartSubcomponentChildClassTypeID = 0
	return
}

func TestPodContainersGroupTestSuite(t *testing.T) {
	suite.Run(t, new(PodContainersGroupTestSuite))
}
