package create_keys

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/keys"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type CreateKeyTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *CreateKeyTestSuite) TestUpdateUserPassword() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	pw := s.Tc.AdminLoginPassword
	userID := 1685378241971196000
	nk := NewCreateKey(userID, pw)
	nk.PublicKeyVerified = true
	nk.PublicKeyName = "userLoginPassword"
	nk.PublicKeyTypeID = keys.PassphraseKeyTypeID
	nk.CreatedAt = time.Now()
	err := nk.UpdateUserSignInKey(ctx)
	s.Require().Nil(err)
}
func (s *CreateKeyTestSuite) TestInsertUserPassword() {
	pw := s.Tc.AdminLoginPassword
	nk := NewCreateKey(s.Tc.ProductionLocalTemporalUserID, pw)
	nk.PublicKeyVerified = true
	nk.PublicKeyName = "userLoginPassword"
	nk.PublicKeyTypeID = keys.PassphraseKeyTypeID
	nk.CreatedAt = time.Now()

	err := nk.InsertUserKey(ctx)
	s.Require().Nil(err)
}

func (s *CreateKeyTestSuite) TestInsertUserSessionID() {
	sessionID := uuid.New()
	nk := NewCreateKey(s.Tc.ProductionLocalTemporalUserID, sessionID.String())
	nk.PublicKeyVerified = true
	nk.PublicKeyName = "sessionID"
	nk.CreatedAt = time.Now()

	_, err := nk.InsertUserSessionKey(ctx)
	s.Require().Nil(err)
}

func (s *CreateKeyTestSuite) TestInsertDiscordKey() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	token := "sssssss"
	nk := NewCreateKey(s.Tc.ProductionLocalTemporalUserID, token)
	nk.PublicKeyVerified = true
	nk.PublicKeyName = "discord"
	nk.CreatedAt = time.Now()

	err := nk.InsertDiscordKey(ctx)
	s.Require().Nil(err)
}

func (s *CreateKeyTestSuite) TestInsertUserKey() {
	uID := s.NewTestUser()
	nk := NewCreateKey(uID, "0x974C0c36265b7aa658b63A6121041AeE9e4DFd1b")
	nk.PublicKeyVerified = false
	nk.PublicKeyName = "test_key"
	nk.PublicKeyTypeID = keys.EcdsaKeyTypeID
	nk.CreatedAt = time.Now()

	err := nk.InsertUserKey(ctx)
	s.Require().Nil(err)
}

func TestCreateKeyTestSuite(t *testing.T) {
	suite.Run(t, new(CreateKeyTestSuite))
}
