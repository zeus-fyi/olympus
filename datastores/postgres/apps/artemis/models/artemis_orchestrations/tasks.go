package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type AITaskLibrary struct {
	TaskID                int                    `db:"task_id" json:"taskID,omitempty"`
	OrgID                 int                    `db:"org_id" json:"orgID,omitempty"`
	UserID                int                    `db:"user_id" json:"userID,omitempty"`
	MaxTokensPerTask      int                    `db:"max_tokens_per_task" json:"maxTokensPerTask"`
	TaskType              string                 `db:"task_type" json:"taskType"`
	TaskName              string                 `db:"task_name" json:"taskName"`
	TaskGroup             string                 `db:"task_group" json:"taskGroup"`
	TokenOverflowStrategy string                 `db:"token_overflow_strategy" json:"tokenOverflowStrategy"`
	Model                 string                 `db:"model" json:"model"`
	Prompt                string                 `db:"prompt" json:"prompt"`
	Schemas               []JsonSchemaDefinition `json:"schemas,omitempty"`
	ResponseFormat        string                 `db:"response_format" json:"responseFormat"`
	CycleCount            int                    `db:"cycle_count" json:"cycleCount,omitempty"`
	RetrievalDependencies []RetrievalItem        `json:"retrievalDependencies,omitempty"`
	EvalFns               []EvalFn               `json:"evalFns"`
}

func InsertTask(ctx context.Context, task *AITaskLibrary) error {
	if task == nil {
		return nil
	}

	opt1 := `DELETE FROM public.ai_json_task_schemas
    		 WHERE task_id IN (SELECT task_id FROM cte_task_wrapper)`

	opt2 := `INSERT INTO public.ai_json_task_schemas (task_id, schema_id)
			 SELECT (SELECT task_id FROM cte_task_wrapper), unnest($11::bigint[])
			 ON CONFLICT (schema_id, task_id) DO NOTHING`

	// Executing the query
	if task.ResponseFormat == "" {
		task.ResponseFormat = "text"
		task.Schemas = nil
	}
	var sids []int
	if task.ResponseFormat == "json" {
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
			DELETE FROM public.ai_json_task_schemas
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
	query := `SELECT 
				tl.task_id, 
				tl.max_tokens_per_task, 
				tl.task_type,
				tl.task_name, 
				tl.task_group, 
				tl.token_overflow_strategy, 
				tl.model, 
				tl.prompt, 
				tl.response_format,
				array_agg(json_schema) AS json_schemas
			FROM 
				public.ai_task_library tl
			LEFT JOIN 
				public.ai_json_task_schemas js ON tl.task_id = js.task_id
			LEFT JOIN (
				SELECT 
					d.schema_id, 
					d.org_id, 
					jsonb_build_object(
						'schemaID', d.schema_id, 
						'schemaName', d.schema_name, 
						'schemaGroup', d.schema_group, 
						'isObjArray', d.is_obj_array,
						'fields', array_agg(
							jsonb_build_object(
								'fieldName', f.field_name, 
								'dataType', f.data_type, 
								'fieldDescription', f.field_description
							)
						)
					) AS json_schema
				FROM 
					public.ai_json_schema_definitions d
				JOIN 
					public.ai_json_schema_fields f ON d.schema_id = f.schema_id
				WHERE 
					d.org_id = $1
				GROUP BY 
					d.schema_id, d.org_id
			) AS schemas ON js.schema_id = schemas.schema_id AND tl.org_id = schemas.org_id
			WHERE 
				tl.org_id = $1
			GROUP BY 
				tl.task_id, 
				tl.max_tokens_per_task, 
				tl.task_type,
				tl.task_name, 
				tl.task_group,
				tl.token_overflow_strategy, 
				tl.model, 
				tl.prompt, 
				tl.response_format
			ORDER BY task_id DESC 
`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, query, ou.OrgID)
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
		var jsonSchema pgtype.JSONBArray

		err = rows.Scan(
			&task.TaskID, &task.MaxTokensPerTask, &task.TaskType, &task.TaskName,
			&task.TaskGroup, &task.TokenOverflowStrategy, &task.Model,
			&task.Prompt, &task.ResponseFormat, &jsonSchema, // Scan into a string
		)
		if err != nil {
			log.Err(err).Msg("failed to scan task")
			return nil, err
		}
		for _, elem := range jsonSchema.Elements {
			var schema JsonSchemaDefinition
			if elem.Bytes == nil {
				continue
			}
			if err = json.Unmarshal(elem.Bytes, &schema); err != nil {
				log.Err(err).Msg("failed to unmarshal json schema element")
				return nil, err
			}
			task.Schemas = append(task.Schemas, schema)
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}
