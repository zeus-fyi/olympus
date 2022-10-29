package create_class_type

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/classes/class_type"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateTopologyClass struct {
	b hestia_test.BaseHestiaTestSuite
	conversions_test.ConversionsTestSuite
}

func (s *CreateTopologyClass) TestInsertTopologyClass() {
	tc := NewCreateTopologyClass()
	cl := tc.GetClassDefinition()

	cl.TopologyClassID = s.Ts.UnixTimeStampNow()
	cl.TopologyClassName = fmt.Sprintf("test_topology_class_%d", cl.TopologyClassTypeID)
	cl.TopologyClassTypeID = class_type.SkeletonBaseClassTypeID

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertTopologyClass", "topology_class_types", "where", 1000, []string{})
	q.TableName = cl.GetTableName()
	q.Columns = cl.GetTableColumns()
	q.Values = []apps.RowValues{cl.GetRowValues("default")}

	tc.SetClassDefinition(cl)
	err := tc.InsertTopologyClass(ctx, q)
	s.Require().Nil(err)
}

func TestCreateTopologyClass(t *testing.T) {
	suite.Run(t, new(CreateTopologyClass))
}
