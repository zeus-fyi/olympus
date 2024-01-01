package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type EvalFn struct {
	EvalID         *int                  `json:"evalID,omitempty"`
	OrgID          int                   `json:"orgID,omitempty"`
	UserID         int                   `json:"userID,omitempty"`
	EvalName       string                `json:"evalName"`
	EvalType       string                `json:"evalType"`
	EvalGroupName  string                `json:"evalGroupName"`
	EvalModel      *string               `json:"evalModel,omitempty"`
	EvalFormat     string                `json:"evalFormat"`
	EvalMetrics    []EvalMetric          `json:"evalMetrics"`
	EvalMetricMap  map[string]EvalMetric `json:"evalMetricMap,omitempty"`
	EvalCycleCount int                   `json:"evalCycleCount,omitempty"`
}

type EvalMetric struct {
	EvalMetricID          *int     `json:"evalMetricID"`
	EvalModelPrompt       string   `json:"evalModelPrompt"`
	EvalMetricName        string   `json:"evalMetricName"`
	EvalMetricResult      string   `json:"evalMetricResult"`
	EvalComparisonBoolean *bool    `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *float64 `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string  `json:"evalComparisonString,omitempty"`
	EvalMetricDataType    string   `json:"evalMetricDataType"`
	EvalOperator          string   `json:"evalOperator"`
	EvalState             string   `json:"evalState"`
}

// this should map eval fn metric name -> to the result after evaluation

type EvalFnMetricResults struct {
	Map map[string]EvalMetricsResult `json:"map"`
}

func InsertOrUpdateEvalFnWithMetrics(ctx context.Context, evalFn *EvalFn) error {
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("failed to begin transaction")
		return err
	}
	defer tx.Rollback(ctx)
	ts := chronos.Chronos{}
	if evalFn.EvalID == nil {
		tv := ts.UnixTimeStampNow()
		evalFn.EvalID = &tv
	}
	// Inserting or updating eval_fns
	evalFnInsertOrUpdateQuery := `
        INSERT INTO public.eval_fns (eval_id, org_id, user_id, eval_name, eval_type, eval_group_name, eval_model, eval_format)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (eval_id) DO UPDATE SET
            org_id = EXCLUDED.org_id,
            user_id = EXCLUDED.user_id,
            eval_name = EXCLUDED.eval_name,
            eval_type = EXCLUDED.eval_type,
            eval_group_name = EXCLUDED.eval_group_name,
            eval_model = EXCLUDED.eval_model,
            eval_format = EXCLUDED.eval_format
        RETURNING eval_id;`
	err = tx.QueryRow(ctx, evalFnInsertOrUpdateQuery, evalFn.EvalID, evalFn.OrgID, evalFn.UserID, evalFn.EvalName, evalFn.EvalType, evalFn.EvalGroupName, evalFn.EvalModel, evalFn.EvalFormat).Scan(&evalFn.EvalID)
	if err != nil {
		log.Err(err).Msg("failed to insert or update eval_fns")
		return err
	}
	// Inserting or updating eval_metrics
	for _, metric := range evalFn.EvalMetrics {
		if metric.EvalMetricID == nil {
			tv := ts.UnixTimeStampNow()
			metric.EvalMetricID = &tv
		}
		evalMetricInsertOrUpdateQuery := `
            INSERT INTO public.eval_metrics (eval_metric_id, eval_id, eval_model_prompt, eval_metric_name, eval_metric_result, eval_comparison_boolean, eval_comparison_number, eval_comparison_string, eval_metric_data_type, eval_operator, eval_state)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
            ON CONFLICT (eval_metric_id, eval_id) DO UPDATE SET
                eval_id = EXCLUDED.eval_id,
                eval_model_prompt = EXCLUDED.eval_model_prompt,
                eval_metric_name = EXCLUDED.eval_metric_name,
                eval_metric_result = EXCLUDED.eval_metric_result,
                eval_comparison_boolean = EXCLUDED.eval_comparison_boolean,
                eval_comparison_number = EXCLUDED.eval_comparison_number,
                eval_comparison_string = EXCLUDED.eval_comparison_string,
                eval_metric_data_type = EXCLUDED.eval_metric_data_type,
                eval_operator = EXCLUDED.eval_operator,
                eval_state = EXCLUDED.eval_state;`
		_, err = tx.Exec(ctx, evalMetricInsertOrUpdateQuery, metric.EvalMetricID, evalFn.EvalID, metric.EvalModelPrompt, metric.EvalMetricName, metric.EvalMetricResult, metric.EvalComparisonBoolean, metric.EvalComparisonNumber, metric.EvalComparisonString, metric.EvalMetricDataType, metric.EvalOperator, metric.EvalState)
		if err != nil {
			log.Err(err).Msg("failed to insert or update eval_fns")
			return err
		}
	}
	return tx.Commit(ctx)
}

func SelectEvalFnsByOrgID(ctx context.Context, ou org_users.OrgUser) ([]EvalFn, error) {
	const query = `
        SELECT eval_id, eval_name, eval_type, eval_group_name, eval_model, eval_format
        FROM public.eval_fns
        WHERE org_id = $1;`
	rows, err := apps.Pg.Query(ctx, query, ou.OrgID)
	if err != nil {
		log.Err(err).Msg("failed to select eval_fns")
		return nil, err
	}
	defer rows.Close()
	var evalFns []EvalFn
	for rows.Next() {
		var ef EvalFn
		err = rows.Scan(&ef.EvalID, &ef.EvalName, &ef.EvalType, &ef.EvalGroupName, &ef.EvalModel, &ef.EvalFormat)
		if err != nil {
			log.Err(err).Msg("failed to select eval_fns")
			return nil, err
		}
		evalFns = append(evalFns, ef)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("failed to select eval_fns")
		return nil, err
	}
	return evalFns, nil
}

func SelectEvalFnsByOrgIDAndID(ctx context.Context, ou org_users.OrgUser, evalFnID int) ([]EvalFn, error) {
	const query = `
    WITH eval_fns_with_metrics AS (
        SELECT f.eval_id, f.org_id, f.user_id, f.eval_name, f.eval_type, f.eval_group_name, f.eval_model, f.eval_format,
               m.eval_metric_id, m.eval_model_prompt, m.eval_metric_name, m.eval_metric_result, m.eval_comparison_boolean,
               m.eval_comparison_number, m.eval_comparison_string, m.eval_metric_data_type, m.eval_operator, m.eval_state
        FROM public.eval_fns f
        LEFT JOIN public.eval_metrics m ON f.eval_id = m.eval_id
        WHERE f.org_id = $1 AND f.eval_id = $2
    )
    SELECT * FROM eval_fns_with_metrics;`

	rows, err := apps.Pg.Query(ctx, query, ou.OrgID, evalFnID)
	if err != nil {
		log.Err(err).Msg("failed to execute query")
		return nil, err
	}
	defer rows.Close()

	evalFnsMap := make(map[int]*EvalFn)
	for rows.Next() {
		var ef EvalFn
		var em EvalMetric
		var evalID int
		err = rows.Scan(&evalID, &ef.OrgID, &ef.UserID, &ef.EvalName, &ef.EvalType, &ef.EvalGroupName, &ef.EvalModel, &ef.EvalFormat,
			&em.EvalMetricID, &em.EvalModelPrompt, &em.EvalMetricName, &em.EvalMetricResult, &em.EvalComparisonBoolean,
			&em.EvalComparisonNumber, &em.EvalComparisonString, &em.EvalMetricDataType, &em.EvalOperator, &em.EvalState)
		if err != nil {
			log.Err(err).Msg("failed to scan row")
			return nil, err
		}
		if existingEvalFn, exists := evalFnsMap[evalID]; exists {
			existingEvalFn.EvalMetrics = append(existingEvalFn.EvalMetrics, em)
			existingEvalFn.EvalMetricMap[em.EvalMetricName] = em
		} else {
			ef.EvalID = &evalID
			ef.EvalMetrics = append(ef.EvalMetrics, em)
			ef.EvalMetricMap = make(map[string]EvalMetric)
			ef.EvalMetricMap[em.EvalMetricName] = em
			evalFnsMap[evalID] = &ef
		}
	}

	var evalFns []EvalFn
	for _, ef := range evalFnsMap {
		if ef == nil {
			continue
		}
		evalFns = append(evalFns, *ef)
	}

	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error in row iteration")
		return nil, err
	}

	return evalFns, nil
}

type EvalContext struct {
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
	EvalMetricResult      string          `json:"evalMetricResult"`
	EvalComparisonBoolean *bool           `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *float64        `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string         `json:"evalComparisonString,omitempty"`
	EvalMetricDataType    string          `json:"evalMetricDataType"`
	EvalOperator          string          `json:"evalOperator"`
	EvalState             string          `json:"evalState"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart,omitempty"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd,omitempty"`
	EvalResultOutcome     bool            `json:"evalResultOutcome"`
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
