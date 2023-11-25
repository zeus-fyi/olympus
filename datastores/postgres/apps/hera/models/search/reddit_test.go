package hera_search

import (
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestInsertRedditSearchQuery() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	query := "devops"
	resp, err := InsertRedditSearchQuery(ctx, ou, defaultTwitterSearchGroupName, query, 100)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}

func (s *SearchAITestSuite) TestSelectRedditSearchQuery() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	ts, err := SelectRedditSearchQuery(ctx, ou, defaultTwitterSearchGroupName)
	s.Require().Nil(err)
	s.Assert().NotNil(ts)
	s.Assert().Equal(100, ts[0].MaxResults)
}

func (s *SearchAITestSuite) TestInsertRedditPosts() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	searchID := 1700685362343578112

	ts := time.Now()
	rts := &reddit.Timestamp{Time: ts}

	posts := []*reddit.Post{
		{
			ID:                    "1700514008821625088",
			FullID:                "1700514008821625088-1",
			Created:               rts,
			Edited:                nil,
			Permalink:             "/kub",
			URL:                   "lkdsd",
			Title:                 "title",
			Body:                  "body",
			Score:                 1,
			UpvoteRatio:           0.7,
			NumberOfComments:      13,
			SubredditName:         "",
			SubredditNamePrefixed: "",
			SubredditID:           "",
			SubredditSubscribers:  0,
			Author:                "zeus",
			AuthorID:              "zeus-1",
			Spoiler:               false,
			Locked:                false,
			NSFW:                  false,
			IsSelfPost:            false,
			Saved:                 false,
			Stickied:              false,
		},
	}
	resp, err := InsertIncomingRedditPosts(ctx, searchID, posts)
	s.Require().Nil(err)
	s.Assert().NotZero(resp)
}
