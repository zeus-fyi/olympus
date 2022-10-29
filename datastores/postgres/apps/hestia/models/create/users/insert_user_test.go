package create_users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateUsersTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateUsersTestSuite) TestInsertUser() {
	var ts chronos.Chronos
	u := NewCreateUser()
	u.UserID = ts.UnixTimeStampNow()

	ctx := context.Background()
	q := sql_query_templates.NewQueryParam("InsertUser", "users", "where", 1000, []string{})
	q.TableName = u.GetTableName()
	q.Columns = u.GetTableColumns()
	q.Values = []apps.RowValues{u.GetRowValues("default")}
	err := u.InsertUser(ctx, q)
	s.Require().Nil(err)

}

func TestCreateUsersTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUsersTestSuite))
}
