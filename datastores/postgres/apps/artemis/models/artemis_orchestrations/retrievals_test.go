package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertRetrieval() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	mockRetrievalItem := RetrievalItem{
		RetrievalName:  "r-2",
		RetrievalGroup: "r-2",
		RetrievalItemInstruction: RetrievalItemInstruction{
			RetrievalPlatform:       "reddit",
			RetrievalPrompt:         "",
			RetrievalPlatformGroups: "",
			RetrievalKeywords:       "",
		},
		Instructions: []byte(`{"key": "value"}`),
	}

	// Step 2: Call InsertRetrieval
	err := InsertRetrieval(ctx, ou, &mockRetrievalItem)
	s.Require().Nil(err)
	s.Require().NotZero(mockRetrievalItem.RetrievalID)
}

func (s *OrchestrationsTestSuite) TestSelectRetrievals() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Step 2: Call InsertRetrieval
	res, err := SelectRetrievals(ctx, ou)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
}
