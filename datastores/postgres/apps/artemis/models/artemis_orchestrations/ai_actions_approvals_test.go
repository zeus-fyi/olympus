package artemis_orchestrations

import (
	"github.com/labstack/echo/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestTriggerActionsApprovals() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	tap := TriggerActionsApproval{
		EvalID:           1704066747085827000,
		TriggerID:        1706487755984811000,
		WorkflowResultID: 1705953610460614000,
		ApprovalState:    "pending",
	}

	// Call the function to test
	err := CreateOrUpdateTriggerActionApproval(ctx, ou, tap)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestSelectTriggerActionApproval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ap, err := SelectTriggerActionApproval(ctx, s.Ou, "finished", 1706641322204777000)
	s.Require().Nil(err)
	s.Require().NotNil(ap)
}

func (s *OrchestrationsTestSuite) TestCreateOrUpdateTriggerActionApprovalWithApiReq() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	wtr := AIWorkflowTriggerResultApiReqResponse{
		ResponseID:  1706566091973014000,
		RetrievalID: 1706487709357339000,
		ReqPayloads: []echo.Map{
			{
				"key": "value1",
			},
			{
				"key": "value2",
			},
		},
	}

	tap := TriggerActionsApproval{
		EvalID:           1704066747085827000,
		TriggerID:        1706487755984811000,
		WorkflowResultID: 1705953610460614000,
		ApprovalID:       1706566091973007000,
		ApprovalState:    "pending",
		RequestSummary:   "pending-summary",
	}
	err := CreateOrUpdateTriggerActionApprovalWithApiReq(ctx, s.Ou, tap, wtr)
	s.Require().Nil(err)
}
func (s *OrchestrationsTestSuite) TestSelectTriggerActionApprovalWithApiReqResp() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	ap, err := SelectTriggerActionApprovalWithReqResponses(ctx, s.Ou, "pending", 1706566091973007000, 0)
	s.Require().Nil(err)
	s.Require().NotNil(ap)
}
