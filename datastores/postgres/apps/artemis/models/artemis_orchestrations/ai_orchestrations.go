package artemis_orchestrations

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowExecParams struct {
	CurrentCycleCount                           int                    `json:"currentCycleCount"`
	UnixStartTime                               int                    `json:"startTimeUnix"`
	AggNormalizedCycleCounts                    map[int]int            `json:"aggNormalizedCycleCounts"`
	TimeStepSize                                time.Duration          `json:"unixTimeStepSize"`
	TotalCyclesPerOneCompleteWorkflow           int                    `json:"totalCyclesPerOneCompleteWorkflow"`
	TotalCyclesPerOneCompleteWorkflowAsDuration time.Duration          `json:"totalCyclesPerOneCompleteWorkflowAsDuration"`
	WorkflowTemplate                            WorkflowTemplate       `json:"workflowTemplate"`
	WorkflowTasks                               []WorkflowTemplateData `json:"workflowTasks"`
}

func GetAiOrchestrationParams(ctx context.Context, ou org_users.OrgUser, unixStartTime int, wfs []WorkflowTemplate) ([]WorkflowExecParams, error) {
	var wfExecParams []WorkflowExecParams
	for _, wf := range wfs {
		wtd, err := SelectWorkflowTemplate(ctx, ou, wf.WorkflowName)
		if err != nil {
			return nil, err
		}
		wfTimeParams := AggregateTasks(wf, wtd)
		wfTimeParams.WorkflowTasks = wtd
		wfTimeParams.WorkflowTemplate = wf
		wfTimeParams.UnixStartTime = unixStartTime
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

func InsertOrchestrationRef(ctx context.Context, oj OrchestrationJob, b []byte) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO orchestrations(org_id, orchestration_name, group_name, type, instructions)
				  VALUES ($1, $2, $3, $4, $5)
				  ON CONFLICT (org_id, orchestration_name) 
				  DO UPDATE SET instructions = EXCLUDED.instructions
				  RETURNING orchestration_id;`

	var id int
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, oj.OrgID, oj.OrchestrationName, oj.GroupName, oj.Type, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&id)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return 0, err
	}
	return id, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}
