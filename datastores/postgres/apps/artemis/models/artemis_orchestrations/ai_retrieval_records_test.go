package artemis_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (s *OrchestrationsTestSuite) TestInsertWorkflowRetrievalResult() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	window := Window{
		UnixStartTime: 1701886760,
		UnixEndTime:   1701887760,
	}
	wr := &AIWorkflowRetrievalResult{
		OrchestrationID:       1701135755183363072,
		RetrievalID:           1701667813254964224, // Corrected field name
		IterationCount:        2,
		RunningCycleNumber:    1,
		SearchWindowUnixStart: window.UnixStartTime,
		SearchWindowUnixEnd:   window.UnixEndTime,
		SkipRetrieval:         false, // Corrected field name
		Metadata:              nil,
	}
	err := InsertWorkflowRetrievalResult(ctx, wr)
	s.Require().NoError(err)
	s.Require().NotZero(wr.WorkflowResultID)
}

func (s *OrchestrationsTestSuite) TestSelectRetrievalResultsIds() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	window := Window{
		UnixStartTime: 1701886760,
		UnixEndTime:   1701887760,
	}
	ojIds := []int{1701135755183363072}        // Example orchestration IDs
	sourceRetIds := []int{1701667813254964224} // Example retrieval IDs
	results, err := SelectRetrievalResultsIds(ctx, window, ojIds, sourceRetIds)
	s.Require().NoError(err)
	s.Require().NotEmpty(results)
}
