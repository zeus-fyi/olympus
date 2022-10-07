package workloads

import (
	"encoding/json"
	"testing"

	v1 "k8s.io/api/apps/v1"

	"github.com/stretchr/testify/suite"

	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions/test"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

type DeploymentPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *DeploymentPackagesTestSuite) TestDeploymentPackagesConversion() {
	packageID := 0
	filepath := "/Users/alex/Desktop/Zeus/olympus/pkg/zeus/core/transformations/deployment.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig := ConvertDeploymentConfigToDB(d)
	s.Require().NotEmpty(dbDeploymentConfig)

	_ = dev_hacks.Use(packageID)
}

func TestDeploymentPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentPackagesTestSuite))
}
