package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowExecParams struct {
	WorkflowTemplate              WorkflowTemplate              `json:"workflowTemplate"`
	WorkflowExecTimekeepingParams WorkflowExecTimekeepingParams `json:"workflowExecTimekeepingParams"`
	CycleCountTaskRelative        CycleCountTaskRelative        `json:"cycleCountTaskRelative"`
	WorkflowTaskRelationships     WorkflowTaskRelationships     `json:"workflowTaskRelationships"`
	WorkflowTasks                 []WorkflowTemplateData        `json:"workflowTasks"`
}

type WorkflowExecTimekeepingParams struct {
	CurrentCycleCount                           int           `json:"currentCycleCount"`
	RunWindow                                   Window        `json:"runWindow"`
	RunTimeDuration                             time.Duration `json:"runTimeDuration"`
	RunCycles                                   int           `json:"runCycles"`
	TimeStepSize                                time.Duration `json:"unixTimeStepSize"`
	TotalCyclesPerOneCompleteWorkflow           int           `json:"totalCyclesPerOneCompleteWorkflow"`
	TotalCyclesPerOneCompleteWorkflowAsDuration time.Duration `json:"totalCyclesPerOneCompleteWorkflowAsDuration"`
}

type WorkflowTaskRelationships struct {
	AnalysisTasks    map[int]AnalysisTaskDB    `json:"analysisTasks,omitempty"`
	AggAnalysisTasks map[int]map[int]AggTaskDb `json:"aggAnalysisTasks,omitempty"`

	AnalysisRetrievals map[int]map[int]bool `json:"analysisRetrievals"`
	AggregateAnalysis  map[int]map[int]bool `json:"aggregateAnalysis"`
}

type CycleCountTaskRelative struct {
	AggNormalizedCycleCounts             map[int]int                 `json:"aggNormalizedCycleCounts"`
	AnalysisEvalNormalizedCycleCounts    map[int]map[int]int         `json:"analysisEvalNormalizedCycleCounts,omitempty"`
	AggEvalNormalizedCycleCounts         map[int]map[int]int         `json:"aggEvalNormalizedCycleCounts,omitempty"`
	AggAnalysisEvalNormalizedCycleCounts map[int]map[int]map[int]int `json:"aggAnalysisEvalNormalizedCycleCounts,omitempty"`
}

type Window struct {
	Start         time.Time `json:"start,omitempty"`
	End           time.Time `json:"end,omitempty"`
	UnixStartTime int       `json:"unixStartTime,omitempty"`
	UnixEndTime   int       `json:"unixEndTime,omitempty"`
}

func (ti *Window) IsWindowEmpty() bool {
	if ti == nil {
		return true
	}
	return ti.Start.IsZero() && ti.End.IsZero() && ti.UnixStartTime == 0 && ti.UnixEndTime == 0
}

func (ti *Window) GetUnixTimestamps() (int, int) {
	if ti == nil {
		return 0, 0
	}
	if ti.UnixStartTime != 0 && ti.UnixEndTime != 0 {
		return ti.UnixStartTime, ti.UnixEndTime
	}
	return int(ti.Start.Unix()), int(ti.End.Unix())
}

// CalculateTimeWindow [0, 1] gives the time window for the first cycle
func CalculateTimeWindow(unixStartTime int, cycleStart, cycleEnd, timeStep time.Duration) Window {
	start := (cycleStart) * timeStep
	end := (cycleEnd) * timeStep
	wind := Window{
		Start: time.Unix(int64(unixStartTime), 0).Add(start),
		End:   time.Unix(int64(unixStartTime), 0).Add(end),
	}
	wind.UnixStartTime = int(wind.Start.Unix())
	wind.UnixEndTime = int(wind.End.Unix())
	return wind
}

func CalculateTimeWindowFromCycles(unixStartTime int, cycleStart, cycleEnd int, timeStep time.Duration) Window {
	start := time.Duration(cycleStart) * timeStep
	end := time.Duration(cycleEnd) * timeStep
	wind := Window{
		Start: time.Unix(int64(unixStartTime), 0).Add(start),
		End:   time.Unix(int64(unixStartTime), 0).Add(end),
	}
	wind.UnixStartTime = int(wind.Start.Unix())
	wind.UnixEndTime = int(wind.End.Unix())
	return wind
}

func GetAiOrchestrationParams(ctx context.Context, ou org_users.OrgUser, unixStartTime, unixEndTime int, wfs []WorkflowTemplate) ([]WorkflowExecParams, error) {
	var wfExecParams []WorkflowExecParams
	for _, wf := range wfs {
		wtd, err := SelectWorkflowTemplateByName(ctx, ou, wf.WorkflowName)
		if err != nil {
			log.Err(err).Msg("error selecting workflow template")
			return nil, err
		}
		if unixStartTime == 0 {
			unixStartTime = int(time.Now().Unix())
		}
		if wtd == nil || len(wtd.WorkflowTemplateSlice) <= 0 {
			return nil, errors.New("workflow template not found")
		}
		wfTimeParams := ConvertTemplateValuesToWorkflowTemplateData(wf, wtd.WorkflowTemplateSlice[0])
		if unixEndTime == 0 {
			unixEndTime = unixStartTime + int(wfTimeParams.WorkflowExecTimekeepingParams.TotalCyclesPerOneCompleteWorkflowAsDuration.Seconds())
		}

		wfTimeParams.WorkflowTemplate = wf
		wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime = unixStartTime
		wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.Start = time.Unix(int64(unixStartTime), 0)
		wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.UnixEndTime = unixEndTime
		wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.End = time.Unix(int64(unixEndTime), 0)
		wfTimeParams.WorkflowExecTimekeepingParams.RunTimeDuration = wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.End.Sub(wfTimeParams.WorkflowExecTimekeepingParams.RunWindow.Start)
		if wfTimeParams.WorkflowExecTimekeepingParams.TimeStepSize == 0 {
			return nil, errors.New("time step size is 0")
		}
		wfTimeParams.WorkflowExecTimekeepingParams.RunCycles = int(wfTimeParams.WorkflowExecTimekeepingParams.RunTimeDuration / wfTimeParams.WorkflowExecTimekeepingParams.TimeStepSize)
		wfExecParams = append(wfExecParams, wfTimeParams)
	}
	return wfExecParams, nil
}

func CalculateStepSizeUnix(stepSize int, stepUnit string) int {
	switch stepUnit {
	case "seconds":
		return stepSize
	case "minutes":
		return stepSize * 60
	case "hours", "hour":
		return stepSize * 60 * 60
	case "days":
		return stepSize * 60 * 60 * 24
	case "weeks":
		return stepSize * 60 * 60 * 24 * 7
	}
	return 0
}

func CalculateAggCycleCount(aggBaseCycleCount int, analysisCycleCounts int) int {
	if analysisCycleCounts > aggBaseCycleCount {
		aggBaseCycleCount = analysisCycleCounts
	}
	return aggBaseCycleCount
}

func ConvertTemplateValuesToWorkflowTemplateData(wf WorkflowTemplate, wfValue WorkflowTemplateValue) WorkflowExecParams {
	aggNormalizedCycleCount := make(map[int]int)
	aggNormalizedEvalsCycleCount := make(map[int]map[int]int)

	wd := make([]WorkflowTemplateData, 0)
	globalMaxCycleLength := 1
	for aggID, aggAnalysisTaskMap := range wfValue.AggAnalysisTasks {
		maxCycleLength := 1
		aggCycleLength := 1
		var ev []EvalFnDB
		for analysisID, aggAnalysisTask := range aggAnalysisTaskMap {
			analysisTask, ok := wfValue.AnalysisTasks[analysisID]
			if !ok {
				continue
			}
			wtd := WorkflowTemplateData{
				AnalysisTaskDB:           analysisTask,
				AnalysisMaxTokensPerTask: aggAnalysisTask.AggMaxTokensPerTask,
				AggTaskID:                &aggAnalysisTask.AggTaskId,
				AggCycleCount:            &aggAnalysisTask.AggCycleCount,
				AggTaskName:              &aggAnalysisTask.AggTaskName,
				AggTaskType:              &aggAnalysisTask.AggTaskType,
				AggPrompt:                &aggAnalysisTask.AggPrompt,
				AggModel:                 &aggAnalysisTask.AggModel,
				AggTokenOverflowStrategy: &aggAnalysisTask.AggTokenOverflowStrategy,
				AggMaxTokensPerTask:      &aggAnalysisTask.AggMaxTokensPerTask,
				AggEvalFns:               aggAnalysisTask.EvalFns,
				AggAnalysisEvalFns:       aggAnalysisTask.AnalysisAggEvalFns,
			}
			ev = append(ev, aggAnalysisTask.EvalFns...)
			if aggAnalysisTask.AggCycleCount > aggCycleLength {
				aggCycleLength = aggAnalysisTask.AggCycleCount
			}
			normalizedCount := CalculateAggCycleCount(aggAnalysisTask.AggCycleCount, analysisTask.AnalysisCycleCount)
			if normalizedCount >= maxCycleLength {
				maxCycleLength = normalizedCount
			}
			if normalizedCount >= globalMaxCycleLength {
				globalMaxCycleLength = normalizedCount
			}
			wd = append(wd, wtd)
		}
		normalizedCount := maxCycleLength * aggCycleLength
		if normalizedCount >= globalMaxCycleLength {
			globalMaxCycleLength = normalizedCount
		}
		aggNormalizedCycleCount[aggID] = normalizedCount
		for _, evalFn := range ev {
			if _, ok := aggNormalizedEvalsCycleCount[aggID]; !ok {
				aggNormalizedEvalsCycleCount[aggID] = make(map[int]int)
			}
			aggNormalizedEvalsCycleCount[aggID][evalFn.EvalID] = evalFn.EvalCycleCount * aggNormalizedCycleCount[aggID]
		}
	}

	aggAnalysisEvalNormalizedCycleCounts := make(map[int]map[int]map[int]int)
	for _, aggAnalysisTaskMap := range wfValue.AggAnalysisTasks {
		for analysisID, aggAnalysisTask := range aggAnalysisTaskMap {
			if _, aok := aggAnalysisEvalNormalizedCycleCounts[aggAnalysisTask.AggTaskId]; !aok {
				aggAnalysisEvalNormalizedCycleCounts[aggAnalysisTask.AggTaskId] = make(map[int]map[int]int)
			}
			for _, evalFn := range aggAnalysisTask.AnalysisAggEvalFns {
				if _, aok := aggAnalysisEvalNormalizedCycleCounts[aggAnalysisTask.AggTaskId][analysisID]; !aok {
					aggAnalysisEvalNormalizedCycleCounts[aggAnalysisTask.AggTaskId][analysisID] = make(map[int]int)
				}
				aggAnalysisEvalNormalizedCycleCounts[aggAnalysisTask.AggTaskId][analysisID][evalFn.EvalID] =
					evalFn.EvalCycleCount * aggAnalysisTask.AnalysisAggCycleCount
			}
		}
	}

	analysisEvalNormalizedCycles := make(map[int]map[int]int)
	for _, analysisTask := range wfValue.AnalysisTasksSlice {
		for _, evalFn := range analysisTask.AnalysisEvalFns {
			if evalFn.EvalCycleCount == 0 {
				evalFn.EvalCycleCount = 1
			}
			if analysisTask.AnalysisCycleCount == 0 {
				analysisTask.AnalysisCycleCount = 1
			}
			if _, ok := analysisEvalNormalizedCycles[analysisTask.AnalysisTaskID]; !ok {
				analysisEvalNormalizedCycles[analysisTask.AnalysisTaskID] = make(map[int]int)
			}
			analysisEvalNormalizedCycles[analysisTask.AnalysisTaskID][evalFn.EvalID] = evalFn.EvalCycleCount * analysisTask.AnalysisCycleCount
		}
	}
	wte := WorkflowExecParams{
		WorkflowTemplate: wf,
		WorkflowExecTimekeepingParams: WorkflowExecTimekeepingParams{
			TimeStepSize:                                time.Duration(CalculateStepSizeUnix(wf.FundamentalPeriod, wf.FundamentalPeriodTimeUnit)) * time.Second,
			TotalCyclesPerOneCompleteWorkflow:           globalMaxCycleLength,
			TotalCyclesPerOneCompleteWorkflowAsDuration: time.Duration(CalculateStepSizeUnix(wf.FundamentalPeriod, wf.FundamentalPeriodTimeUnit)*globalMaxCycleLength) * time.Second,
		},
		CycleCountTaskRelative: CycleCountTaskRelative{
			AggNormalizedCycleCounts:             aggNormalizedCycleCount,
			AnalysisEvalNormalizedCycleCounts:    analysisEvalNormalizedCycles,
			AggEvalNormalizedCycleCounts:         aggNormalizedEvalsCycleCount,
			AggAnalysisEvalNormalizedCycleCounts: aggAnalysisEvalNormalizedCycleCounts,
		},
		WorkflowTaskRelationships: WorkflowTaskRelationships{
			AggAnalysisTasks: wfValue.AggAnalysisTasks,
			AnalysisTasks:    wfValue.AnalysisTasks,
		},
		WorkflowTasks: wd,
	}
	return wte
}

func UpsertAiOrchestration(ctx context.Context, ou org_users.OrgUser, wfParentID string, wfExec WorkflowExecParams) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO orchestrations(org_id, orchestration_name, group_name, type, active, instructions)
              VALUES ($1, $2, $3, $4, $5, $6)
              ON CONFLICT (org_id, orchestration_name) 
			  DO UPDATE SET 
				  instructions = EXCLUDED.instructions,
				  active = EXCLUDED.active
			  RETURNING orchestration_id;`

	var id int
	b, err := json.Marshal(wfExec)
	if err != nil {
		log.Err(err).Msg("error marshalling workflow execution params")
		return 0, err
	}
	active := false
	tn := time.Now().Unix()
	if wfExec.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime == 0 || wfExec.WorkflowExecTimekeepingParams.RunWindow.UnixStartTime >= int(tn) {
		active = true
	}
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, wfParentID, wfExec.WorkflowTemplate.WorkflowGroup, wfExec.WorkflowTemplate.WorkflowName, active, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&id)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return 0, err
	}
	return id, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}
