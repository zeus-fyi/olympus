package ai_platform_service_orchestrations

import (
	"encoding/json"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (t *ZeusWorkerTestSuite) TestGenerateCycleReports() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	act := NewZeusAiPlatformActivities()

	dbg := AiAggregateAnalysisRetrievalTaskInputDebug{}
	fp := dbg.OpenFp()
	bv := fp.ReadFileInPath()
	err := json.Unmarshal(bv, &dbg)
	t.Require().Nil(err)
	t.Require().NotNil(dbg.Cp)
	err = act.GenerateCycleReports(ctx, dbg.Cp)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestReportExport() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	fnv := "CsvIteratorDebug-cycle-1-chunk-0-1714614239098819000.json"
	dbg := OpenCsvIteratorDebug(fnv)

	gens, err := GetGlobalEntitiesFromRef(ctx, dbg.Cp.Ou, dbg.Cp.WfExecParams.WorkflowOverrides.WorkflowEntityRefs)
	t.Require().Nil(err)
	t.Require().NotNil(gens)
}
