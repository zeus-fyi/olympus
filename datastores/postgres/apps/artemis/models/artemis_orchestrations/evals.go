package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
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
	TriggerActions []TriggerAction       `json:"triggerFunctions,omitempty"`
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

func DeleteEvalMetricsAndTriggers(ctx context.Context, ou org_users.OrgUser, tx pgx.Tx, evalFn *EvalFn) (pgx.Tx, error) {
	if evalFn == nil || tx == nil || evalFn.EvalID == nil || *evalFn.EvalID == 0 {
		return nil, nil
	}
	var keepMetricIDs []int
	for _, metric := range evalFn.EvalMetrics {
		if metric.EvalMetricID == nil {
			continue
		}
		keepMetricIDs = append(keepMetricIDs, *metric.EvalMetricID)
	}
	metricIDsArray := pq.Array(keepMetricIDs)

	var keepTriggerIDs []int
	for _, tgr := range evalFn.TriggerActions {
		keepTriggerIDs = append(keepTriggerIDs, tgr.TriggerID)
	}
	keepTriggerIDsArray := pq.Array(keepTriggerIDs)
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
	), cte_get_metrics_to_delete AS (
		SELECT em.eval_metric_id, ef.eval_id
		FROM eval_metrics em
		JOIN eval_fns ef ON em.eval_id = ef.eval_id
		WHERE ef.eval_id = $1 AND ef.org_id = $2 AND em.eval_metric_id != ANY($3)
	) DELETE FROM eval_metrics
	  WHERE eval_id = $1 AND eval_metric_id IN (SELECT eval_metric_id FROM cte_get_metrics_to_delete)`

	_, err := tx.Exec(ctx, deleteDanglingMetricAndTriggerActionsQuery, evalFn.EvalID, ou.OrgID, metricIDsArray, keepTriggerIDsArray)
	if err == pgx.ErrNoRows {
		err = nil
	}
	if err != nil {
		log.Err(err).Msg("failed to delete eval fn trigger eval actions")
		return tx, err
	}

	return tx, nil
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
            INSERT INTO eval_metrics (eval_metric_id, eval_id, eval_model_prompt, eval_metric_name, eval_metric_result, eval_comparison_boolean, eval_comparison_number, eval_comparison_string, eval_metric_data_type, eval_operator, eval_state)
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
	for _, eta := range evalFn.TriggerActions {
		for _, evTrig := range eta.EvalTriggerActions {
			query := `
            INSERT INTO ai_trigger_actions_evals(eval_id, trigger_id, eval_trigger_state, eval_results_trigger_on)
            VALUES ($1, $2, $3, $4)
         	ON CONFLICT (eval_id, trigger_id)
    		DO UPDATE SET
				eval_trigger_state = EXCLUDED.eval_trigger_state,
				eval_results_trigger_on = EXCLUDED.eval_results_trigger_on;` // Adjust as needed
			_, err = tx.Exec(ctx, query, evalFn.EvalID, eta.TriggerID, evTrig.EvalTriggerState, evTrig.EvalResultsTriggerOn)
			if err != nil {
				log.Err(err).Msg("failed to insert eval trigger action")
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
func SelectEvalFnsByOrgIDAndID(ctx context.Context, ou org_users.OrgUser, evalFnID int) ([]EvalFn, error) {
	params := []interface{}{
		ou.OrgID,
	}
	addOnQuery := ""
	if evalFnID != 0 {
		params = append(params, evalFnID)
		addOnQuery = "AND f.eval_id = $2"
	}

	query := `
    WITH eval_fns_with_metrics AS (
         SELECT	f.eval_id, f.org_id, f.user_id, f.eval_name, f.eval_type, f.eval_group_name, f.eval_model, f.eval_format,
                COALESCE(m.eval_metric_id, 0) AS eval_metric_id,
				COALESCE(m.eval_model_prompt, '') AS eval_model_prompt,
				COALESCE(m.eval_metric_name, '') AS eval_metric_name,
				COALESCE(m.eval_metric_result, '') AS eval_metric_result,
				COALESCE(m.eval_comparison_boolean, FALSE) AS eval_comparison_boolean,
				COALESCE(m.eval_comparison_number, 0.0) AS eval_comparison_number,
				COALESCE(m.eval_comparison_string, '') AS eval_comparison_string,
				COALESCE(m.eval_metric_data_type, '') AS eval_metric_data_type,
				COALESCE(m.eval_operator, '') AS eval_operator,
				COALESCE(m.eval_state, '') AS eval_state,
			   	COALESCE(tab.trigger_id, 0), COALESCE(tab.trigger_name, ''), COALESCE(tab.trigger_group, ''),
 			   	COALESCE(tab.trigger_action, ''), COALESCE(ta.eval_trigger_state, ''), COALESCE(ta.eval_results_trigger_on, '')
        FROM public.eval_fns f
        LEFT JOIN public.eval_metrics m ON f.eval_id = m.eval_id
        LEFT JOIN public.ai_trigger_actions_evals ta ON f.eval_id = ta.eval_id
		LEFT JOIN public.ai_trigger_actions tab ON ta.trigger_id = tab.trigger_id
        WHERE f.org_id = $1 ` + addOnQuery + `
    )
    SELECT * FROM eval_fns_with_metrics;`

	rows, err := apps.Pg.Query(ctx, query, params...)
	if err != nil {
		log.Err(err).Msg("failed to execute query")
		return nil, err
	}
	defer rows.Close()

	tm := make(map[int]map[int]*TriggerAction)
	evalFnsMap := make(map[int]*EvalFn)
	for rows.Next() {
		var ef EvalFn
		var em EvalMetric
		var ta TriggerAction
		var eta EvalTriggerActions
		var evalID int
		err = rows.Scan(&evalID, &ef.OrgID, &ef.UserID, &ef.EvalName, &ef.EvalType, &ef.EvalGroupName, &ef.EvalModel, &ef.EvalFormat,
			&em.EvalMetricID, &em.EvalModelPrompt, &em.EvalMetricName, &em.EvalMetricResult, &em.EvalComparisonBoolean,
			&em.EvalComparisonNumber, &em.EvalComparisonString, &em.EvalMetricDataType, &em.EvalOperator, &em.EvalState,
			&ta.TriggerID, &ta.TriggerName, &ta.TriggerGroup,
			&ta.TriggerAction, &eta.EvalTriggerState, &eta.EvalResultsTriggerOn) // Scan for TriggerActions
		if err != nil {
			log.Err(err).Msg("failed to scan row")
			return nil, err
		}
		eta.EvalID = evalID
		eta.TriggerID = ta.TriggerID
		if _, ok := tm[evalID]; !ok {
			tm[evalID] = make(map[int]*TriggerAction)
		}

		if ta.TriggerID != 0 {
			if _, tok := tm[evalID][ta.TriggerID]; !tok {
				if eta.EvalTriggerState != "" && eta.EvalResultsTriggerOn != "" {
					ta.EvalTriggerActions = append(ta.EvalTriggerActions, eta)
				}
				tm[evalID][ta.TriggerID] = &ta
			}
		}

		if existingEvalFn, exists := evalFnsMap[evalID]; exists {
			if em.EvalMetricID != nil && *em.EvalMetricID > 0 {
				existingEvalFn.EvalMetrics = append(existingEvalFn.EvalMetrics, em)
				existingEvalFn.EvalMetricMap[em.EvalMetricName] = em
			}

		} else {
			ef.EvalID = &evalID
			if em.EvalMetricID != nil && *em.EvalMetricID > 0 {
				ef.EvalMetrics = append(ef.EvalMetrics, em)
				ef.EvalMetricMap = make(map[string]EvalMetric)
				ef.EvalMetricMap[em.EvalMetricName] = em
			}
			evalFnsMap[evalID] = &ef
		}
	}
	var evalFns []EvalFn
	for _, ef := range evalFnsMap {
		if ef == nil || ef.EvalID == nil || *ef.EvalID == 0 {
			continue
		}
		for _, efts := range tm[*ef.EvalID] {
			if efts == nil {
				continue
			}
			ef.TriggerActions = append(ef.TriggerActions, *efts)
		}
		evalFns = append(evalFns, *ef)
	}
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error in row iteration")
		return nil, err
	}
	sortEvalFnsByID(evalFns)
	return evalFns, nil
}
func sortEvalFnsByID(efs []EvalFn) {
	sort.Slice(efs, func(i, j int) bool {
		return *efs[i].EvalID > *efs[j].EvalID
	})
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
