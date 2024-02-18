package artemis_orchestrations

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (s *OrchestrationsTestSuite) TestInsertWorkflowStageReference() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	wfStageIO := &WorkflowStageReference{
		InputID:       1,
		InputStrID:    "testInputStrID",
		WorkflowRunID: 1679548290001220864,
		ChildWfID:     "childWfID1",
		RunCycle:      1,
		InputData:     json.RawMessage(`{"key":"value"}`),
	}

	// Insert the mock data into the database
	err := InsertWorkflowStageReference(ctx, wfStageIO)
	s.Require().NoError(err)

	wr, err := SelectWorkflowStageReference(ctx, wfStageIO.InputID)
	s.Require().NoError(err)
	s.Require().NotNil(wr)
	s.Require().Equal(wfStageIO.InputID, wr.InputID)
	s.Require().Equal(json.RawMessage(`{"key":"value"}`), wfStageIO.InputData)

	wfStageIO.InputData = nil
	// Insert the mock data into the database
	err = InsertWorkflowStageReference(ctx, wfStageIO)
	s.Require().NoError(err)

	wr, err = SelectWorkflowStageReference(ctx, wfStageIO.InputID)
	s.Require().NoError(err)
	s.Require().NotNil(wr)
	s.Require().Equal(json.RawMessage(json.RawMessage(nil)), wfStageIO.InputData)
}
