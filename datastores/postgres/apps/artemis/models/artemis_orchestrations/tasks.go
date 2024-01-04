package artemis_orchestrations

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type AITaskLibrary struct {
	TaskID                int             `db:"task_id" json:"taskID,omitempty"`
	OrgID                 int             `db:"org_id" json:"orgID,omitempty"`
	UserID                int             `db:"user_id" json:"userID,omitempty"`
	MaxTokensPerTask      int             `db:"max_tokens_per_task" json:"maxTokensPerTask"`
	TaskType              string          `db:"task_type" json:"taskType"`
	TaskName              string          `db:"task_name" json:"taskName"`
	TaskGroup             string          `db:"task_group" json:"taskGroup"`
	TokenOverflowStrategy string          `db:"token_overflow_strategy" json:"tokenOverflowStrategy"`
	Model                 string          `db:"model" json:"model"`
	Prompt                string          `db:"prompt" json:"prompt"`
	ResponseFormat        string          `db:"response_format" json:"responseFormat"`
	CycleCount            int             `db:"cycle_count" json:"cycleCount,omitempty"`
	RetrievalDependencies []RetrievalItem `json:"retrievalDependencies,omitempty"`
	EvalFns               []EvalFn        `json:"evalFns"`
}

func InsertTask(ctx context.Context, task *AITaskLibrary) error {
	if task == nil {
		return nil
	}
	query := `
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
        RETURNING task_id;`
	// Executing the query
	if task.ResponseFormat == "" {
		task.ResponseFormat = "text"
	}

	err := apps.Pg.QueryRowWArgs(ctx, query,
		task.OrgID, task.UserID, task.MaxTokensPerTask, task.TaskType, task.TaskName, task.TaskGroup, task.TokenOverflowStrategy, task.Model, task.Prompt, task.ResponseFormat).
		Scan(&task.TaskID)
	if err != nil {
		log.Err(err).Msg("failed to insert task")
		return err
	}
	return nil
}

func SelectTasks(ctx context.Context, orgID int) ([]AITaskLibrary, error) {
	query := `
        SELECT task_id, org_id, user_id, max_tokens_per_task, task_type, task_name, task_group, token_overflow_strategy, model, prompt, response_format
        FROM public.ai_task_library
        WHERE org_id = $1
        ORDER BY task_id, task_group, task_name DESC;`

	// Executing the query
	rows, err := apps.Pg.Query(ctx, query, orgID)
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
		err = rows.Scan(&task.TaskID, &task.OrgID, &task.UserID, &task.MaxTokensPerTask, &task.TaskType, &task.TaskName,
			&task.TaskGroup, &task.TokenOverflowStrategy, &task.Model, &task.Prompt, &task.ResponseFormat)
		if err != nil {
			log.Err(err).Msg("failed to scan task")
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
