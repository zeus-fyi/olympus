package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type EvalFn struct {
	EvalID        *int         `json:"evalID,omitempty"`
	OrgID         int          `json:"orgID,omitempty"`
	UserID        int          `json:"userID,omitempty"`
	EvalName      string       `json:"evalName"`
	EvalType      string       `json:"evalType"`
	EvalGroupName string       `json:"evalGroupName"`
	EvalModel     *string      `json:"evalModel,omitempty"`
	EvalFormat    string       `json:"evalFormat"`
	EvalMetrics   []EvalMetric `json:"evalMetrics"`
}

type EvalMetric struct {
	EvalMetricID          *int    `json:"evalMetricID"`
	EvalModelPrompt       string  `json:"evalModelPrompt"`
	EvalMetricName        string  `json:"evalMetricName"`
	EvalMetricResult      string  `json:"evalMetricResult"`
	EvalComparisonBoolean *bool   `json:"evalComparisonBoolean,omitempty"`
	EvalComparisonNumber  *int    `json:"evalComparisonNumber,omitempty"`
	EvalComparisonString  *string `json:"evalComparisonString,omitempty"`
	EvalMetricDataType    string  `json:"evalMetricDataType"`
	EvalOperator          string  `json:"evalOperator"`
	EvalState             string  `json:"evalState"`
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

func SelectEvalFnsByOrgID(ctx context.Context, orgID int) ([]EvalFn, error) {
	const query = `
        SELECT eval_id, org_id, user_id, eval_name, eval_type, eval_group_name, eval_model, eval_format
        FROM public.eval_fns
        WHERE org_id = $1;`
	rows, err := apps.Pg.Query(ctx, query, orgID)
	if err != nil {
		log.Err(err).Msg("failed to select eval_fns")
		return nil, err
	}
	defer rows.Close()
	var evalFns []EvalFn
	for rows.Next() {
		var ef EvalFn
		err = rows.Scan(&ef.EvalID, &ef.OrgID, &ef.UserID, &ef.EvalName, &ef.EvalType, &ef.EvalGroupName, &ef.EvalModel, &ef.EvalFormat)
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
