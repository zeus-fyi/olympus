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
		RetrievalName:  "TestRetrieval",
		RetrievalGroup: "TestGroup1",
	}

	// Example instructions in byte array
	instructions := []byte(`{"key": "value"}`)

	// Step 2: Call InsertRetrieval
	err := InsertRetrieval(ctx, ou, &mockRetrievalItem, instructions)
	s.Require().Nil(err)
	s.Require().NotZero(mockRetrievalItem.RetrievalID)
}
