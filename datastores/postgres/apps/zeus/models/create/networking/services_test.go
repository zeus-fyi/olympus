package networking

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/core/v1"
)

type NetworkingTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *NetworkingTestSuite) TestConvertServiceAndInsert() {
	filepath := s.TestDirectory + "/mocks/test/service_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var svc v1.Service
	err = json.Unmarshal(jsonBytes, &svc)

	s.Require().Nil(err)
	s.Require().NotEmpty(svc)

	zeusService := Service{
		networking.NewService(),
	}
	zeusService.K8sService = svc
	zeusService.ConvertK8sServiceToDB()
	s.Require().Nil(err)

	s.Require().NotEmpty(zeusService.Metadata)
	s.Require().NotEmpty(zeusService.ServiceSpec)
	s.Require().NotEmpty(zeusService.Service.ServiceSpec.Ports)

	mockC := s.MockChart()
	s.Require().Nil(err)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertService", "table", "where", 1000, []string{})

	err = zeusService.InsertService(ctx, q, &mockC)
	s.Require().Nil(err)
}

func TestNetworkingTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkingTestSuite))
}
