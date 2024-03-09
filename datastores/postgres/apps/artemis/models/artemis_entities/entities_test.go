package artemis_entities

import (
	"context"
	"encoding/json"
	"fmt"
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
	res, err := SelectUserMetadataByProvidedFields(ctx, s.Ou,
		"test-nickname", "test-platform", []string{"test-label1", "test-label2"},
		1)
	s.Require().Nil(err)
	s.Require().Len(res, 1)

	tsNow := chronos.Chronos{}
	res, err = SelectUserMetadataByProvidedFields(ctx, s.Ou,
		"test-nickname", "test-platform", []string{"test-label1", "test-label2"},
		tsNow.UnixTimeStampNow())
	s.Require().Nil(err)
	s.Require().Len(res, 0)
}

func (s *EntitiesTestSuite) TestSelectEntitiesWithAnyData() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	res, err := SelectUserMetadataByProvidedFields(ctx, s.Ou,
		"", "", nil, -36000)

	s.Require().Nil(err)
	s.Require().NotEmpty(res)
}

func (s *EntitiesTestSuite) TestInsertUserEntity() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Test data for insertion
	testUserEntity := UserEntityWrapper{
		UserEntity: UserEntity{
			Nickname: "+17575828406",
			Platform: "twillio",
			MdSlice: []UserEntityMetadata{
				{
					JsonData: json.RawMessage(`{"vt": "me"}`),
					TextData: nil, // Demonstrating handling of NULL
					Labels: []UserEntityMetadataLabel{
						{Label: "virginia-tech"},
						{Label: "mockingbird"},
						{Label: "twillio"},
						{Label: "indexer:twillio"},
					},
				},
			},
		},
		Ou: s.Ou,
	}

	err := InsertUserEntityLabeledMetadata(ctx, &testUserEntity)
	s.Require().Nil(err)
	s.Require().NotZero(testUserEntity.EntityID)
}

func (s *EntitiesTestSuite) TestSelectLatestTwillioIndexTime() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)

	res, err := SelectHighestLabelIdForLabelAndPlatform(ctx, s.Ou, "twillio", "indexer:twillio")
	s.Require().Nil(err)
	s.Require().NotZero(res)

	ch := chronos.Chronos{}

	tv := ch.ConvertUnixTimeStampToDate(res)
	fmt.Println(tv.String())
}

func TestEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(EntitiesTestSuite))
}
