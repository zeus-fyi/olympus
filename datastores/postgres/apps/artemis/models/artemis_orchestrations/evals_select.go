package artemis_orchestrations

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type DbJsonTriggerField struct {
	TriggerID            int    `json:"triggerID"`
	TriggerName          string `json:"triggerName"`
	TriggerGroup         string `json:"triggerGroup"`
	TriggerAction        string `json:"triggerAction"`
	EvalTriggerState     string `json:"evalTriggerState"`
	EvalResultsTriggerOn string `json:"evalResultsTriggerOn"`
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
			WITH cte_metrics AS (
				SELECT 
					m.eval_id,
					m.field_id,
					af.schema_id,
					af.field_name,
					af.field_description,
					af.data_type,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
								'evalMetricID', COALESCE(m.eval_metric_id, 0),
								'evalExpectedResultState', COALESCE(m.eval_metric_result, ''),
								'evalMetricComparisonValues', JSONB_BUILD_OBJECT(
										'evalComparisonBoolean', COALESCE(m.eval_comparison_boolean, FALSE),
										'evalComparisonNumber', COALESCE(m.eval_comparison_number, 0.0),
										'evalComparisonString', COALESCE(m.eval_comparison_string, ''),
										'evalComparisonInteger', COALESCE(m.eval_comparison_integer, 0)
									),
								'evalOperator', COALESCE(m.eval_operator, ''),
								'evalState', COALESCE(m.eval_state, '')
							)
					) AS fields_metrics_jsonb
				FROM public.eval_fns f
				JOIN public.eval_metrics m ON m.eval_id = f.eval_id
				JOIN public.ai_fields af ON af.field_id = m.field_id AND af.is_field_archived = false
				WHERE m.is_eval_metric_archived = false AND f.org_id = $1 ` + addOnQuery + `
				GROUP BY m.eval_id, m.field_id,
					af.schema_id,
					af.field_name,
					af.field_description,
					af.data_type
			), cte_fields_and_metrics AS (
				SELECT 
					f.eval_id,
					jsd.schema_id,
					jsd.schema_name,
					jsd.schema_group,
					jsd.schema_description,
					jsd.is_obj_array,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'fieldID', COALESCE(f.field_id, 0),
							'fieldName', COALESCE(f.field_name, ''),
							'fieldDescription', COALESCE(f.field_description, ''),
							'dataType', COALESCE(f.data_type, ''),
							'evalMetrics', f.fields_metrics_jsonb
						)
					) AS fields_jsonb
				FROM cte_metrics f
				JOIN public.ai_json_schema_definitions jsd ON f.schema_id = jsd.schema_id
				GROUP BY f.eval_id, jsd.schema_id,jsd.schema_name,
					jsd.schema_group,
					jsd.schema_description,
					jsd.is_obj_array
			), eval_fns_with_metrics AS (
				SELECT 
					f.eval_id, 
					f.eval_name, 
					f.eval_type, 
					f.eval_group_name, 
					f.eval_model, 
					f.eval_format,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'schemaID', COALESCE(fm.schema_id, 0),
							'schemaName', COALESCE(fm.schema_name, ''),
							'schemaGroup', COALESCE(fm.schema_group, 'default'),
							'schemaDescription', COALESCE(fm.schema_description, ''),
							'isObjArray', COALESCE(fm.is_obj_array, false),
							'fields', COALESCE(fm.fields_jsonb, '[]'::jsonb)
						)
					) AS metrics_jsonb
				FROM cte_fields_and_metrics fm
				JOIN eval_fns f ON f.eval_id = fm.eval_id
				WHERE f.org_id = $1 ` + addOnQuery + `
				GROUP BY f.eval_id, f.eval_name, f.eval_type, f.eval_group_name, f.eval_model, f.eval_format
			),
						cte_eval_fn_slice AS (
							SELECT 
								ct.eval_id,
							 ct.metrics_jsonb AS eval_fn_metrics_jsonb
							FROM 
								eval_fns_with_metrics ct
							JOIN eval_fns evf ON evf.eval_id = ct.eval_id
							GROUP BY 
							ct.eval_id, ct.metrics_jsonb
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
								em.eval_name, 
								em.eval_type, 
								em.eval_group_name, 
								em.eval_model, 
								em.eval_format,
				ct.triggers_list AS triggers_list, 
				COALESCE(cs.eval_fn_metrics_jsonb, '[]'::jsonb) AS json_schemas
			FROM eval_fns_with_metrics em
			LEFT JOIN cte_triggers ct ON ct.eval_id = em.eval_id
			LEFT JOIN cte_eval_fn_slice cs ON cs.eval_id = em.eval_id
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
		err = rows.Scan(
			&ef.EvalID,
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
		if ef.EvalID == nil {
			continue
		} else {
			ef.EvalStrID = aws.String(fmt.Sprintf("%d", *ef.EvalID))
		}
		for _, trigger := range dbTriggersHelper {
			if trigger.TriggerID == 0 {
				continue
			}
			ta := TriggerAction{
				TriggerID:     trigger.TriggerID,
				TriggerStrID:  fmt.Sprintf("%d", trigger.TriggerID),
				TriggerName:   trigger.TriggerName,
				TriggerGroup:  trigger.TriggerGroup,
				TriggerAction: trigger.TriggerAction,
				EvalTriggerAction: EvalTriggerActions{
					EvalID:               *ef.EvalID,
					EvalStrID:            fmt.Sprintf("%d", *ef.EvalID),
					TriggerID:            trigger.TriggerID,
					TriggerStrID:         fmt.Sprintf("%d", trigger.TriggerID),
					EvalTriggerState:     trigger.EvalTriggerState,
					EvalResultsTriggerOn: trigger.EvalResultsTriggerOn,
				},
				EvalTriggerActions: []EvalTriggerActions{
					{
						EvalID:               *ef.EvalID,
						EvalStrID:            fmt.Sprintf("%d", *ef.EvalID),
						TriggerID:            trigger.TriggerID,
						TriggerStrID:         fmt.Sprintf("%d", trigger.TriggerID),
						EvalTriggerState:     trigger.EvalTriggerState,
						EvalResultsTriggerOn: trigger.EvalResultsTriggerOn,
					},
				},
				TriggerActionsApprovals: nil,
			}
			ef.TriggerActions = append(ef.TriggerActions, ta)
		}

		var sc []*JsonSchemaDefinition
		for _, schema := range ef.Schemas {
			if schema.SchemaID == 0 || len(schema.Fields) <= 0 {
				continue
			}
			if schema.SchemaStrID == "" {
				schema.SchemaStrID = fmt.Sprintf("%d", schema.SchemaID)
			}
			for fi, _ := range schema.Fields {
				if schema.Fields[fi].FieldID == 0 {
					continue
				}
				if schema.Fields[fi].FieldStrID == "" {
					schema.Fields[fi].FieldStrID = fmt.Sprintf("%d", schema.Fields[fi].FieldID)
				}
			}
			if ef.SchemasMap == nil {
				ef.SchemasMap = make(map[string]*JsonSchemaDefinition)
			}
			ef.SchemasMap[strconv.Itoa(schema.SchemaID)] = schema
			sc = append(sc, schema)
		}
		ef.Schemas = sc
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
