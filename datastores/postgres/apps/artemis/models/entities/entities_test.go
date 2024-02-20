package entities

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type EntitiesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *EntitiesTestSuite) TestSelectUserEntityWithMd() {
	res, err := SelectUserMetadataByNicknameAndPlatform(ctx,
		"test-nickname", "test-platform", []string{"test-label1", "test-label2"},
		1)
	s.Require().Nil(err)
	s.Require().Len(res, 1)

	tsNow := chronos.Chronos{}
	res, err = SelectUserMetadataByNicknameAndPlatform(ctx,
		"test-nickname", "test-platform", []string{"test-label1", "test-label2"},
		tsNow.UnixTimeStampNow())
	s.Require().Nil(err)
	s.Require().Len(res, 0)
}

func (s *EntitiesTestSuite) TestInsertUserEntity() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Test data for insertion
	testUserEntity := UserEntityWrapper{
		UserEntity: UserEntity{
			Nickname:  "test-nickname",
			Platform:  "test-platform",
			FirstName: nil, // Assuming no first name to demonstrate handling of NULL
			LastName:  nil, // Assuming no last name to demonstrate handling of NULL
		},
		Ou: s.Ou,
		MdSlice: []UserEntityMetadata{
			{
				JsonData: json.RawMessage(`{"key": "value"}`),
				TextData: nil, // Demonstrating handling of NULL
				Labels: []UserEntityMetadataLabel{
					{Label: "test-label1"},
					{Label: "test-label2"},
				},
			},
		},
	}

	err := InsertUserEntityLabeledMetadata(ctx, &testUserEntity)
	s.Require().Nil(err)
	s.Require().NotZero(testUserEntity.EntityID)
}

func TestEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(EntitiesTestSuite))
}
