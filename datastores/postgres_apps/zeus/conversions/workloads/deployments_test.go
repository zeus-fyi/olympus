package workloads

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/conversions"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	v1 "k8s.io/api/apps/v1"
)

type DeploymentPackagesTestSuite struct {
	conversions.ChartPackagesTestSuite
}

func (s *DeploymentPackagesTestSuite) TestDeploymentPackagesConversion() {
	packageID := 0
	filepath := "/Users/alex/Desktop/Zeus/olympus/pkg/zeus/core/transformations/deployment.yaml"
	jsonBytes, err := s.yr.ReadYamlConfig(filepath)

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
