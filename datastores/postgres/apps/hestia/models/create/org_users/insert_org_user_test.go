package create_org_users

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	create_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type CreateOrgUserTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgUserTestSuite) TestInsertOrgUser() {
	ctx := context.Background()

	ou := NewCreateOrgUser()

	ou.OrgID = s.NewTestOrg()
	ou.UserID = s.NewTestUser()

	user := create_users.NewCreateUser()
	user.Metadata = `{"name": "test"}`

	b, err := json.Marshal(user.Metadata)
	s.Require().Nil(err)
	err = ou.InsertOrgUser(ctx, b)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ou.UserID)
}

func TestCreateOrgUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgUserTestSuite))
}
