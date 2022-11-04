package create_class_type

import (
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
)

type CreateTopologyClassTestSuite struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateTopologyClassTestSuite) TestInsertTopologyClass() {
	//tc := NewCreateTopologyClass()
	//cl := tc.GetClassDefinition()
	//
	//cl.TopologyClassID = s.Ts.UnixTimeStampNow()
	//cl.TopologyClassName = fmt.Sprintf("test_topology_class_%d", cl.TopologyClassTypeID)
	//cl.TopologyClassTypeID = class_type.SkeletonBaseClassTypeID
	//
	//ctx := context.Background()
	//q := sql_query_templates.NewQueryParam("InsertTopologyClass", "topology_class_types", "where", 1000, []string{})
	//q.TableName = cl.GetTableName()
	//q.Columns = cl.GetTableColumns()
	//q.Values = []apps.RowValues{cl.GetRowValues("default")}
	//
	//tc.SetClassDefinition(cl)
	//err := tc.InsertTopologyClass(ctx, q)
	//s.Require().Nil(err)
}

func TestCreateTopologyClassTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTopologyClassTestSuite))
}
