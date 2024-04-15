package artemis_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

func (s *OrchestrationsTestSuite) TestArchiveRuns() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	err := UpdateOrchestrationsToArchive(ctx, ou, []string{"ai-workflow"}, true)
	s.Require().Nil(err)
}

func (s *OrchestrationsTestSuite) TestSelectRuns() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrations(ctx, ou, 0)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}

func (t *OrchestrationsTestSuite) TestSaveWfStatus() {
	//fnv := "CsvIteratorDebug-cycle-1-chunk-0-1712604784507256000.json"
	//dbg := OpenCsvIteratorDebug(fnv)
	//wfs := WfStatus{
	//	TotalApiRequests:    10,
	//	CompleteApiRequests: 9,
	//	TotalCsvElements:    100,
	//	CompleteCsvElements: 99,
	//}
	//err := SaveWfStatus(ctx, dbg.Cp, wfs)
	//t.Require().Nil(err)
	//wfrs, err := GetWfStatus(ctx, dbg.Cp)
	//t.Require().Nil(err)
	//t.Require().NotNil(wfrs)
}

func (s *OrchestrationsTestSuite) TestSelectRun() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = 1710298581127603000
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrations(ctx, ou, 1713167856080170000)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)

}

func (s *OrchestrationsTestSuite) TestSelectRunsUI() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = 1710298581127603000
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrationsUI(ctx, ou, 0)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}

func (s *OrchestrationsTestSuite) TestSelectRunWithRet() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	ou := org_users.OrgUser{}
	ou.OrgID = 1685378241971196000
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	ojs, err := SelectAiSystemOrchestrations(ctx, ou, 1713167856080170000)
	s.Require().Nil(err)
	s.Assert().NotEmpty(ojs)
}
