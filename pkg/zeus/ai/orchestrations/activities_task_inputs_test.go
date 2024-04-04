package ai_platform_service_orchestrations

import "encoding/json"

/*
	func getAnalysisDeps(aggInst artemis_orchestrations.WorkflowTemplateData, wfExecParams artemis_orchestrations.WorkflowExecParams) []int {
		depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
		var analysisDep []int
		for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
			analysisDep = append(analysisDep, k)
		}
		return analysisDep
	}
	todo: needs to set ^ when adding csv merge stage
*/

func (t *ZeusWorkerTestSuite) TestAiAggregateAnalysisRetrievalTask() {
	dbg := AiAggregateAnalysisRetrievalTaskInputDebug{}
	fp := dbg.OpenFp()
	b := fp.ReadFileInPath()
	err := json.Unmarshal(b, &dbg)
	t.Require().Nil(err)
	na := NewZeusAiPlatformActivities()
	_, err = na.AiAggregateAnalysisRetrievalTask(ctx, dbg.Cp, nil)
	t.Require().Nil(err)
}
