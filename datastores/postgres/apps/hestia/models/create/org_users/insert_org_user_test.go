package create_org_users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_keys "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/keys"
	create_orgs "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/orgs"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"k8s.io/apimachinery/pkg/util/rand"
)

type CreateOrgUserTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgUserTestSuite) TestInsertDemoOrgUserWithSignUp() {
	ctx := context.Background()
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := OrgUser{}
	us := UserSignup{
		FirstName:    "alex",
		LastName:     "g",
		EmailAddress: rand.String(10) + "@zeus.fyi",
		Password:     "password",
	}
	userID, err := DoesUserExist(ctx, us.EmailAddress)
	s.Require().Nil(err)
	s.Assert().Zero(userID)

	key, err := ou.InsertSignUpOrgUserAndVerifyEmail(ctx, us)
	s.Require().Nil(err)
	s.Assert().NotEmpty(key)

	err = create_keys.UpdateKeysFromVerifyEmail(ctx, key)
	s.Require().Nil(err)

	userID, err = DoesUserExist(ctx, us.EmailAddress)
	s.Require().Nil(err)
	s.Assert().NotZero(userID)
}

func (s *CreateOrgUserTestSuite) TestInsertDemoOrgUserWithKey() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	ou := OrgUser{}
	key, err := ou.InsertDemoOrgUserWithNewKey(ctx, []byte(`{"email": "alex@zeus.fyi", "ethereumAddress": "0x974C0c36265b7aa658b63A6121041AeE9e4DFd1b", "validatorCount": "3"}`), "userDemo", EthereumEphemeryServiceID)
	s.Require().Nil(err)
	s.Assert().NotEmpty(key)
}

func (s *CreateOrgUserTestSuite) TestInsertOrgUserWithKeyToService() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := OrgUser{}
	ou.OrgID = 1677111877017029000
	key, err := ou.InsertOrgUserWithNewKeyForService(ctx, []byte(`{"name": "zeus-webhooks"}`), "zeus-webhooks", ZeusWebhooksServiceID)
	s.Require().Nil(err)
	s.Assert().NotEmpty(key)
}

func (s *CreateOrgUserTestSuite) TestInsertOrg() {
	ctx := context.Background()
	var ts chronos.Chronos

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	o := create_orgs.NewCreateOrg()
	o.OrgID = ts.UnixTimeStampNow()

	o.Org.Name = "zeus-webhooks"
	err := o.InsertOrg(ctx)
	s.Require().Nil(err)
}

func TestCreateOrgUserTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrgUserTestSuite))
}
