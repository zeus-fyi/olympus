package package_inserts

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/deployments"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type PackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (p *PackagesTestSuite) TestInsert() {
	pkg := Packages{
		Chart:      charts.Chart{},
		Deployment: deployments.Deployment{},
		Service:    networking.Service{},
	}
	filepath := p.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := p.Yr.ReadYamlConfig(filepath)
	p.Require().Nil(err)
	err = json.Unmarshal(jsonBytes, &pkg.K8sDeployment)
	p.Require().Nil(err)

	p.Assert().NotEmpty(pkg)

}
func TestPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(PackagesTestSuite))
}
