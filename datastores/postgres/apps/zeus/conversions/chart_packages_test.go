package conversions

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/transformations"
	v1 "k8s.io/api/apps/v1"

	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

var PgTestDB apps.Db

type ChartPackagesTestSuite struct {
	base.TestSuite
	yr transformations.YamlReader
}

func (s *ChartPackagesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	s.yr = transformations.YamlReader{}
}

func (s *ChartPackagesTestSuite) TestChartPackagesQuery() {
	packageID := 0
	ctx := context.Background()
	conn := PgTestDB.InitPG(ctx, s.Tc.LocalDbPgconn)
	s.Assert().NotNil(conn)
	defer conn.Close()

	pkg, err := FetchQueryPackage(ctx, packageID)
	s.Require().Nil(err)
	s.Require().NotEmpty(pkg)
}

func (s *ChartPackagesTestSuite) TestConvertYamlConfig() {
	filepath := "/Users/alex/Desktop/Zeus/olympus/pkg/zeus/core/transformations/deployment.yaml"
	jsonBytes, err := s.yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)
	s.Require().Nil(err)
	s.Require().NotEmpty(d)
}

func TestChartPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ChartPackagesTestSuite))
}
