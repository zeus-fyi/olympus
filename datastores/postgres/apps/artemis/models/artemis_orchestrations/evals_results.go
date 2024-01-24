package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AIWorkflowEvalResultResponse struct {
	EvalResultID     int `db:"eval_result_id" json:"evalResultId"`
	WorkflowResultID int `db:"workflow_result_id" json:"workflowResultId"`
	EvalID           int `db:"eval_id" json:"evalId"`
	ResponseID       int `db:"response_id" json:"responseId"`
}

func UpsertEvalMetricsResults(ctx context.Context, evCtx EvalContext, emrs []EvalMetric) error {
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	const query = `
       INSERT INTO public.eval_metrics_results (
           eval_metrics_result_id,
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
       ON CONFLICT (eval_metric_id, source_task_id, orchestration_id, running_cycle_number, chunk_offset, eval_iteration_count)
       DO UPDATE SET
           running_cycle_number = EXCLUDED.running_cycle_number,
           search_window_unix_start = EXCLUDED.search_window_unix_start,
           search_window_unix_end = EXCLUDED.search_window_unix_end,
           eval_result_outcome = EXCLUDED.eval_result_outcome,
           eval_metadata = EXCLUDED.eval_metadata;
   `
	for _, emr := range emrs {
		ts := chronos.Chronos{}
		tsNow := ts.UnixTimeStampNow()

		var pgTemp *pgtype.JSONB

		if emr.EvalMetricResult == nil {
			pgTemp = &pgtype.JSONB{Bytes: []byte{}, Status: IsNull(nil)}
		} else {
			pgTemp = &pgtype.JSONB{Bytes: sanitizeBytesUTF8(emr.EvalMetricResult.EvalMetadata), Status: IsNull(emr.EvalMetricResult.EvalMetadata)}
		}
		_, err = tx.Exec(ctx, query,
			tsNow,
			evCtx.AIWorkflowAnalysisResult.OrchestrationsID,
			evCtx.AIWorkflowAnalysisResult.SourceTaskID,
			emr.EvalMetricID,
			evCtx.AIWorkflowAnalysisResult.RunningCycleNumber,
			evCtx.AIWorkflowAnalysisResult.SearchWindowUnixStart,
			evCtx.AIWorkflowAnalysisResult.SearchWindowUnixEnd,
			emr.EvalExpectedResultState,
			pgTemp,
			evCtx.AIWorkflowAnalysisResult.ChunkOffset,
			evCtx.EvalIterationCount,
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
	q.RawQuery = `INSERT INTO ai_workflow_eval_result_response(eval_result_id, workflow_result_id, eval_id, response_id)
                  VALUES ($1, $2, $3, $4)
                  ON CONFLICT (eval_result_id) 
                  DO UPDATE SET 
                      workflow_result_id = EXCLUDED.workflow_result_id,
                      eval_id = EXCLUDED.eval_id,
                      response_id = EXCLUDED.response_id
                  RETURNING eval_result_id;`

	if errr.EvalResultID <= 0 {
		ch := chronos.Chronos{}
		errr.EvalResultID = ch.UnixTimeStampNow()
	}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, errr.EvalResultID, errr.WorkflowResultID, errr.EvalID, errr.ResponseID).Scan(&errr.EvalResultID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("AIWorkflowEvalResultResponse")); returnErr != nil {
		log.Err(returnErr).Interface("errr", errr).Msg(q.LogHeader("AIWorkflowEvalResultResponse"))
		return errr.EvalResultID, err
	}
	return errr.EvalResultID, nil
}

// eval_metrics_results
/*
CREATE TABLE public.ai_workflow_eval_result_responses(
    eval_result_id BIGINT NOT NULL PRIMARY KEY,
    workflow_result_id BIGINT NOT NULL REFERENCES ai_workflow_analysis_results(workflow_result_id),
    eval_id BIGINT NOT NULL REFERENCES eval_fns(eval_id),
    response_id int8 NOT NULL REFERENCES completion_responses(response_id)
);

CREATE TABLE public.eval_metrics_results(
    eval_metrics_result_id int8 NOT NULL DEFAULT next_id() PRIMARY KEY,
    orchestration_id int8 NOT NULL REFERENCES orchestrations(orchestration_id),
    source_task_id int8 NOT NULL REFERENCES ai_task_library(task_id),
    eval_metric_id BIGINT NOT NULL REFERENCES public.eval_metrics(eval_metric_id),
    running_cycle_number int8 NOT NULL DEFAULT 1,
    search_window_unix_start int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    search_window_unix_end int8 NOT NULL CHECK (search_window_unix_start < search_window_unix_end),
    eval_iteration_count int8 NOT NULL DEFAULT 0,
    chunk_offset int8 NOT NULL DEFAULT 0,
    eval_result_outcome boolean NOT NULL,
    eval_metadata jsonb
);
ALTER TABLE public.eval_metrics_results
    ADD COLUMN eval_iteration int8 NOT NULL DEFAULT 0;

CREATE INDEX eval_result_outcome_idx ON public.eval_metrics_results("eval_result_outcome");
CREATE INDEX eval_result_metric_idx ON public.eval_metrics_results("eval_metric_id");
CREATE INDEX eval_result_orch_id_idx ON public.eval_metrics_results("orchestration_id");
CREATE INDEX eval_result_cycle_idx ON public.eval_metrics_results("running_cycle_number");
CREATE INDEX eval_result_source_search_start_idx ON public.eval_metrics_results("search_window_unix_start");
CREATE INDEX eval_result_source_search_end_idx ON public.eval_metrics_results("search_window_unix_end");
CREATE INDEX eval_result_eval_iter_idx ON public.eval_metrics_results("eval_iteration_count");
CREATE INDEX eval_result_eval_chunk_idx ON public.eval_metrics_results("chunk_offset");
ALTER TABLE public.eval_metrics_results
ADD CONSTRAINT unique_eval_metrics_combination UNIQUE (eval_metric_id, source_task_id, orchestration_id, running_cycle_number, chunk_offset, eval_iteration_count);


*/
