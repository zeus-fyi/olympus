package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestSelectWorkflowTemplate() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	newTemplate := WorkflowTemplate{
		WorkflowName:              "Example Workflow1",
		FundamentalPeriod:         5,
		WorkflowGroup:             "TestGroup1",
		FundamentalPeriodTimeUnit: "days",
	}

	res, err := SelectWorkflowTemplate(ctx, ou, newTemplate.WorkflowName)
	s.Require().Nil(err)
	s.Require().NotEmpty(res)
}
