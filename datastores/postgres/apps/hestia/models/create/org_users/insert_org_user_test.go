package create_org_users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	create_orgs "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/orgs"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type CreateOrgUserTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *CreateOrgUserTestSuite) TestInsertDemoOrgUserWithKey() {
	ctx := context.Background()

	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ou := OrgUser{}
	key, err := ou.InsertDemoOrgUserWithNewKey(ctx, []byte("{}"), "userDemo", EthereumEphemeryServiceID)
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
