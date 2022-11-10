package create_deployments

import (
	"encoding/json"
	"testing"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/misc/dev_hacks"
	v1 "k8s.io/api/apps/v1"

	"github.com/stretchr/testify/suite"
)

type ConvertDeploymentPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertDeploymentPackagesTestSuite) TestConvertDeployment() {
	packageID := 0
	filepath := s.TestDirectory + "/mocks/demo/deployment.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig, err := ConvertDeploymentConfigToDB(d)
	s.Require().Nil(err)

	s.Require().NotEmpty(dbDeploymentConfig.Metadata.Name.ChartSubcomponentValue)
	s.Require().NotEmpty(dbDeploymentConfig)

	s.Require().NotEmpty(dbDeploymentConfig.Spec)

	_ = dev_hacks.Use(packageID)
}

func TestConvertDeploymentPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertDeploymentPackagesTestSuite))
}
