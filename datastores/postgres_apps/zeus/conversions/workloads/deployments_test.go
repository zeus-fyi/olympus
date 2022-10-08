package workloads

import (
	"encoding/json"
	"testing"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/test"
	v1 "k8s.io/api/apps/v1"

	"github.com/stretchr/testify/suite"

	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

type ConvertDeploymentPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertDeploymentPackagesTestSuite) TestConvertDeployment() {
	packageID := 0
	filepath := s.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig := ConvertDeploymentConfigToDB(d)
	s.Require().NotEmpty(dbDeploymentConfig)

	_ = dev_hacks.Use(packageID)
}

func TestConvertDeploymentPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertDeploymentPackagesTestSuite))
}
