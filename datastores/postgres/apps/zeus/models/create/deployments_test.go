package create

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/workloads"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/dev_hacks"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	v1 "k8s.io/api/apps/v1"
)

type ConvertDeploymentPackagesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *ConvertDeploymentPackagesTestSuite) TestConvertDeploymentAndInsert() {
	packageID := 0
	filepath := s.TestDirectory + "/mocks/test/deployment_eth_indexer.yaml"
	jsonBytes, err := s.Yr.ReadYamlConfig(filepath)

	var d *v1.Deployment
	err = json.Unmarshal(jsonBytes, &d)

	s.Require().Nil(err)
	s.Require().NotEmpty(d)

	dbDeploymentConfig := workloads.ConvertDeploymentConfigToDB(d)
	s.Require().NotEmpty(dbDeploymentConfig)

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertDeployment", "table", "where", 1000, []string{})
	err = InsertDeployment(ctx, q, dbDeploymentConfig)

	s.Require().Nil(err)
	_ = dev_hacks.Use(packageID)
}

func TestConvertDeploymentPackagesTestSuite(t *testing.T) {
	suite.Run(t, new(ConvertDeploymentPackagesTestSuite))
}
