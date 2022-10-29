package create_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateInfraTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateInfraTestSuite) TestInsertInfra() {
	tID, _ := s.SeedTopology()
	inf := NewCreateInfrastructure()
	inf.TopologyID = tID
	// manually seeding now
	chartPackageID := 6672899785140184951
	inf.ChartPackageID = chartPackageID
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertInfrastructureBase", "topology_infrastructure_components", "where", 1000, []string{})
	q.TableName = inf.GetTableName()
	q.Columns = inf.GetTableColumns()
	q.Values = []apps.RowValues{inf.GetRowValues("default")}
	err := inf.InsertInfrastructureBase(ctx, q)
	s.Require().Nil(err)
}

func TestCreateInfraTestSuite(t *testing.T) {
	suite.Run(t, new(CreateInfraTestSuite))
}
