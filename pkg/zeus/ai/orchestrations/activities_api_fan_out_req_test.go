package ai_platform_service_orchestrations

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

func (t *ZeusWorkerTestSuite) TestFanOutApiCallRequestTask() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	na := NewZeusAiPlatformActivities()
	fnv := "CsvIteratorDebug-cycle-1-chunk-0-1712604784507256000.json"
	dbg := OpenCsvIteratorDebug(fnv)

	rts, err := na.AiWebRetrievalGetRoutesTask(ctx, getOrgRetIfFlows(dbg.Cp), dbg.Cp.Tc.Retrieval)
	t.Require().Nil(err)
	t.Require().NotNil(rts)

	_, err = na.FanOutApiCallRequestTask(ctx, rts, dbg.Cp)
	t.Require().Nil(err)
}
