package artemis_orchestrations

import (
	"context"
	"errors"

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
        INSERT INTO public.ai_trigger_actions (org_id, user_id, trigger_name, trigger_group, trigger_action)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (org_id, trigger_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            trigger_action = EXCLUDED.trigger_action,
            trigger_group = EXCLUDED.trigger_group
        RETURNING trigger_id;`

	err = tx.QueryRow(ctx, q.RawQuery, ou.OrgID, ou.UserID, trigger.TriggerName, trigger.TriggerGroup, trigger.TriggerAction).Scan(&trigger.TriggerID)
	if err != nil {
		log.Err(err).Msg("failed to insert ai trigger action")
		return err
	}
	for _, eta := range trigger.EvalTriggerActions {
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
		if eta.EvalID != 0 {
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
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}
