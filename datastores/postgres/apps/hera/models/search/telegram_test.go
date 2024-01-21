package hera_search

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *SearchAITestSuite) TestSelectTelegramResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	si := artemis_orchestrations.Window{}
	si.Start = time.Now().AddDate(0, 0, -60)
	si.End = time.Now()

	sp := AiSearchParams{
		Retrieval: artemis_orchestrations.RetrievalItem{
			//RetrievalID:    0,
			RetrievalName:  "",
			RetrievalGroup: "",
			RetrievalItemInstruction: artemis_orchestrations.RetrievalItemInstruction{
				RetrievalPlatform: "telegram",
				//RetrievalPrompt:         "",
				//RetrievalPlatformGroups: "Ze",
				//RetrievalKeywords:       "",
				//RetrievalUsernames:      "",
				DiscordFilters: nil,
			},
		},
		TimeRange: "",
		Window:    si,
	}
	res, err := SearchTelegram(ctx, ou, sp)
	s.Require().Nil(err)
	s.Assert().NotEmpty(res)

	ou.OrgID = 0
	res, err = SearchTelegram(ctx, ou, sp)
	s.Require().Nil(err)
	s.Assert().Nil(res)
}
