package artemis_entities

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (s *EntitiesTestSuite) TestSelectUserEntityCache() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Test data for insertion

	rt := iris_models.RouteInfo{
		RoutePath: "https://example.com",
		RouteExt:  "/?url={q}",
		Payload: map[string]interface{}{
			"key": "value",
		},
	}

	ht, err := HashWebRequestResultsAndParams(s.Ou, rt)
	s.Require().Nil(err)
	mockFakeResp := "fake"
	uew := &UserEntityWrapper{
		UserEntity: UserEntity{
			Nickname: ht.RequestCache,
			Platform: "mb-cache",
		},
		Ou: s.Ou,
	}
	err = SelectEntitiesCaches(ctx, uew)
	s.Require().Nil(err)
	s.Assert().NotEmpty(uew.MdSlice)

	s.Require().Len(uew.MdSlice, 1)
	s.Equal(uew.MdSlice[0].TextData, aws.String(mockFakeResp))
}

func (s *EntitiesTestSuite) TestInsertUserEntityCache() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	// Test data for insertion

	rt := iris_models.RouteInfo{
		RoutePath: "https://example.com",
		RouteExt:  "/?url={q}",
		Payload: map[string]interface{}{
			"key": "value",
		},
	}

	ht, err := HashWebRequestResultsAndParams(s.Ou, rt)
	s.Require().Nil(err)

	mockFakeResp := "fake"
	uew := UserEntityWrapper{
		UserEntity: UserEntity{
			Nickname: ht.RequestCache,
			Platform: "mb-cache",
			MdSlice: []UserEntityMetadata{
				{
					JsonData: nil,
					TextData: aws.String(mockFakeResp),
					Labels: []UserEntityMetadataLabel{
						{
							Label: ht.RequestCache,
						},
					},
				},
			},
		},
		Ou: s.Ou,
	}

	_, err = InsertEntitiesCaches(ctx, &uew)
	s.Require().Nil(err)
	//
	//resp, err := SelectEntitiesCaches(ctx, uew, ht)
	//s.Require().Nil(err)
	//s.Require().Nil(resp)
}
