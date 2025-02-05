package hera_search

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

var ctx = context.Background()

type SearchAITestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

func (s *SearchAITestSuite) TestSelectAll() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	ts, err := GetAllActiveSearchIndexers(ctx)
	s.Require().Nil(err)
	s.Assert().NotNil(ts)
}

func (s *SearchAITestSuite) TestSelectAllByOrd() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	ts, err := GetSearchIndexersByOrg(ctx, ou)
	s.Require().Nil(err)
	s.Assert().NotNil(ts)
}

//func (s *SearchAITestSuite) TestSelectTelegramResults() {
//	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
//	ou := org_users.OrgUser{}
//	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
//	ou.UserID = s.Tc.ProductionLocalTemporalUserID
//
//	si := TimeInterval{}
//	si[0] = time.Now().AddDate(0, 0, -60)
//
//	fmt.Println(si[0].Unix())
//	si[1] = time.Now()
//	fmt.Println(si[1].Unix())
//
//	sp := AiSearchParams{
//		GroupFilter:    "Ze",
//		SearchInterval: si,
//	}
//	res, err := SearchTelegram(ctx, ou, sp)
//	s.Require().Nil(err)
//	s.Assert().NotZero(res)
//}

//func (s *SearchAITestSuite) TestHashSearchResultAndParams() {
//	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
//
//	sp := AiSearchParams{
//		GroupFilter: "Ze",
//	}
//	ou := org_users.OrgUser{}
//	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
//	ou.UserID = s.Tc.ProductionLocalTemporalUserID
//
//	res, err := SearchTelegram(ctx, ou, sp)
//	s.Require().Nil(err)
//	s.Assert().NotZero(res)
//
//	hash, err := HashParams(ou.OrgID, []interface{}{sp, res})
//	s.Require().Nil(err)
//	s.Assert().NotEmpty(hash)
//	fmt.Println(hash)
//	response := openai.ChatCompletionResponse{
//		Choices: []openai.ChatCompletionChoice{
//			{
//				Message: openai.ChatCompletionMessage{
//					Role:    "chat",
//					Content: "sdfsdfsdfsd",
//					Name:    "kjkdd",
//				},
//			},
//		},
//	}
//	hash2, err := HashParams(ou.OrgID, []interface{}{sp, res, response})
//	s.Require().Nil(err)
//	s.Assert().NotEmpty(hash)
//	s.Assert().NotEqual(hash, hash2)
//	fmt.Println(hash2)
//
//	hrp, err := HashAiSearchResponseResultsAndParams(ou, response, sp, res)
//	s.Require().Nil(err)
//	s.Assert().NotNil(hrp)
//	s.Assert().Equal(hash, hrp.SearchAndResultsHash)
//	s.Assert().Equal(hash2, hrp.SearchAnalysisHash)
//
//	err = InsertCompletionResponseChatGptFromSearch(ctx, ou, response, sp, res)
//	s.Require().Nil(err)
//}

// HashSearchResultAndParams
func TestSearchAITestSuite(t *testing.T) {
	suite.Run(t, new(SearchAITestSuite))
}
