package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowExecParams struct {
	CurrentCycleCount                           int                    `json:"currentCycleCount"`
	RunWindow                                   Window                 `json:"runWindow"`
	RunTimeDuration                             time.Duration          `json:"runTimeDuration"`
	RunCycles                                   int                    `json:"runCycles"`
	AggNormalizedCycleCounts                    map[int]int            `json:"aggNormalizedCycleCounts"`
	TimeStepSize                                time.Duration          `json:"unixTimeStepSize"`
	TotalCyclesPerOneCompleteWorkflow           int                    `json:"totalCyclesPerOneCompleteWorkflow"`
	TotalCyclesPerOneCompleteWorkflowAsDuration time.Duration          `json:"totalCyclesPerOneCompleteWorkflowAsDuration"`
	WorkflowTemplate                            WorkflowTemplate       `json:"workflowTemplate"`
	WorkflowTasks                               []WorkflowTemplateData `json:"workflowTasks"`
}

type Window struct {
	Start         time.Time `json:"start"`
	End           time.Time `json:"end"`
	UnixStartTime int       `json:"unixStartTime"`
	UnixEndTime   int       `json:"unixEndTime"`
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
		wtd, err := SelectWorkflowTemplate(ctx, ou, wf.WorkflowName)
		if err != nil {
			log.Err(err).Msg("error selecting workflow template")
			return nil, err
		}
		if unixStartTime == 0 {
			unixStartTime = int(time.Now().Unix())
		}
		wfTimeParams := AggregateTasks(wf, wtd)
		if unixEndTime == 0 {
			unixEndTime = unixStartTime + int(wfTimeParams.TotalCyclesPerOneCompleteWorkflowAsDuration.Seconds())
		}

		wfTimeParams.WorkflowTasks = wtd
		wfTimeParams.WorkflowTemplate = wf
		wfTimeParams.RunWindow.UnixStartTime = unixStartTime
		wfTimeParams.RunWindow.Start = time.Unix(int64(unixStartTime), 0)
		wfTimeParams.RunWindow.UnixEndTime = unixEndTime
		wfTimeParams.RunWindow.End = time.Unix(int64(unixEndTime), 0)
		wfTimeParams.RunTimeDuration = wfTimeParams.RunWindow.End.Sub(wfTimeParams.RunWindow.Start)
		wfTimeParams.RunCycles = int(wfTimeParams.RunTimeDuration / wfTimeParams.TimeStepSize)
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
	case "days":
		return stepSize * 60 * 24
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

func AggregateTasks(wf WorkflowTemplate, wd []WorkflowTemplateData) WorkflowExecParams {
	aggMap := make(map[int]int)
	aggNormalizedCycleCount := make(map[int]int)

	maxCycleLength := 1
	for _, w := range wd {
		if w.AnalysisCycleCount > maxCycleLength {
			maxCycleLength = w.AnalysisCycleCount
		}
		if w.AggTaskID != nil && w.AggCycleCount != nil {
			aggVal := aggMap[*w.AggTaskID]
			aggNormalizedCycleCount[*w.AggTaskID] = *w.AggCycleCount
			aggMap[*w.AggTaskID] = CalculateAggCycleCount(aggVal, w.AnalysisCycleCount)
		}
	}
	for k, v := range aggMap {
		aggNormalizedCycleCount[k] = v * aggNormalizedCycleCount[k]
		if aggNormalizedCycleCount[k] > maxCycleLength {
			maxCycleLength = aggNormalizedCycleCount[k]
		}
	}
	return WorkflowExecParams{
		AggNormalizedCycleCounts:                    aggNormalizedCycleCount,
		TotalCyclesPerOneCompleteWorkflow:           maxCycleLength,
		TimeStepSize:                                time.Duration(CalculateStepSizeUnix(wf.FundamentalPeriod, wf.FundamentalPeriodTimeUnit)) * time.Second,
		TotalCyclesPerOneCompleteWorkflowAsDuration: time.Duration(CalculateStepSizeUnix(wf.FundamentalPeriod, wf.FundamentalPeriodTimeUnit)*maxCycleLength) * time.Second,
	}
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
	if wfExec.RunWindow.UnixStartTime == 0 || wfExec.RunWindow.UnixStartTime >= int(tn) {
		active = true
	}
	err = apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, wfParentID, wfExec.WorkflowTemplate.WorkflowGroup, wfExec.WorkflowTemplate.WorkflowName, active, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&id)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return 0, err
	}
	return id, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}
