package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (t *ZeusWorkerTestSuite) TestCsvIterator() {
	fnv := "test.json"
	dbg := OpenCsvIteratorDebug(fnv)
	na := NewZeusAiPlatformActivities()
	mb := dbg.Cp

	mb.Oj.OrchestrationID = 1715620241520326000
	mb.Window.UnixStartTime = 1715620241
	mb.Window.UnixEndTime = 1715634641
	mb.Tc.TaskID = 1711419223167298000
	mb.Tc.TaskName = "website-analysis"
	sv, serr := artemis_orchestrations.SelectAiWorkflowAnalysisResultsIds(ctx, mb.Window, []int{mb.Oj.OrchestrationID}, []int{mb.Tc.TaskID}, 0)
	t.Require().Nil(serr)
	fmt.Println(sv)
	sm := make(map[int]map[int]bool)
	for _, vi := range sv {
		if _, ok := sm[vi.ChunkOffset]; !ok {
			sm[vi.ChunkOffset] = make(map[int]bool)
		}
		sm[vi.ChunkOffset][vi.IterationCount] = true
	}
	fmt.Println(na)
	mb.Tc.ChunkIterator = 203
	mb.Tc.Model = "gpt-3.5-turbo-0125"
	err := na.CsvIterator(ctx, mb)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestCsvIteratorReports() {
	fnv := "CsvIteratorDebug-cycle-1-chunk-0-1714195427378636000.json"
	dbg := OpenCsvIteratorDebug(fnv)
	na := NewZeusAiPlatformActivities()
	mb := dbg.Cp
	sv, serr := artemis_orchestrations.SelectAiWorkflowAnalysisResultsIds(ctx, mb.Window, []int{mb.Oj.OrchestrationID}, []int{mb.Tc.TaskID}, 0)
	t.Require().Nil(serr)
	fmt.Println(sv)
	sm := make(map[int]map[int]bool)
	for _, vi := range sv {
		if _, ok := sm[vi.ChunkOffset]; !ok {
			sm[vi.ChunkOffset] = make(map[int]bool)
		}
		sm[vi.ChunkOffset][vi.IterationCount] = true
	}

	fmt.Println(na)
	err := na.GenerateCycleReports(ctx, mb)
	t.Require().Nil(err)
}

func (t *ZeusWorkerTestSuite) TestC() {
	//plms := []map[string]interface{}{
	//	{"entity": "https://www.homeenvironmentsolutions.com", "prompt1": "home insulation"},
	//	{"entity": "https://www.homeenvironmentsolutions.com", "prompt2": "Noticed on your website you provide energy-saving and air quality services to keep homes comfortable and healthy."},
	//	{"entity": "https://www.aol.com", "prompt1": "try me"},
	//	{"entity": "https://www.aol.com", "prompt2": "Noticed test."},
	//}
	//fmt.Println(res)
}
