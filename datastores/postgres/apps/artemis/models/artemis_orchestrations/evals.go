package artemis_orchestrations

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

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
			for _, metric := range field.EvalMetrics {
				if metric == nil {
					continue
				}
				if metric.EvalMetricID == nil || aws.ToInt(metric.EvalMetricID) == 0 {
					tv := ts.UnixTimeStampNow()
					metric.EvalMetricID = &tv
				}
				if metric.EvalState == "" {
					metric.EvalState = "info"
				}

				if metric.EvalMetricComparisonValues == nil {
					metric.EvalMetricComparisonValues = &EvalMetricComparisonValues{}
				}
				evalMetricInsertOrUpdateQuery := `
        INSERT INTO eval_metrics (eval_metric_id, eval_id, field_id, eval_metric_result,
                                  eval_comparison_integer, eval_comparison_boolean, eval_comparison_number, eval_comparison_string, eval_operator, eval_state)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (eval_metric_id) DO UPDATE SET
            eval_metric_result = EXCLUDED.eval_metric_result,
            eval_comparison_integer = EXCLUDED.eval_comparison_integer,
            eval_comparison_boolean = EXCLUDED.eval_comparison_boolean,
            eval_comparison_number = EXCLUDED.eval_comparison_number,
            eval_comparison_string = EXCLUDED.eval_comparison_string,
            eval_operator = EXCLUDED.eval_operator,
            eval_state = EXCLUDED.eval_state;`
				_, err = tx.Exec(ctx, evalMetricInsertOrUpdateQuery,
					metric.EvalMetricID,
					evalFn.EvalID,
					field.FieldID,
					metric.EvalExpectedResultState,
					metric.EvalMetricComparisonValues.EvalComparisonInteger,
					metric.EvalMetricComparisonValues.EvalComparisonBoolean,
					metric.EvalMetricComparisonValues.EvalComparisonNumber,
					metric.EvalMetricComparisonValues.EvalComparisonString,
					metric.EvalOperator,
					metric.EvalState)
				if err != nil {
					log.Err(err).Msg("failed to insert or update eval_metrics")
					return err
				}
			}
		}
	}
	for ei, eta := range evalFn.TriggerActions {
		if eta.TriggerStrID != "" {
			eta.TriggerID, err = strconv.Atoi(eta.TriggerStrID)
			if err != nil {
				log.Err(err).Msg("failed to parse int")
				return err
			}
			evalFn.TriggerActions[ei].TriggerID = eta.TriggerID
		}
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
			for _, metric := range field.EvalMetrics {
				if metric == nil {
					continue
				}
				if metric.EvalMetricID != nil {
					keepMetricIDs = append(keepMetricIDs, *metric.EvalMetricID)
				}
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
		log.Warn().Msg("no rows found")
		err = nil
	}
	if err != nil {
		log.Err(err).Msg("failed to delete eval fn trigger eval actions")
		return tx, err
	}

	return tx, nil
}
