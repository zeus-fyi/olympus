package artemis_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (s *OrchestrationsTestSuite) TestSelectEvalMetricResults() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	//zz := AIWorkflowEvalResultResponse{
	//	EvalID:             1705978298687209000,
	//	WorkflowResultID:   1705978298687209000,
	//	ResponseID:         1672188679693780000,
	//	EvalIterationCount: 0,
	//}
	_, err := SelectEvalMetricResults(ctx, s.Ou)
	s.Require().Nil(err)
}
