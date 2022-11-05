package containers

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	v1 "k8s.io/api/apps/v1"
)

type PodContainersGroupTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PodContainersGroupTestSuite) TestContainersInsertFromParsedDeploymentFile() {
	//ctx := context.Background()
	filepath := p.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := p.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	p.Require().Nil(err)
	p.Require().NotEmpty(d)

	//ps := PodTemplateSpec{}
	//pts := d.Spec.Template.Spec
	//dbDeploymentConfig, err := ps.ConvertPodTemplateSpecConfigToDB(&pts)
	//p.Require().Nil(err)
	//p.Require().NotEmpty(dbDeploymentConfig)
	//
	//q := sql_query_templates.NewQueryParam("InsertPodResourceContainers", "table", "where", 1000, []string{})
	//
	//// TODO update with another chart later?
	//c := charts.Chart{}
	//// specific to test, above code is just setting up
	//err = dbDeploymentConfig.InsertPodTemplateSpec(ctx, q, &c)
	//p.Require().Nil(err)
}

func TestPodContainersGroupTestSuite(t *testing.T) {
	suite.Run(t, new(PodContainersGroupTestSuite))
}
