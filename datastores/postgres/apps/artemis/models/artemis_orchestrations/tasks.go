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
	EvalFns               []EvalFn                `json:"evalFns"`
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
        cte_schemas AS (
            SELECT 
                ats.task_id,
                jsonb_agg(
                    jsonb_build_object(
                        'schemaID', jsd.schema_id,
                        'schemaName', jsd.schema_name,
                        'schemaGroup', jsd.schema_group,
                        'isObjArray', jsd.is_obj_array,
                        'fields', (
                            SELECT jsonb_agg(
                                jsonb_build_object(
                                    'fieldID', af.field_id,
                                    'fieldName', af.field_name,
                                    'fieldDescription', af.field_description,
                                    'dataType', af.data_type
                                )
                            ) 
                            FROM public.ai_fields af 
                            WHERE af.schema_id = jsd.schema_id AND af.is_field_archived = false
                        )
                    )
                ) AS schemas_jsonb
            FROM 
                cte_tasks ct
            JOIN 
                public.ai_task_schemas ats ON ct.task_id = ats.task_id
            JOIN 
                public.ai_json_schema_definitions jsd ON ats.schema_id = jsd.schema_id
            GROUP BY 
                ats.task_id
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
			cs.schemas_jsonb
		FROM 
			cte_tasks ct
		LEFT JOIN 
			cte_schemas cs ON ct.task_id = cs.task_id
		ORDER BY ct.task_id DESC 
    `
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
			&task.Prompt, &task.ResponseFormat, &task.Schemas, // Scan into a string
		)
		if err != nil {
			log.Err(err).Msg("failed to scan task")
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}
