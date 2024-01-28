package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
)

type AIWorkflowEvalResultResponse struct {
	EvalResultsID      int `db:"eval_result_id" json:"evalResultID"`
	EvalID             int `db:"eval_id" json:"evalID"`
	WorkflowResultID   int `db:"workflow_result_id" json:"workflowResultID"`
	ResponseID         int `db:"response_id" json:"responseID"`
	EvalIterationCount int `db:"eval_iteration_count" json:"evalIterationCount"`
}

const Sn = "UpsertEvalMetricsResults"

func UpsertEvalMetricsResults(ctx context.Context, emrs *EvalMetricsResults) error {
	fp := filepaths.Path{
		DirOut: "/Users/alex/go/Olympus/olympus/datastores/postgres/apps/artemis/models/artemis_orchestrations",
		FnOut:  "debug-1.json",
	}

	b, _ := json.Marshal(emrs)
	fmt.Println(string(b))
	err := fp.WriteToFileOutPath(b)
	if err != nil {
		log.Err(err).Msg("UpsertEvalMetricsResults: failed to write to file")
		return err
	}

	q := sql_query_templates.NewQueryParam("UpsertEvalMetricsResults", "eval_metrics_results", "where", 1000, []string{})
	cte := sql_query_templates.CTE{Name: "UpsertEvalMetricsResults"}
	cte.SubCTEs = []sql_query_templates.SubCTE{}
	cte.Params = []interface{}{}
	ch := chronos.Chronos{}
	for _, emr := range emrs.EvalMetricsResults {
		queryResourceId := fmt.Sprintf("eval_metrics_results_insert_%d", ch.UnixTimeStampNow())
		scteRe := sql_query_templates.NewSubInsertCTE(queryResourceId)
		scteRe.TableName = "eval_metrics_results"
		scteRe.Columns = []string{
			"eval_metrics_results_id",
			"orchestration_id",
			"source_task_id",
			"eval_metric_id",
			"running_cycle_number",
			"search_window_unix_start",
			"search_window_unix_end",
			"eval_result_outcome",
			"eval_metadata",
			"chunk_offset",
			"eval_iteration_count",
		}
		if emr == nil || aws.ToInt(emr.EvalMetricID) == 0 || emr.EvalMetricResult == nil || emr.EvalMetricResult.EvalResultOutcomeBool == nil {
			continue
		}
		ts := chronos.Chronos{}
		tsNow := ts.UnixTimeStampNow()
		if aws.ToInt(emr.EvalMetricResult.EvalMetricResultID) <= 0 {
			emr.EvalMetricResult.EvalMetricResultID = &tsNow
			emr.EvalMetricResult.EvalMetricResultStrID = aws.String(fmt.Sprintf("%d", tsNow))
		}
		var pgTemp *pgtype.JSONB
		if emr.EvalMetricResult == nil {
			pgTemp = &pgtype.JSONB{Bytes: []byte{}, Status: IsNull(nil)}
		} else {
			pgTemp = &pgtype.JSONB{Bytes: sanitizeBytesUTF8(emr.EvalMetricResult.EvalMetadata), Status: IsNull(emr.EvalMetricResult.EvalMetadata)}
		}
		rv := apps.RowValues{
			aws.ToInt(emr.EvalMetricResult.EvalMetricResultID),
			emrs.EvalContext.AIWorkflowAnalysisResult.OrchestrationID,
			emrs.EvalContext.AIWorkflowAnalysisResult.SourceTaskID,
			aws.ToInt(emr.EvalMetricID),
			emrs.EvalContext.AIWorkflowAnalysisResult.RunningCycleNumber,
			emrs.EvalContext.AIWorkflowAnalysisResult.SearchWindowUnixStart,
			emrs.EvalContext.AIWorkflowAnalysisResult.SearchWindowUnixEnd,
			aws.ToBool(emr.EvalMetricResult.EvalResultOutcomeBool),
			pgTemp,
			emrs.EvalContext.AIWorkflowAnalysisResult.ChunkOffset,
			emrs.EvalContext.EvalIterationCount,
		}
		scteRe.Values = []apps.RowValues{rv}

		cte.SubCTEs = append(cte.SubCTEs, scteRe)
	}
	cte.OnConflicts = []string{"eval_metrics_results_id"}
	cte.OnConflictsUpdateColumns = []string{"running_cycle_number", "search_window_unix_start", "search_window_unix_end", "eval_result_outcome", "eval_metadata"}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("MetricResults: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertOrUpdateAiWorkflowEvalResultResponse(ctx context.Context, errr AIWorkflowEvalResultResponse) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO eval_results_responses(eval_results_id, workflow_result_id, eval_id, response_id, eval_iteration_count)
                  VALUES ($1, $2, $3, $4, $5)
                  ON CONFLICT (workflow_result_id, eval_id, eval_iteration_count, response_id)
                  DO UPDATE SET 
                      workflow_result_id = EXCLUDED.workflow_result_id,
                      eval_id = EXCLUDED.eval_id,
                      response_id = EXCLUDED.response_id
                  RETURNING eval_results_id;`

	if errr.EvalResultsID <= 0 {
		ch := chronos.Chronos{}
		errr.EvalResultsID = ch.UnixTimeStampNow()
	}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, errr.EvalResultsID, errr.WorkflowResultID, errr.EvalID, errr.ResponseID, errr.EvalIterationCount).Scan(&errr.EvalResultsID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("AIWorkflowEvalResultResponse")); returnErr != nil {
		log.Err(returnErr).Interface("errr", errr).Msg(q.LogHeader("AIWorkflowEvalResultResponse"))
		return errr.EvalResultsID, err
	}
	return errr.EvalResultsID, nil
}
