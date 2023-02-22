package create_org_users

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_orgs "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/orgs"
	create_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type CreateOrgUserTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgUserTestSuite) TestInsertOrgUserWithKey() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := OrgUser{}
	key, err := ou.InsertDemoOrgUserWithNewKey(ctx, []byte("{}"), "userdemotestkey", EthereumEphemeryServiceID)
	s.Require().Nil(err)
	s.Assert().NotEmpty(key)
}

func (s *CreateOrgUserTestSuite) TestInsertOrgUser() {
	ctx := context.Background()
	var ts chronos.Chronos

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	o := create_orgs.NewCreateOrg()
	o.OrgID = ts.UnixTimeStampNow()

	o.Org.Name = "orgName"
	err := o.InsertOrg(ctx)
	s.Require().Nil(err)

	user := create_users.NewCreateUser()
	user.Metadata = `{"name": "usersname"}`

	ou := NewCreateOrgUserWithOrgID(o.OrgID)
	err = ou.InsertOrgUser(ctx, []byte(user.Metadata))
	b, err := json.Marshal(user.Metadata)
	s.Require().Nil(err)
	err = ou.InsertOrgUser(ctx, b)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ou.UserID)
}

func TestCreateOrgUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgUserTestSuite))
}
