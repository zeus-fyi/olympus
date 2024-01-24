package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type AITaskLibrary struct {
	TaskID                int                     `db:"task_id" json:"taskID,omitempty"`
	OrgID                 int                     `db:"org_id" json:"orgID,omitempty"`
	UserID                int                     `db:"user_id" json:"userID,omitempty"`
	MaxTokensPerTask      int                     `db:"max_tokens_per_task" json:"maxTokensPerTask"`
	TaskType              string                  `db:"task_type" json:"taskType"`
	TaskName              string                  `db:"task_name" json:"taskName"`
	TaskGroup             string                  `db:"task_group" json:"taskGroup"`
	TokenOverflowStrategy string                  `db:"token_overflow_strategy" json:"tokenOverflowStrategy"`
	Model                 string                  `db:"model" json:"model"`
	Prompt                string                  `db:"prompt" json:"prompt"`
	Schemas               []*JsonSchemaDefinition `json:"schemas,omitempty"`
	ResponseFormat        string                  `db:"response_format" json:"responseFormat"`
	CycleCount            int                     `db:"cycle_count" json:"cycleCount,omitempty"`
	RetrievalDependencies []RetrievalItem         `json:"retrievalDependencies,omitempty"`
	EvalFns               []EvalFn                `json:"evalFns,omitempty"`
}

func InsertTask(ctx context.Context, task *AITaskLibrary) error {
	if task == nil {
		return nil
	}

	opt1 := `DELETE FROM public.ai_task_schemas
    		 WHERE task_id IN (SELECT task_id FROM cte_task_wrapper)`

	opt2 := `INSERT INTO public.ai_task_schemas (task_id, schema_id)
			 SELECT (SELECT task_id FROM cte_task_wrapper), unnest($11::bigint[])
			 ON CONFLICT (schema_id, task_id) DO NOTHING`

	// Executing the query
	if task.ResponseFormat == "" {
		task.ResponseFormat = "text"
		task.Schemas = nil
	}
	var sids []int
	if task.ResponseFormat == "json" || task.ResponseFormat == "social-media-engagement" {
		if len(task.Schemas) == 0 {
			return nil
		}
		for _, schema := range task.Schemas {
			sids = append(sids, schema.SchemaID)
		}
		opt1 = opt2
	}
	query := `
		WITH cte_task_wrapper AS (
			INSERT INTO public.ai_task_library 
				(org_id, user_id, max_tokens_per_task, task_type, task_name, task_group, token_overflow_strategy, model, prompt, response_format)
			VALUES 
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (org_id, task_group, task_name) 
			DO UPDATE SET 
				max_tokens_per_task = EXCLUDED.max_tokens_per_task,
				token_overflow_strategy = EXCLUDED.token_overflow_strategy,
				model = EXCLUDED.model,
				response_format = EXCLUDED.response_format,
				prompt = EXCLUDED.prompt
			RETURNING task_id
		), cte_cleanup_json_schema AS (
			DELETE FROM public.ai_task_schemas
			WHERE task_id IN (SELECT task_id FROM cte_task_wrapper) AND schema_id != ANY($11)
		), cte_insert_json_schema AS (
			` + opt1 + `
		) SELECT task_id FROM cte_task_wrapper;`

	err := apps.Pg.QueryRowWArgs(ctx, query,
		task.OrgID, task.UserID, task.MaxTokensPerTask, task.TaskType, task.TaskName, task.TaskGroup, task.TokenOverflowStrategy, task.Model, task.Prompt, task.ResponseFormat, pq.Array(sids)).
		Scan(&task.TaskID)
	if err != nil {
		log.Err(err).Msg("failed to insert task")
		return err
	}
	return nil
}

func SelectTasks(ctx context.Context, ou org_users.OrgUser) ([]AITaskLibrary, error) {
	return SelectTask(ctx, ou, 0)
}

func SelectTask(ctx context.Context, ou org_users.OrgUser, taskID int) ([]AITaskLibrary, error) {
	queryAddOn := ""
	params := []interface{}{ou.OrgID}
	if taskID != 0 {
		params = append(params, taskID)
		queryAddOn = "AND tl.task_id = $2"
	}
	query := `
			WITH cte_tasks AS (
				SELECT 
					tl.task_id, 
					tl.max_tokens_per_task, 
					tl.task_type,
					tl.task_name, 
					tl.task_group, 
					tl.token_overflow_strategy, 
					tl.model, 
					tl.prompt, 
					tl.response_format
				FROM 
					public.ai_task_library tl
        		WHERE 
            		tl.org_id = $1 ` + queryAddOn + ` 
			),
			cte_tasks_0 AS (
				SELECT 
					tl.task_id, 
					tl.max_tokens_per_task, 
					tl.task_type,
					tl.task_name, 
					tl.task_group, 
					tl.token_overflow_strategy, 
					tl.model, 
					tl.prompt, 
					tl.response_format,
					ate.eval_id
				FROM  
					cte_tasks tl
				LEFT JOIN ai_workflow_template_eval_task_relationships ate ON tl.task_id = ate.task_id
				GROUP BY 
					tl.task_id, tl.max_tokens_per_task, tl.task_type, tl.task_name, tl.task_group,
					tl.token_overflow_strategy, tl.model, tl.prompt, tl.response_format, ate.eval_id
			),
			cte_1 AS (
				SELECT 
					tl.task_id,
					tl.max_tokens_per_task, 
					tl.task_type,
					tl.task_name, 
					tl.task_group, 
					tl.token_overflow_strategy, 
					tl.model, 
					tl.prompt, 
					tl.response_format,
					tl.eval_id,
					af.field_id,
					jsd.schema_id,
					jsd.is_obj_array,
					jsd.schema_name,
					jsd.schema_group,
					jsd.schema_description,
					af.field_name, 
					af.field_description,
					af.data_type
				FROM 
					cte_tasks_0 tl
				JOIN ai_task_schemas ats ON ats.task_id = tl.task_id
				JOIN public.ai_json_schema_definitions jsd ON ats.schema_id = jsd.schema_id
				JOIN public.ai_fields af ON ats.schema_id = af.schema_id AND af.is_field_archived = false
				GROUP BY 
					tl.task_id, 
					tl.max_tokens_per_task, 
					tl.task_type,
					tl.task_name, 
					tl.task_group, 
					tl.token_overflow_strategy, 
					tl.model,
					tl.response_format,
					tl.prompt,
					tl.eval_id,
					jsd.schema_id, jsd.schema_name, jsd.schema_group, jsd.is_obj_array, jsd.schema_description,
					af.field_id
			),
			cte_2 AS (
				SELECT
					tl.task_id,
					tl.eval_id,
					tl.field_id,
					tl.field_name,
					tl.field_description,
					tl.data_type,
					tl.schema_id,
					tl.is_obj_array,
					tl.schema_name,
					tl.schema_group,
					tl.schema_description,
					jsonb_agg(
						jsonb_build_object(
							'evalMetricID', evm.eval_metric_id,
							'evalComparisonNumber', evm.eval_comparison_number,
							'evalComparisonBoolean', evm.eval_comparison_boolean,
							'evalComparisonString', evm.eval_comparison_string,
							'evalOperator', evm.eval_operator,
							'evalState', evm.eval_state,
							'evalExpectedResultState', evm.eval_metric_result
					)) AS eval_metrics_jsonb
				FROM  
					cte_1 tl
					LEFT JOIN public.eval_metrics evm ON evm.eval_id = tl.eval_id AND evm.is_eval_metric_archived = false AND evm.eval_state != 'ignore'
				GROUP BY 
					tl.task_id, tl.eval_id, tl.schema_id, tl.field_id, tl.field_name,
					tl.field_description, tl.data_type, tl.schema_name, tl.schema_group, tl.is_obj_array, tl.schema_description
			),
			cte_schema_definitions AS (
				SELECT 
					te.task_id,
					te.eval_id,
					te.schema_id,
					te.is_obj_array,
					te.schema_name,
					te.schema_group,
					te.schema_description,
					jsonb_agg(
						jsonb_build_object(
							'fieldID', te.field_id,
							'fieldName', te.field_name,
							'fieldDescription', te.field_description,
							'dataType', te.data_type,
							'evalMetrics', COALESCE(te.eval_metrics_jsonb, '[]'::jsonb)
						)
					) AS eval_fn_metrics_schemas_jsonb
				FROM 
					cte_2 te
				GROUP BY 
					te.task_id, te.eval_id, te.schema_id, te.schema_name, te.schema_group, te.schema_description, te.is_obj_array
			),
			cte_xy AS (
				SELECT 
					ct.task_id,
					ct.eval_id,
					ct.schema_id,
					jsonb_build_object(
						'schemaID', ct.schema_id,
						'schemaName', ct.schema_name,
						'schemaGroup', ct.schema_group,
						'schemaDescription', ct.schema_description,
						'isObjArray', ct.is_obj_array,
						'fields', ct.eval_fn_metrics_schemas_jsonb 
					) AS task_schemas_jsonb
				FROM 
					cte_schema_definitions ct
				GROUP BY 
					ct.task_id, ct.eval_id, ct.schema_id, ct.schema_name, ct.schema_group, ct.schema_description, ct.is_obj_array, ct.eval_fn_metrics_schemas_jsonb
			),
			cte_eval_fn_slice AS (
				SELECT 
					ct.task_id,
					ct.eval_id,
					jsonb_build_object(
						'evalID', ct.eval_id,
						'evalName', evf.eval_name,
						'evalType', evf.eval_type,
						'evalGroupName', evf.eval_group_name,
						'evalModel', evf.eval_model,
						'evalFormat', evf.eval_format,
						'schemas', jsonb_agg(ct.task_schemas_jsonb)
					) AS eval_fn_metrics_jsonb
				FROM 
					cte_xy ct
				JOIN eval_fns evf ON evf.eval_id = ct.eval_id
				GROUP BY 
					ct.task_id, ct.eval_id, evf.eval_name, evf.eval_group_name, evf.eval_type, evf.eval_model, evf.eval_format
			),
			cte_abc AS (
				SELECT 
					ct.task_id,
					jsonb_agg(eval_fn_metrics_jsonb) AS eval_fn_metrics_jsonb
				FROM 
					cte_eval_fn_slice ct
				GROUP BY 
					ct.task_id
			),
			cte_task_schema_definitions AS (
				SELECT 
					te.task_id,
					jsd.schema_id,
					jsd.schema_name,
					jsd.schema_group,
					jsd.schema_description,
					jsd.is_obj_array,
					jsonb_agg(
						jsonb_build_object(
							'fieldID', af.field_id,
							'fieldName', af.field_name,
							'fieldDescription', af.field_description,
							'dataType', af.data_type
						)
					) AS fields
				FROM 
					cte_tasks te
					JOIN ai_task_schemas ats ON ats.task_id = te.task_id
					JOIN public.ai_json_schema_definitions jsd ON ats.schema_id = jsd.schema_id
					JOIN public.ai_fields af ON ats.schema_id = af.schema_id AND af.is_field_archived = false
				GROUP BY 
					te.task_id, jsd.schema_id, jsd.schema_name, jsd.schema_group, jsd.schema_description, jsd.is_obj_array
			),
			cte_task_schemas AS (
				SELECT 
					ct.task_id,
					ct.schema_id,
					jsonb_build_object(
						'schemaID', ct.schema_id,
						'schemaName', ct.schema_name,
						'schemaGroup', ct.schema_group,
						'schemaDescription', ct.schema_description,
						'isObjArray', ct.is_obj_array,
						'fields', ct.fields
					) AS task_schemas_jsonb
				FROM 
					cte_task_schema_definitions ct
				GROUP BY 
					ct.task_id, ct.schema_id, ct.schema_name, ct.schema_group, ct.schema_description, ct.is_obj_array, ct.fields
			),
			cte_x AS (
				SELECT 
					ct.task_id,
					jsonb_agg(ct.task_schemas_jsonb) AS task_schemas_jsonb
				FROM 
					cte_task_schemas ct
				GROUP BY 
					ct.task_id
			)
			SELECT 
				ct.task_id, 
				ct.max_tokens_per_task, 
				ct.task_type,
				ct.task_name, 
				ct.task_group,
				ct.token_overflow_strategy, 
				ct.model, 
				ct.prompt, 
				ct.response_format,
				COALESCE(css.task_schemas_jsonb, '[]'::jsonb) AS task_schemas, 
				COALESCE(cefs.eval_fn_metrics_jsonb, '[]'::jsonb) AS eval_fns
			FROM 
				cte_tasks ct
				LEFT JOIN cte_abc cefs ON ct.task_id = cefs.task_id
				LEFT JOIN cte_x css ON ct.task_id = css.task_id  
			ORDER BY 
				ct.task_id DESC;`
	// Executing the query

	rows, err := apps.Pg.Query(ctx, query, params...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	var tasks []AITaskLibrary

	// Iterating over the result set
	for rows.Next() {
		var task AITaskLibrary

		err = rows.Scan(
			&task.TaskID, &task.MaxTokensPerTask, &task.TaskType, &task.TaskName,
			&task.TaskGroup, &task.TokenOverflowStrategy, &task.Model,
			&task.Prompt, &task.ResponseFormat, &task.Schemas, &task.EvalFns, // Scan into a string
		)
		if err != nil {
			log.Err(err).Msg("failed to scan task")
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}
