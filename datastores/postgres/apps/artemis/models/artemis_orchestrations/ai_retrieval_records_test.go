package artemis_orchestrations

import "github.com/zeus-fyi/olympus/datastores/postgres/apps"

func (s *OrchestrationsTestSuite) TestInsertWorkflowRetrievalResult() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	window := Window{
		UnixStartTime: 1712702165,
		UnixEndTime:   1712716565,
	}
	wr := &AIWorkflowRetrievalResult{
		OrchestrationID:       1712702165698519000,
		RetrievalID:           1712533371223555000,
		ChunkOffset:           0,
		IterationCount:        0,
		RunningCycleNumber:    1,
		Status:                "complete",
		SearchWindowUnixStart: window.UnixStartTime,
		SearchWindowUnixEnd:   window.UnixEndTime,
		SkipRetrieval:         false,
		Metadata:              nil,
	}
	err := InsertWorkflowRetrievalResult(ctx, wr)
	s.Require().NoError(err)
	s.Require().NotZero(wr.WorkflowResultID)
}

func (s *OrchestrationsTestSuite) TestInsertWorkflowRetrievalResultError() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	window := Window{
		UnixStartTime: 1712702165,
		UnixEndTime:   1712716565,
	}
	wr := &AIWorkflowRetrievalResult{
		OrchestrationID:       1712702165698519000,
		RetrievalID:           1712533371223555000,
		ChunkOffset:           0,
		IterationCount:        0,
		Attempts:              0,
		RunningCycleNumber:    1,
		Status:                "error",
		SearchWindowUnixStart: window.UnixStartTime,
		SearchWindowUnixEnd:   window.UnixEndTime,
		SkipRetrieval:         false,
		Metadata:              nil,
	}
	err := InsertWorkflowRetrievalResultError(ctx, wr)
	s.Require().NoError(err)
	s.Require().NotZero(wr.WorkflowResultID)
}

func (s *OrchestrationsTestSuite) TestSelectRetrievalResultsIdsErrs() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	window := Window{
		UnixStartTime: 1701886760,
		UnixEndTime:   1701887760,
	}
	ojIds := []int{1712702165698519000}        // Example orchestration IDs
	sourceRetIds := []int{1712533371223555000} // Example retrieval IDs
	results, err := SelectRetrievalResultsIds(ctx, window, ojIds, sourceRetIds)
	s.Require().NoError(err)
	s.Require().NotEmpty(results)
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
