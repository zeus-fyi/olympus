package create_orgs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateOrgsTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgsTestSuite) TestInsertOrg() {
	var ts chronos.Chronos
	o := NewCreateOrg()
	o.OrgID = ts.UnixTimeStampNow()

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertOrg", "orgs", "where", 1000, []string{})
	q.TableName = o.GetTableName()
	q.Columns = o.GetTableColumns()
	q.Values = []apps.RowValues{o.GetRowValues("default")}
	err := o.InsertOrg(ctx)
	s.Require().Nil(err)

}

func TestCreateOrgsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgsTestSuite))
}
