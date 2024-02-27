package artemis_orchestrations

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func CreateOrUpdateTriggerAction(ctx context.Context, ou org_users.OrgUser, trigger *TriggerAction) error {
	if trigger == nil {
		return errors.New("trigger cannot be nil")
	}
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	// Defer a rollback in case of failure
	defer tx.Rollback(ctx)
	// Insert or update the ai_trigger_actions
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
        INSERT INTO public.ai_trigger_actions (org_id, user_id, trigger_name, trigger_group, trigger_action, expires_after_seconds)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, trigger_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            trigger_action = EXCLUDED.trigger_action,
            trigger_group = EXCLUDED.trigger_group,
        	expires_after_seconds = EXCLUDED.expires_after_seconds
        RETURNING trigger_id, trigger_id::text;`

	if trigger.TriggerExpirationDuration > 0 && trigger.TriggerExpirationTimeUnit != "" {
		trigger.TriggerExpiresAfterSeconds = CalculateStepSizeUnix(int(trigger.TriggerExpirationDuration), trigger.TriggerExpirationTimeUnit)
	} else {
		trigger.TriggerExpiresAfterSeconds = 0
	}
	err = tx.QueryRow(ctx, q.RawQuery, ou.OrgID, ou.UserID, trigger.TriggerName,
		trigger.TriggerGroup, trigger.TriggerAction, trigger.TriggerExpiresAfterSeconds,
	).Scan(&trigger.TriggerID, &trigger.TriggerStrID)
	if err != nil {
		log.Err(err).Msg("failed to insert ai trigger action")
		return err
	}
	for ein, eta := range trigger.EvalTriggerActions {
		if eta.EvalResultsTriggerOn == "" || eta.EvalTriggerState == "" {
			continue
		}
		q.RawQuery = `
            INSERT INTO public.ai_trigger_eval(trigger_id, eval_trigger_state, eval_results_trigger_on)
            VALUES ($1, $2, $3)
         	ON CONFLICT (trigger_id)
         	DO UPDATE SET
				eval_trigger_state = EXCLUDED.eval_trigger_state,
				eval_results_trigger_on = EXCLUDED.eval_results_trigger_on;`
		_, err = tx.Exec(ctx, q.RawQuery, trigger.TriggerID, eta.EvalTriggerState, eta.EvalResultsTriggerOn)
		if err != nil {
			log.Err(err).Msg("failed to insert eval trigger action")
			return err
		}
		if eta.EvalID != 0 || eta.EvalStrID != "" {
			if eta.EvalStrID != "" {
				ei, aerr := strconv.Atoi(eta.EvalStrID)
				if aerr != nil {
					log.Err(aerr).Msg("failed to convert eval id to int")
					return aerr
				}
				eta.EvalID = ei
				trigger.EvalTriggerActions[ein].EvalID = ei
			}

			q.RawQuery = `
            INSERT INTO public.ai_trigger_actions_evals(eval_id, trigger_id)
            VALUES ($1, $2)
         	ON CONFLICT (eval_id, trigger_id)
    		DO NOTHING;`
			_, err = tx.Exec(ctx, q.RawQuery, eta.EvalID, trigger.TriggerID)
			if err != nil {
				log.Err(err).Msg("failed to insert eval trigger action")
				return err
			}
		}
	}
	var rids []int
	for i, retrieval := range trigger.TriggerRetrievals {
		if retrieval.RetrievalID == nil {
			continue
		}
		if aws.StringValue(retrieval.RetrievalStrID) != "" {
			ri, aerr := strconv.Atoi(*retrieval.RetrievalStrID)
			if aerr != nil {
				log.Err(aerr).Msg("failed to convert retrieval id to int")
				return aerr
			}
			trigger.TriggerRetrievals[i].RetrievalID = &ri
			retrieval.RetrievalID = &ri
		}
		rids = append(rids, *retrieval.RetrievalID)
		q.RawQuery = `
			INSERT INTO public.ai_trigger_actions_api(trigger_id, retrieval_id)
			VALUES ($1, $2)
		 	ON CONFLICT (trigger_id, retrieval_id)
			DO NOTHING;`
		_, err = tx.Exec(ctx, q.RawQuery, trigger.TriggerID, retrieval.RetrievalID)
		if err != nil {
			log.Err(err).Msg("failed to insert eval trigger action")
			return err
		}
	}
	if len(rids) > 0 {
		// Building a query string with placeholder for array
		query := `
				DELETE FROM public.ai_trigger_actions_api
				WHERE trigger_id = $1 AND retrieval_id NOT IN (SELECT UNNEST($2::bigint[]));`
		// Execute the delete query
		_, err = tx.Exec(ctx, query, trigger.TriggerID, pq.Array(rids))
		if err != nil {
			log.Err(err).Msg("failed to delete unwanted eval trigger actions")
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
