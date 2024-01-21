package artemis_orchestrations

import (
	"context"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type DbJsonTriggerField struct {
	TriggerID            int64  `json:"schemaID"`
	TriggerGroup         string `json:"triggerGroup"`
	TriggerName          string `json:"triggerName"`
	TriggerAction        string `json:"triggerAction"`
	EvalTriggerAction    string `json:"evalTriggerAction"`
	EvalResultsTriggerOn string `json:"evalResultsTriggerOn"`
}

func SelectEvalFnsByOrgIDAndID(ctx context.Context, ou org_users.OrgUser, evalFnID int) ([]EvalFn, error) {
	params := []interface{}{
		ou.OrgID,
	}
	addOnQuery := ""
	addOnQuery2 := ""
	if evalFnID != 0 {
		params = append(params, evalFnID)
		addOnQuery = "AND f.eval_id = $2"
		addOnQuery2 = "AND m.eval_id = $2"
	}

	query := `
			WITH cte_fields_and_metrics AS (
				SELECT 
					m.eval_id,
					jsd.schema_id,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'fieldID', COALESCE(af.field_id, 0),
							'fieldName', COALESCE(af.field_name, ''),
							'fieldDescription', COALESCE(af.field_description, ''),
							'dataType', COALESCE(af.data_type, ''),
							'evalMetric', JSONB_BUILD_OBJECT(
								'evalMetricID', COALESCE(m.eval_metric_id, 0),
								'evalMetricResult', COALESCE(m.eval_metric_result, ''),
								'evalComparisonBoolean', COALESCE(m.eval_comparison_boolean, FALSE),
								'evalComparisonNumber', COALESCE(m.eval_comparison_number, 0.0),
								'evalComparisonString', COALESCE(m.eval_comparison_string, ''),
								'evalMetricDataType', COALESCE(af.data_type, ''),
								'evalOperator', COALESCE(m.eval_operator, ''),
								'evalState', COALESCE(m.eval_state, '')
							)
						)
					) AS fields_metrics_jsonb
				FROM public.eval_metrics m
				JOIN public.ai_fields af ON m.field_id = af.field_id
				JOIN public.ai_json_schema_definitions jsd ON af.schema_id = jsd.schema_id
				WHERE m.is_eval_metric_archived = false AND jsd.org_id = $1 ` + addOnQuery2 + `
				GROUP BY m.eval_id, jsd.schema_id
			),
			eval_fns_with_metrics AS (
				SELECT 
					f.eval_id, 
					f.org_id,
					f.user_id, 
					f.eval_name, 
					f.eval_type, 
					f.eval_group_name, 
					f.eval_model, 
					f.eval_format,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'schemaID', COALESCE(jsd.schema_id, 0),
							'schemaName', COALESCE(jsd.schema_name, ''),
							'schemaGroup', COALESCE(jsd.schema_group, 'default'),
							'isObjArray', COALESCE(jsd.is_obj_array, false),
							'fields', fm.fields_metrics_jsonb
						)
					) AS metrics_jsonb
				FROM public.eval_fns f
				JOIN cte_fields_and_metrics fm ON f.eval_id = fm.eval_id
				JOIN public.ai_json_schema_definitions jsd ON fm.schema_id = jsd.schema_id
				WHERE f.org_id = $1 ` + addOnQuery + `
				GROUP BY f.eval_id, f.org_id, f.user_id, f.eval_name, f.eval_type, f.eval_group_name, f.eval_model, f.eval_format
			), 
			cte_triggers AS (
				SELECT 
					em.eval_id,
					JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'triggerID', COALESCE(tab.trigger_id, 0),
						'triggerName', COALESCE(tab.trigger_name, ''),
						'triggerGroup', COALESCE(tab.trigger_group, ''),
						'triggerAction', COALESCE(tab.trigger_action, ''),
						'evalTriggerState', COALESCE(ta.eval_trigger_state, ''),
						'evalResultsTriggerOn', COALESCE(ta.eval_results_trigger_on, '')
					)
				) AS triggers_list
				FROM eval_fns_with_metrics em
				JOIN public.ai_trigger_actions_evals tae ON em.eval_id = tae.eval_id
				JOIN public.ai_trigger_eval ta ON ta.trigger_id = tae.trigger_id
				JOIN public.ai_trigger_actions tab ON tab.trigger_id = ta.trigger_id
				GROUP BY em.eval_id
			)
			SELECT 
				em.eval_id, 
				em.org_id, 
				em.user_id, 
				em.eval_name, 
				em.eval_type, 
				em.eval_group_name, 
				em.eval_model, 
				em.eval_format,
				COALESCE(ct.triggers_list, '[]'::jsonb) AS triggers_list, 
				COALESCE(em.metrics_jsonb, '[]'::jsonb) AS json_schemas
			FROM eval_fns_with_metrics em
			LEFT JOIN cte_triggers ct ON ct.eval_id = em.eval_id
			ORDER BY em.eval_id DESC
	`

	rows, err := apps.Pg.Query(ctx, query, params...)
	if err != nil {
		log.Err(err).Msg("failed to execute query")
		return nil, err
	}
	defer rows.Close()
	var evalFns []EvalFn

	for rows.Next() {
		ef := &EvalFn{}
		var dbTriggersHelper []DbJsonTriggerField
		//var dbFieldHelper []JsonSchemaDefinition
		err = rows.Scan(
			&ef.EvalID,
			&ef.OrgID,
			&ef.UserID,
			&ef.EvalName,
			&ef.EvalType,
			&ef.EvalGroupName,
			&ef.EvalModel,
			&ef.EvalFormat,
			&dbTriggersHelper,
			&ef.Schemas,
		)

		if err != nil {
			log.Err(err).Msg("failed to scan row")
			return nil, err
		}

		for _, trigger := range dbTriggersHelper {
			ta := TriggerAction{
				TriggerID:    int(trigger.TriggerID),
				TriggerGroup: trigger.TriggerGroup,
				TriggerName:  trigger.TriggerName,
			}
			ef.TriggerActions = append(ef.TriggerActions, ta)
		}

		evalFns = append(evalFns, *ef)
	}

	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error in row iteration")
		return nil, err
	}
	//sortEvalFnsByID(evalFns)
	return evalFns, nil
}
func sortEvalFnsByID(efs []EvalFn) {
	sort.Slice(efs, func(i, j int) bool {
		return *efs[i].EvalID > *efs[j].EvalID
	})
}
