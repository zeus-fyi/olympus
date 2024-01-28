package artemis_orchestrations

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestTriggerInserts() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	triggerAction := TriggerAction{
		OrgID:        ou.OrgID,
		UserID:       ou.UserID,
		TriggerName:  "TestTrigger",
		TriggerGroup: "TestGroup",
		EvalTriggerActions: []EvalTriggerActions{
			{
				EvalID:               1703922045959259000, // Use an appropriate EvalID
				EvalTriggerState:     "info",
				EvalResultsTriggerOn: "all-pass",
			},
		},
	}

	// Call the function to test
	err := CreateOrUpdateTriggerAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, ou, 0)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}

func (s *OrchestrationsTestSuite) TestCreateTriggerApiRetrieval() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	// Create a TriggerActions instance
	triggerAction := TriggerAction{
		TriggerStrID:  "",
		TriggerID:     0,
		OrgID:         ou.OrgID,
		UserID:        ou.UserID,
		TriggerName:   "TestTrigger",
		TriggerGroup:  "TestGroup",
		TriggerAction: "",
		TriggerRetrievals: []RetrievalItem{
			{
				RetrievalStrID: aws.String("1703922045959259000"),
				RetrievalID:    aws.Int(1703922045959259000),
				RetrievalName:  "test",
				RetrievalGroup: "test",
				RetrievalItemInstruction: RetrievalItemInstruction{
					RetrievalPlatform:         "",
					RetrievalPrompt:           nil,
					RetrievalPlatformGroups:   nil,
					RetrievalKeywords:         nil,
					RetrievalNegativeKeywords: nil,
					RetrievalUsernames:        nil,
					DiscordFilters:            nil,
					WebFilters:                nil,
					Instructions:              nil,
				},
			},
		},
		TriggerPlatformReference: TriggerPlatformReference{},
		EvalTriggerAction:        EvalTriggerActions{},
		EvalTriggerActions: []EvalTriggerActions{
			{
				EvalID:               1703922045959259000, // Use an appropriate EvalID
				EvalTriggerState:     "info",
				EvalResultsTriggerOn: "all-pass",
			},
		},
		TriggerActionsApprovals: nil,
	}

	// Call the function to test
	err := CreateOrUpdateTriggerAction(ctx, ou, &triggerAction)
	s.Require().Nil(err)

	res2, err := SelectTriggerActionsByOrgAndOptParams(ctx, ou, 0)
	s.Require().Nil(err)
	s.Require().NotNil(res2)
}
