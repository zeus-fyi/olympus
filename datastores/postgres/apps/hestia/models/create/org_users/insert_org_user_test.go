package create_org_users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type CreateOrgUserTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgUserTestSuite) TestInsertOrgUser() {
	ctx := context.Background()

	ou := NewCreateOrgUser()
	ou.OrgID = s.NewTestOrg()
	ou.UserID = s.NewTestUser()
	quo := sql_query_templates.NewQueryParam("NewTestOrgUser", "org_users", "where", 1000, []string{})
	quo.TableName = ou.GetTableName()
	quo.Columns = ou.GetTableColumns()
	quo.Values = []apps.RowValues{ou.GetRowValues("default")}

	err := ou.InsertOrgUser(ctx, quo)
	s.Require().Nil(err)
}

func TestCreateOrgUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgUserTestSuite))
}
