package create_kns

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	conversions_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateKnsTestSuite struct {
	conversions_test.ConversionsTestSuite
}

func (s *CreateKnsTestSuite) TestInsertKns() {

	topID, _ := s.SeedTopology()
	newKns := NewCreateKns()
	newKns.CloudProvider = "do"
	newKns.Region = "sfo3"
	newKns.Context = "context"
	newKns.Env = "test"
	newKns.Namespace = "testnamespace"
	newKns.TopologyID = topID
	fmt.Println(topID)
	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertKns", "kns", "where", 1000, []string{})
	q.TableName = newKns.GetTableName()
	q.Columns = newKns.GetTableColumns()
	q.Values = []apps.RowValues{newKns.GetRowValues("default")}
	err := newKns.InsertKns(ctx, q)
	s.Require().Nil(err)

	newKns.Namespace = "new"
	err = InsertKns(ctx, &newKns.TopologyKubeCtxNs)
	s.Require().Nil(err)
}

func TestCreateKnsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateKnsTestSuite))
}
