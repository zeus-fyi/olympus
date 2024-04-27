package ai_platform_service_orchestrations

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
)

func (t *ZeusWorkerTestSuite) TestCsvIterator() {
	fnv := "CsvIteratorDebug-cycle-1-chunk-0-1712604784507256000.json"
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
	plms := []map[string]interface{}{
		{"entity": "https://www.homeenvironmentsolutions.com", "prompt1": "home insulation"},
		{"entity": "https://www.homeenvironmentsolutions.com", "prompt2": "Noticed on your website you provide energy-saving and air quality services to keep homes comfortable and healthy."},
		{"entity": "https://www.aol.com", "prompt1": "try me"},
		{"entity": "https://www.aol.com", "prompt2": "Noticed test."},
	}
	res := convEntityToCsvCol("Website", plms, 0)
	fmt.Println(res)
}
