package artemis_orchestrations

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AIWorkflowEvalResultResponse struct {
	EvalResultsID      int `db:"eval_result_id" json:"evalResultID"`
	EvalID             int `db:"eval_id" json:"evalID"`
	WorkflowResultID   int `db:"workflow_result_id" json:"workflowResultID"`
	ResponseID         int `db:"response_id" json:"responseID"`
	EvalIterationCount int `db:"eval_iteration_count" json:"evalIterationCount"`
}

func UpsertEvalMetricsResults(ctx context.Context, emrs *EvalMetricsResults) error {
	if emrs == nil || emrs.EvalMetricsResults == nil {
		log.Info().Msg("UpsertEvalMetricsResults: emr is nil")
		return nil
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	const query = `
       INSERT INTO public.eval_metrics_results (
           eval_metrics_results_id,
           orchestration_id,
           source_task_id,
           eval_metric_id,
           running_cycle_number,
           search_window_unix_start,
           search_window_unix_end,
           eval_result_outcome,
           eval_metadata,
           chunk_offset,
           eval_iteration_count
       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
       ON CONFLICT (eval_metrics_results_id, eval_metric_id, source_task_id, orchestration_id, running_cycle_number, chunk_offset, eval_iteration_count)
       DO UPDATE SET
           running_cycle_number = EXCLUDED.running_cycle_number,
           search_window_unix_start = EXCLUDED.search_window_unix_start,
           search_window_unix_end = EXCLUDED.search_window_unix_end,
           eval_result_outcome = EXCLUDED.eval_result_outcome,
           eval_metadata = EXCLUDED.eval_metadata;
   `
	for _, emr := range emrs.EvalMetricsResults {
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
		_, err = tx.Exec(ctx, query,
			emr.EvalMetricResult.EvalMetricResultID,
			emrs.EvalContext.AIWorkflowAnalysisResult.OrchestrationID,
			emrs.EvalContext.AIWorkflowAnalysisResult.SourceTaskID,
			emr.EvalMetricID,
			emrs.EvalContext.AIWorkflowAnalysisResult.RunningCycleNumber,
			emrs.EvalContext.AIWorkflowAnalysisResult.SearchWindowUnixStart,
			emrs.EvalContext.AIWorkflowAnalysisResult.SearchWindowUnixEnd,
			emr.EvalMetricResult.EvalResultOutcomeBool,
			pgTemp,
			emrs.EvalContext.AIWorkflowAnalysisResult.ChunkOffset,
			emrs.EvalContext.EvalIterationCount,
		)
		if err != nil {
			log.Err(err).Msg("failed to execute query")
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
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
