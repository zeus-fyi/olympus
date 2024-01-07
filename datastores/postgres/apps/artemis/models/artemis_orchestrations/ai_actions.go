package artemis_orchestrations

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type TriggerAction struct {
	TriggerID          int                  `db:"trigger_id" json:"triggerId"`
	OrgID              int                  `db:"org_id" json:"orgId"`
	UserID             int                  `db:"user_id" json:"userId"`
	TriggerName        string               `db:"trigger_name" json:"triggerName"`
	TriggerGroup       string               `db:"trigger_group" json:"triggerGroup"`
	EvalTriggerActions []EvalTriggerActions `db:"eval_trigger_actions" json:"evalTriggerActions"`
}

type EvalTriggerActions struct {
	EvalID               int    `db:"eval_id" json:"evalId"`
	TriggerID            int    `db:"trigger_id" json:"triggerId"`
	EvalTriggerState     string `db:"eval_trigger_state" json:"evalTriggerState"`
	EvalResultsTriggerOn string `db:"eval_results_trigger_on" json:"evalResultsTriggerOn"`
}

func CreateOrUpdateAction(ctx context.Context, ou org_users.OrgUser, trigger *TriggerAction) error {
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
        INSERT INTO public.ai_trigger_actions (org_id, user_id, trigger_name, trigger_group)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (org_id, trigger_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            trigger_group = EXCLUDED.trigger_group
        RETURNING trigger_id;`

	err = tx.QueryRow(ctx, q.RawQuery, ou.OrgID, ou.UserID, trigger.TriggerName, trigger.TriggerGroup).Scan(&trigger.TriggerID)
	if err != nil {
		log.Err(err).Msg("failed to insert ai trigger action")
		return err
	}

	// Insert eval trigger actions if any
	for _, eta := range trigger.EvalTriggerActions {
		q.RawQuery = `
            INSERT INTO public.ai_eval_trigger_actions (eval_id, trigger_id, eval_trigger_state, eval_results_trigger_on)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (eval_id, trigger_id) 
            DO NOTHING;` // Adjust as needed

		_, err = tx.Exec(ctx, q.RawQuery, eta.EvalID, trigger.TriggerID, eta.EvalTriggerState, eta.EvalResultsTriggerOn)
		if err != nil {
			log.Err(err).Msg("failed to insert eval trigger action")
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
