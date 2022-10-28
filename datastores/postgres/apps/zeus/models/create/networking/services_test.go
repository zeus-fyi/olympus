package networking

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/charts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/networking/services"
	create_charts "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/charts"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/rand"
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
		services.NewService(),
	}
	zeusService.K8sService = svc
	zeusService.ConvertK8sServiceToDB()
	s.Require().Nil(err)

	s.Require().NotEmpty(zeusService.Metadata)
	s.Require().NotEmpty(zeusService.ServiceSpec)
	s.Require().NotEmpty(zeusService.Service.ServiceSpec.Ports)

	ns := sql.NullString{}
	c := create_charts.Chart{ChartPackages: autogen_bases.ChartPackages{
		ChartPackageID:   0,
		ChartName:        rand.String(10),
		ChartVersion:     rand.String(10),
		ChartDescription: ns,
	}}
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertChart", "table", "where", 1000, []string{})
	err = c.InsertChart(ctx, q)
	s.Require().Nil(err)

	mockC := charts.Chart{}
	mockC.ChartPackageID = c.GetChartPackageID()
	err = zeusService.InsertService(ctx, q, &mockC)
	s.Require().Nil(err)
}

func TestNetworkingTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkingTestSuite))
}
