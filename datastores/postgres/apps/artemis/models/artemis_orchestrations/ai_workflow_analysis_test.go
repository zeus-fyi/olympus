package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestInsertAiWorkflowAnalysisResult() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	wr := &AIWorkflowAnalysisResult{
		OrchestrationsID:      1701135755183363072,
		ResponseID:            1693500147586158080,
		SourceTaskID:          1701657830780669952,
		IterationCount:        2,
		RunningCycleNumber:    1,
		SearchWindowUnixStart: 1701886760,
		SearchWindowUnixEnd:   1701887760,
		SkipAnalysis:          false,
		Metadata:              nil,
		CompletionChoices:     nil,
	}
	err := InsertAiWorkflowAnalysisResult(ctx, wr)
	s.Require().Nil(err)
	s.Require().NotZero(wr.WorkflowResultID)
}
