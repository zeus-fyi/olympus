package read_keys

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ReadOrgUsersKeyTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *ReadOrgUsersKeyTestSuite) TestVerifyBearerToken() {
	ctx := context.Background()

	oID, uID, bearerToken := s.NewTestOrgAndUserWithBearer()

	k := NewKeyReader()
	k.PublicKey = bearerToken

	err := k.VerifyUserBearerToken(ctx)
	s.Require().Nil(err)

	s.Assert().Equal(oID, k.OrgID)
	s.Assert().Equal(uID, k.OrgUser.UserID)
	s.Assert().True(k.PublicKeyVerified)
}

func TestReadOrgUsersKeyTestSuite(t *testing.T) {
	suite.Run(t, new(ReadOrgUsersKeyTestSuite))
}
