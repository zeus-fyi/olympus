package read_keys

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	create_org_users "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/create/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ReadOrgUsersKeyTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ReadOrgUsersKeyTestSuite) TestVerifyBearerTokenService() {
	serviceName := create_org_users.EthereumEphemeryService
	bearerToken := "75ghwvjxw9lnzlz57nlblbrb2kw4v99nk7nbx2zbq5pwj4dntdpl962j8qfj6rg8456zdcghgnwmnp46sj6bd7rcgcc8ddkrvl2nvqbb5dg2nwtzj6vph7lj"

	k := NewKeyReader()
	k.PublicKey = bearerToken

	err := k.VerifyUserBearerTokenService(ctx, serviceName)
	s.Require().Nil(err)

	s.Assert().Equal(create_org_users.UserDemoOrgID, k.OrgID)
	s.Assert().Equal(1677099883023927040, k.UserID)
	s.Assert().True(k.PublicKeyVerified)

	serviceName = create_org_users.EthereumMainnetService
	err = k.VerifyUserBearerTokenService(ctx, serviceName)
	s.Require().Error(err)
	s.Assert().False(k.PublicKeyVerified)
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
