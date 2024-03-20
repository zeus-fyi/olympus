package artemis_entities

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/suite"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
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
		"", "", nil, 0)

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
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	var user, pw string
	sn := "api-twillio"
	ps, err := aws_secrets.GetMockingbirdPlatformSecrets(context.Background(), s.Ou, sn)
	if len(ps.TwillioAccount) > 0 {
		user = ps.TwillioAccount
	}
	if len(ps.TwillioAuth) > 0 {
		pw = ps.TwillioAuth
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: user,
		Password: pw,
	})

	tvUnix, err := SelectHighestLabelIdForLabelAndPlatform(ctx, s.Ou, "twillio", "indexer:twillio")
	s.Require().Nil(err)
	//s.Require().NotZero(res)
	if tvUnix == 0 {
		tvUnix = int(time.Now().Add(-time.Hour * 24 * 7).UnixNano())
	}
	ch := chronos.Chronos{}

	tv := ch.ConvertUnixTimeStampToDate(tvUnix)
	fmt.Println(tv.String())

	resp, err := client.Api.ListMessage(&twilioApi.ListMessageParams{
		DateSentAfter: &tv,
		PageSize:      aws.Int(1000),
		Limit:         aws.Int(100),
	})
	fmt.Println(resp, err)
	s.Require().Nil(err)

}

func TestEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(EntitiesTestSuite))
}
