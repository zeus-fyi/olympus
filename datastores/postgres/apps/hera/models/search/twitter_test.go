package hera_search

import (
	twitter2 "github.com/cvcio/twitter"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestInsertTweetSearch() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	query := `(("Kubernetes" OR "k8s" OR "#kube" OR "container orchestration") -is:retweet (has:links OR has:media OR has:mentions) (lang:en OR lang:es))`
	resp, err := InsertTwitterSearchQuery(ctx, ou, defaultTwitterSearchGroupName, query, 100)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)

}

func (s *SearchAITestSuite) TestInsertIncomingTweets() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	searchID := 1700514008821625088
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
