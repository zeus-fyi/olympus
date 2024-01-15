package hera_search

import (
	"time"

	twitter2 "github.com/cvcio/twitter"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestInsertTweetSearch() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en))`
	resp, err := InsertTwitterSearchQuery(ctx, ou, defaultTwitterSearchGroupName, query, 100)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}

func (s *SearchAITestSuite) TestSelectTwitterSearchQuery() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en))`
	ts, err := SelectTwitterSearchQuery(ctx, ou, defaultTwitterSearchGroupName)
	s.Require().Nil(err)
	s.Assert().NotZero(ts.SearchID)
	s.Assert().Equal(100, ts.MaxResults)
	s.Assert().Equal(query, ts.Query)
	s.Assert().NotZero(ts.MaxTweetID)
}

func (s *SearchAITestSuite) TestInsertIncomingTweets() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	searchID := 1700515651883687936
	tweets := []*twitter2.Tweet{
		{
			ID:   "1700514008821625088",
			Text: "tweet 1",
		},
		{
			ID:   "1700514008821625081",
			Text: "tweet 2",
		},
	}
	resp, err := InsertIncomingTweets(ctx, searchID, tweets)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}

func (s *SearchAITestSuite) TestSelectTweets() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	si := artemis_orchestrations.Window{}
	si.Start = time.Now().AddDate(0, 0, -30)
	si.End = time.Now()

	// Call the function
	sp := AiSearchParams{
		Retrieval: artemis_orchestrations.RetrievalItem{
			RetrievalName:  "",
			RetrievalGroup: "",
			RetrievalItemInstruction: artemis_orchestrations.RetrievalItemInstruction{
				DiscordFilters: nil,
			},
		},
		Window: si,
	}

	results, err := SearchTwitter(ctx, ou, sp)

	resp := FormatSearchResultsV3(results)
	s.NotEmpty(resp)
	// Assert expected outcomes
	s.Require().NoError(err, "SearchTwitter should not return an error")
	s.Require().NotNil(results, "Results should not be nil")

	ou.OrgID = 0
	results, err = SearchTwitter(ctx, ou, sp)
	s.Require().Nil(err)
	s.Require().Nil(results)
}
