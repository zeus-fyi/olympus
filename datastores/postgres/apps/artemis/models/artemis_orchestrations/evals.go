package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type EvalFn struct {
	EvalID         *int                    `json:"evalID,omitempty"`
	OrgID          int                     `json:"orgID,omitempty"`
	UserID         int                     `json:"userID,omitempty"`
	EvalName       string                  `json:"evalName"`
	EvalType       string                  `json:"evalType"`
	EvalGroupName  string                  `json:"evalGroupName"`
	EvalModel      *string                 `json:"evalModel,omitempty"`
	EvalFormat     string                  `json:"evalFormat"`
	EvalCycleCount int                     `json:"evalCycleCount,omitempty"`
	TriggerActions []TriggerAction         `json:"triggerFunctions,omitempty"`
	Schemas        []*JsonSchemaDefinition `json:"schemas,omitempty"`
}

type EvalMetric struct {
	EvalMetricID          *int     `json:"evalMetricID"`
	EvalMetricResult      string   `json:"evalMetricResult"`
	EvalComparisonBoolean *bool    `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *float64 `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string  `json:"evalComparisonString,omitempty"`
	EvalOperator          string   `json:"evalOperator"`
	EvalState             string   `json:"evalState"`
}

// this should map eval fn metric name -> to the result after evaluation

type EvalFnMetricResults struct {
	Map map[string]EvalMetricsResult `json:"map"`
}

func InsertOrUpdateEvalFnWithMetrics(ctx context.Context, ou org_users.OrgUser, evalFn *EvalFn) error {
	if evalFn == nil {
		return nil
	}
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)
	if evalFn.EvalID != nil && *evalFn.EvalID > 0 {
		tx, err = DeleteEvalMetricsAndTriggers(ctx, ou, tx, evalFn)
		if err != nil {
			log.Err(err).Msg("failed to delete eval metrics and triggers")
			return err
		}
	}
	ts := chronos.Chronos{}
	if evalFn.EvalID == nil {
		tv := ts.UnixTimeStampNow()
		evalFn.EvalID = &tv
	}
	// Inserting or updating eval_fns
	evalFnInsertOrUpdateQuery := `
        INSERT INTO eval_fns (eval_id, org_id, user_id, eval_name, eval_type, eval_group_name, eval_model, eval_format)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (eval_id) DO UPDATE SET
            eval_name = EXCLUDED.eval_name,
            eval_type = EXCLUDED.eval_type,
            eval_group_name = EXCLUDED.eval_group_name,
            eval_model = EXCLUDED.eval_model,
            eval_format = EXCLUDED.eval_format
        RETURNING eval_id;`
	err = tx.QueryRow(ctx, evalFnInsertOrUpdateQuery, evalFn.EvalID, evalFn.OrgID,
		evalFn.UserID, evalFn.EvalName, evalFn.EvalType, evalFn.EvalGroupName,
		evalFn.EvalModel, evalFn.EvalFormat).Scan(&evalFn.EvalID)
	if err != nil {
		log.Err(err).Msg("failed to insert or update eval_fns")
		return err
	}
	// Inserting or updating eval_metrics from json schema
	for _, schema := range evalFn.Schemas {
		for _, field := range schema.Fields {
			metric := field.EvalMetric
			if metric == nil {
				continue
			}
			if metric.EvalMetricID == nil {
				tv := ts.UnixTimeStampNow()
				metric.EvalMetricID = &tv
			}
			if metric.EvalState == "" {
				metric.EvalState = "info"
			}
			evalMetricInsertOrUpdateQuery := `
        INSERT INTO eval_metrics (eval_metric_id, eval_id, field_id, eval_metric_result, eval_comparison_boolean, eval_comparison_number, eval_comparison_string, eval_operator, eval_state)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (eval_metric_id) DO UPDATE SET
            eval_metric_result = EXCLUDED.eval_metric_result,
            eval_comparison_boolean = EXCLUDED.eval_comparison_boolean,
            eval_comparison_number = EXCLUDED.eval_comparison_number,
            eval_comparison_string = EXCLUDED.eval_comparison_string,
            eval_operator = EXCLUDED.eval_operator,
            eval_state = EXCLUDED.eval_state;`
			_, err = tx.Exec(ctx, evalMetricInsertOrUpdateQuery, metric.EvalMetricID, evalFn.EvalID, field.FieldID,
				metric.EvalMetricResult,
				metric.EvalComparisonBoolean, metric.EvalComparisonNumber, metric.EvalComparisonString,
				metric.EvalOperator, metric.EvalState)
			if err != nil {
				log.Err(err).Msg("failed to insert or update eval_metrics")
				return err
			}
		}
	}

	for _, eta := range evalFn.TriggerActions {
		for _, evTrig := range eta.EvalTriggerActions {
			query := `
            INSERT INTO ai_trigger_actions_evals(eval_id, trigger_id)
            VALUES ($1, $2)
         	ON CONFLICT (eval_id, trigger_id)
    		DO NOTHING ` // Adjust as needed
			_, err = tx.Exec(ctx, query, evalFn.EvalID, eta.TriggerID)
			if err != nil {
				log.Err(err).Interface("evTrig", evTrig).Msg("failed to insert eval trigger action")
				return err
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}

type EvalContext struct {
	EvalID                int `json:"evalID,omitempty"`
	WorkflowResultID      int `json:"workflowResultID,omitempty"`
	OrchestrationID       int `json:"orchestrationID"`
	SourceTaskID          int `json:"sourceTaskId"`
	RunningCycleNumber    int `json:"runningCycleNumber"`
	SearchWindowUnixStart int `json:"searchWindowUnixStart"`
	SearchWindowUnixEnd   int `json:"searchWindowUnixEnd"`
}

type EvalMetricsResults struct {
	EvalContext        EvalContext         `json:"evalContext"`
	EvalMetricsResults []EvalMetricsResult `json:"evalMetricsResults"`
}

type EvalMetricsResult struct {
	EvalName              string          `json:"evalName,omitempty"`
	EvalMetricName        string          `json:"evalMetricName"`
	EvalMetricID          int             `json:"evalMetricID,omitempty"`
	EvalMetricsResultID   int             `json:"evalMetricsResultId"`
	EvalMetricResult      string          `json:"evalMetricResult"` // pass or fail expected result
	EvalComparisonBoolean *bool           `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *float64        `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string         `json:"evalComparisonString,omitempty"`
	EvalComparisonInt     *float64        `json:"evalComparisonInt,omitempty"`
	EvalMetricDataType    string          `json:"evalMetricDataType"`
	EvalOperator          string          `json:"evalOperator"`
	EvalState             string          `json:"evalState"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart,omitempty"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd,omitempty"`
	EvalResultOutcome     bool            `json:"evalResultOutcome"` // true if eval passed, false if eval failed
	EvalMetadata          json.RawMessage `json:"evalMetadata,omitempty"`
}

func UpsertEvalMetricsResults(ctx context.Context, evCtx EvalContext, emrs []EvalMetricsResult) error {
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
            eval_metadata
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (eval_metric_id, source_task_id, orchestration_id, running_cycle_number)
        DO UPDATE SET
            orchestration_id = EXCLUDED.orchestration_id,
            source_task_id = EXCLUDED.source_task_id,
            eval_metric_id = EXCLUDED.eval_metric_id,
            running_cycle_number = EXCLUDED.running_cycle_number,
            search_window_unix_start = EXCLUDED.search_window_unix_start,
            search_window_unix_end = EXCLUDED.search_window_unix_end,
            eval_result_outcome = EXCLUDED.eval_result_outcome,
            eval_metadata = EXCLUDED.eval_metadata;
    `
	for _, emr := range emrs {
		ts := chronos.Chronos{}
		tsNow := ts.UnixTimeStampNow()
		_, err = tx.Exec(ctx, query,
			tsNow,
			evCtx.OrchestrationID,
			evCtx.SourceTaskID,
			emr.EvalMetricID,
			evCtx.RunningCycleNumber,
			evCtx.SearchWindowUnixStart,
			evCtx.SearchWindowUnixEnd,
			emr.EvalResultOutcome,
			&pgtype.JSONB{Bytes: sanitizeBytesUTF8(emr.EvalMetadata), Status: IsNull(emr.EvalMetadata)},
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

func DeleteEvalMetricsAndTriggers(ctx context.Context, ou org_users.OrgUser, tx pgx.Tx, evalFn *EvalFn) (pgx.Tx, error) {
	if evalFn == nil || tx == nil || evalFn.EvalID == nil || *evalFn.EvalID == 0 {
		return nil, nil
	}
	var keepMetricIDs []int

	var keepTriggerIDs []int
	for _, tgr := range evalFn.TriggerActions {
		keepTriggerIDs = append(keepTriggerIDs, tgr.TriggerID)
	}
	var keepFieldIds []int
	for _, schema := range evalFn.Schemas {
		for _, field := range schema.Fields {
			if field.EvalMetric != nil && field.EvalMetric.EvalMetricID != nil {
				keepMetricIDs = append(keepMetricIDs, *field.EvalMetric.EvalMetricID)
			}
			keepFieldIds = append(keepFieldIds, field.FieldID)
		}
	}
	// Using keepTriggerIDsArray in the delete query
	deleteDanglingMetricAndTriggerActionsQuery := `
	WITH cte_trigger_actions AS (
		SELECT ef.eval_id, te.trigger_id
		FROM ai_trigger_actions_evals te
		JOIN eval_fns ef ON ef.eval_id = te.eval_id 
		WHERE te.eval_id = $1 AND ef.org_id = $2 AND te.trigger_id = ANY($4)
	), cte_delete_trigger_actions AS (
		DELETE FROM ai_trigger_actions_evals te
		WHERE te.eval_id = $1 AND te.trigger_id NOT IN (SELECT trigger_id FROM cte_trigger_actions)
	)
		UPDATE public.eval_metrics
		SET is_eval_metric_archived = true,
			archived_at = NOW()
		WHERE eval_metrics.eval_id = $1
			AND eval_metric_id NOT IN (SELECT UNNEST($3::bigint[]))
			AND is_eval_metric_archived = false;`

	_, err := tx.Exec(ctx, deleteDanglingMetricAndTriggerActionsQuery, *evalFn.EvalID, ou.OrgID, pq.Array(keepMetricIDs), pq.Array(keepTriggerIDs))
	if err == pgx.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Err(err).Msg("failed to delete eval fn trigger eval actions")
		return tx, err
	}

	return tx, nil
}
