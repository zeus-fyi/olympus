package create_topologies

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type TopologiesTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (t *TopologiesTestSuite) TestInsert() {
	top := NewCreateTopology()
	top.Name = "testTopology"
	top.TopologyID = t.Ts.UnixTimeStampNow()
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertTopology", "topologies", "where", 1000, []string{})
	q.TableName = top.GetTableName()
	q.Columns = top.GetTableColumns()
	q.Values = []apps.RowValues{top.GetRowValues("default")}
	err := top.InsertTopology(ctx, q)
	t.Require().Nil(err)
}

func TestTopologiesTestSuite(t *testing.T) {
	suite.Run(t, new(TopologiesTestSuite))
}
