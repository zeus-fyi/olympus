package create_state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateTopologyStateTestSuite struct {
	h hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateTopologyStateTestSuite) TestInsertTopologyState() {
	ctx := context.Background()
	topID, _ := s.SeedTopology()
	topState := NewCreateState()
	topState.TopologyID = topID
	topState.TopologyStatus = "InProgress"

	q := sql_query_templates.NewQueryParam("InsertState", "topologies_deployed", "where", 1000, []string{})
	q.TableName = topState.GetTableName()
	q.Columns = topState.GetTableColumns()
	q.Values = []apps.RowValues{topState.GetRowValues("default")}
	err := topState.InsertState(ctx, q)

	s.Require().Nil(err)
}

func TestCreateTopologyStateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTopologyStateTestSuite))
}
