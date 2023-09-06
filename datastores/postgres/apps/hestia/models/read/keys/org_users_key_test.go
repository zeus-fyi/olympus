package read_keys

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ReadOrgUsersKeyTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ReadOrgUsersKeyTestSuite) TestGetAuthedServices() {

	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	k := NewKeyReader()

	ou := org_users.NewOrgUser()
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	services, _, err := k.QueryUserAuthedServices(ctx, s.Tc.ProductionLocalTemporalBearerToken)
	s.Require().Nil(err)
	s.NotEmpty(services)
	s.Assert().Equal(7138983863666903883, k.OrgID)
	s.Assert().Equal(7138958574876245567, k.UserID)
	s.Assert().True(k.PublicKeyVerified)
}

func (s *ReadOrgUsersKeyTestSuite) TestVerifyBearerTokenService() {
	serviceName := "zeus"
	bearerToken := "n2cnnxxzsdsqhz6qwnq977z228pxrzjb6wxxlp7ksg5gnxhfdk54x8crmt7lpw4j"

	k := NewKeyReader()
	k.PublicKey = bearerToken

	err := k.VerifyUserTokenService(ctx, serviceName)
	s.Require().Nil(err)

	s.Assert().Equal(7138983863666903883, k.OrgID)
	s.Assert().Equal(7138958574876245567, k.UserID)
	s.Assert().True(k.PublicKeyVerified)

}
func (s *ReadOrgUsersKeyTestSuite) TestVerifyPassword() {

	k := NewKeyReader()
	k.PublicKey = s.Tc.AdminLoginPassword
	err := k.VerifyUserPassword(ctx, "alex@zeus.fyi")
	s.Require().Nil(err)
	s.Assert().Equal(s.Tc.ProductionLocalTemporalOrgID, k.OrgID)
	s.Assert().Equal(s.Tc.ProductionLocalTemporalUserID, k.UserID)
	s.Assert().True(k.PublicKeyVerified)
}

func (s *ReadOrgUsersKeyTestSuite) TestVerifyBearerToken() {
	oID, uID, bearerToken := s.NewTestOrgAndUserWithBearer()

	k := NewKeyReader()
	k.PublicKey = bearerToken

	err := k.VerifyUserBearerToken(ctx)
	s.Require().Nil(err)

	s.Assert().Equal(oID, k.OrgID)
	s.Assert().Equal(uID, k.UserID)
	s.Assert().True(k.PublicKeyVerified)
}

func TestReadOrgUsersKeyTestSuite(t *testing.T) {
	suite.Run(t, new(ReadOrgUsersKeyTestSuite))
}
