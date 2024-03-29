package artemis_entities

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	iris_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/iris"
)

func (s *EntitiesTestSuite) TestSelectUserEntityCache() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	// Test data for insertion

	//rt := iris_models.RouteInfo{
	//	RoutePath: "https://example.com",
	//	RouteExt:  "/?url={q}",
	//	Payload: map[string]interface{}{
	//		"key": "value",
	//	},
	//}
	//
	//ht, err := HashWebRequestResultsAndParams(s.Ou, rt)
	//s.Require().Nil(err)
	uew := &UserEntityWrapper{
		UserEntity: UserEntity{
			Nickname: "17848f7a3d99c0e2fae7b0d474d054aa37b40849b01f0ea941707329f91110bf17848f7a3d99c0e2fae7b0d474d054aa37b40849b01f0ea941707329f91110bf",
			Platform: "mb-cache",
		},
		Ou: s.Ou,
	}
	uew.Ou.OrgID = 1685378241971196000
	ef := EntitiesFilter{
		Platform: "mb-cache",
	}
	et := ef.SetSinceOffsetNowTimestamp("hours", 3)
	fmt.Println(et)
	s.Require().NotZero(et)
	// ef 1711511508136852000

	// db 1711504617297472000
	err := SelectEntitiesCaches(ctx, uew, ef)
	s.Require().Nil(err)
	s.Assert().NotEmpty(uew.MdSlice)
	//
	s.Require().Len(uew.MdSlice, 1)
	mockFakeResp := "fake"
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
