package ai_platform_service_orchestrations

import (
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations"
	hera_search "github.com/zeus-fyi/olympus/datastores/postgres/apps/hera/models/search"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type InputDataAnalysisToAgg struct {
	TextInput                   *string                        `json:"textInput,omitempty"`
	ChatCompletionQueryResponse *ChatCompletionQueryResponse   `json:"chatCompletionQueryResponse,omitempty"`
	SearchResultGroup           *hera_search.SearchResultGroup `json:"baseSearchResultsGroup,omitempty"`
}

func getDefaultRetryPolicy() workflow.ActivityOptions {
	ao := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 24, // Setting a valid non-zero timeout
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second * 5,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute * 5,
			MaximumAttempts:    1000,
		},
	}
	return ao
}

func getWr(cp *MbChildSubProcessParams, chunkOffset int) *artemis_orchestrations.AIWorkflowAnalysisResult {
	wr := &artemis_orchestrations.AIWorkflowAnalysisResult{
		OrchestrationID:       cp.Oj.OrchestrationID,
		SourceTaskID:          cp.Tc.TaskID,
		IterationCount:        0,
		ChunkOffset:           chunkOffset,
		RunningCycleNumber:    cp.Wsr.RunCycle,
		SearchWindowUnixStart: cp.Window.UnixStartTime,
		SearchWindowUnixEnd:   cp.Window.UnixEndTime,
		ResponseID:            cp.Tc.ResponseID,
	}
	return wr
}

const (
	aggAllCsvTaskID = 100
)

func getAnalysisDeps(aggInst artemis_orchestrations.WorkflowTemplateData, wfExecParams artemis_orchestrations.WorkflowExecParams) []int {
	var analysisDep []int
	depM := artemis_orchestrations.MapDependencies(wfExecParams.WorkflowTasks)
	for k, _ := range depM.AggregateAnalysis[*aggInst.AggTaskID] {
		analysisDep = append(analysisDep, k)
	}
	return analysisDep
}

func isInvalidAggInst(aggInst artemis_orchestrations.WorkflowTemplateData, md artemis_orchestrations.WorkflowTaskRelationships, wfExecParams artemis_orchestrations.WorkflowExecParams) bool {
	if aggInst.AggTaskID == nil || aggInst.AggCycleCount == nil || aggInst.AggModel == nil || aggInst.AggTaskName == nil {
		return true
	}
	if md.AggregateAnalysis[*aggInst.AggTaskID] == nil || md.AggregateAnalysis[*aggInst.AggTaskID][aggInst.AnalysisTaskID] == false || wfExecParams.WorkflowTaskRelationships.AggAnalysisTasks == nil {
		return true
	}
	return false
}
