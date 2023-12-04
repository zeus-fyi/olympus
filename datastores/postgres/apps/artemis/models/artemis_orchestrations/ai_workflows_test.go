package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertAiWorkflow() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup1",
		FundamentalPeriodTimeUnit: "days",
	}

	err := InsertWorkflowTemplate(ctx, ou, &newTemplate)
	s.Require().Nil(err)
	s.Require().NotZero(newTemplate.WorkflowTemplateID)
}
