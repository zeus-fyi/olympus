package networking

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"

	"github.com/zeus-fyi/olympus/datastores/postgres_apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
)

type ConvertServiceTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertServiceTestSuite) TestConvertService() {
	packageID := 0
	filepath := s.TestDirectory + "/mocks/test/service_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var svc *v1.Service
	err = json.Unmarshal(jsonBytes, &svc)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)

	dbServiceConfig := ConvertServiceConfigToDB(svc)
	s.Require().NotEmpty(dbServiceConfig)

	_ = dev_hacks.Use(packageID)
}

func TestConvertServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertServiceTestSuite))
}
