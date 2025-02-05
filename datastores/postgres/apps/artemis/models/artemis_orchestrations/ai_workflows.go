package artemis_orchestrations

import (
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/artemis_entities"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

type WorkflowTemplate struct {
	WorkflowTemplateStrID     string `json:"workflowTemplateStrID,omitempty"`
	WorkflowTemplateID        int    `json:"workflowID,omitempty"`
	WorkflowName              string `json:"workflowName"`
	WorkflowGroup             string `json:"workflowGroup"`
	FundamentalPeriod         int    `json:"fundamentalPeriod"`
	FundamentalPeriodTimeUnit string `json:"fundamentalPeriodTimeUnit"`
	Tasks                     []Task `json:"tasks"` // Array of Task to hold the JSON aggregated tasks
}

type WorkflowUserEntitiesOverrides map[string][]artemis_entities.UserEntity

type WorkflowSchemaOverrides map[string]SchemaOverrides

type SchemaOverrides map[string]map[string][]string

type RetrievalOverrides map[string]RetrievalOverride

type RetrievalOverride struct {
	Payload  map[string]interface{}   `json:"retrievalPayload,omitempty"`
	Payloads []map[string]interface{} `json:"retrievalPayloads,omitempty"`
}

type TaskOverrides map[string]TaskOverride
type TaskOverride struct {
	ReplacePrompts  map[string]string `json:"replacePrompts,omitempty"`
	SystemPromptExt string            `json:"systemPromptExt,omitempty"`
}

type Task struct {
	TaskStrID         string     `json:"taskStrID,omitempty"`
	TaskID            int        `json:"taskID"`
	TaskName          string     `json:"taskName"`
	TaskType          string     `json:"taskType"`
	MarginBuffer      float64    `json:"marginBuffer,omitempty"`
	Temperature       float64    `json:"temperature,omitempty"`
	ResponseFormat    string     `json:"responseFormat"`
	Model             string     `json:"model"`
	Prompt            string     `json:"prompt"`
	CycleCount        int        `json:"cycleCount"`
	RetrievalName     string     `json:"retrievalName,omitempty"`
	RetrievalPlatform string     `json:"retrievalPlatform,omitempty"`
	EvalFnDBs         []EvalFnDB `json:"evalFns,omitempty"`
}

func InsertWorkflowTemplate(ctx context.Context, ou org_users.OrgUser, template *WorkflowTemplate) error {
	// SQL query to insert a new workflow template or update existing
	query := `
        INSERT INTO public.ai_workflow_template (workflow_name, workflow_group, org_id, user_id, fundamental_period, fundamental_period_time_unit)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, workflow_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            fundamental_period = EXCLUDED.fundamental_period,
            fundamental_period_time_unit = EXCLUDED.fundamental_period_time_unit
        RETURNING workflow_template_id;`
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, query, template.WorkflowName, template.WorkflowGroup, ou.OrgID, ou.UserID, template.FundamentalPeriod, template.FundamentalPeriodTimeUnit).Scan(&template.WorkflowTemplateID)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow template")
		return err
	}
	return nil
}

type AggTask struct {
	AggId      int             `json:"aggID"`
	CycleCount int             `json:"cycleCount"`
	EvalFns    []EvalFn        `json:"evalFns,omitempty"`
	Tasks      []AITaskLibrary `json:"tasks"`
}

type WorkflowTasks struct {
	AggTasks          []AggTask       `json:"aggTasks,omitempty"`
	AnalysisOnlyTasks []AITaskLibrary `json:"analysisOnlyTasks"`
}

func InsertWorkflowWithComponents(ctx context.Context, ou org_users.OrgUser, workflowTemplate *WorkflowTemplate, tasks WorkflowTasks) error {
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}

	// Rollback in case of any error
	defer tx.Rollback(ctx)
	// Insert the workflow template and get its ID
	query := `
        INSERT INTO public.ai_workflow_template (workflow_name, workflow_group, org_id, user_id, fundamental_period, fundamental_period_time_unit)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (org_id, workflow_name) 
        DO UPDATE SET 
            user_id = EXCLUDED.user_id,
            fundamental_period = EXCLUDED.fundamental_period,
            fundamental_period_time_unit = EXCLUDED.fundamental_period_time_unit
        RETURNING workflow_template_id;`

	err = tx.QueryRow(ctx, query,
		workflowTemplate.WorkflowName, workflowTemplate.WorkflowGroup, ou.OrgID, ou.UserID, workflowTemplate.FundamentalPeriod, workflowTemplate.FundamentalPeriodTimeUnit).Scan(&workflowTemplate.WorkflowTemplateID)
	if err != nil {
		log.Err(err).Msg("failed to insert workflow template")
		return err
	}
	for i, at := range tasks.AnalysisOnlyTasks {
		// Link component to the workflow template
		if at.TaskID == 0 && at.TaskStrID != "" {
			at.TaskID, err = strconv.Atoi(at.TaskStrID)
			if err != nil {
				log.Err(err).Msg("failed to parse int")
				return err
			}
			tasks.AnalysisOnlyTasks[i].TaskID = at.TaskID
		}
		for ir, rd := range at.RetrievalDependencies {
			if aws.ToInt(rd.RetrievalID) == 0 && aws.ToString(rd.RetrievalStrID) != "" {
				rid, rerr := strconv.Atoi(aws.ToString(rd.RetrievalStrID))
				if rerr != nil {
					log.Err(rerr).Msg("failed to parse int")
					return rerr
				}
				rd.RetrievalID = aws.Int(rid)
				at.RetrievalDependencies[ir].RetrievalID = aws.Int(rid)
			}
			err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_analysis_tasks(workflow_template_id, task_id, retrieval_id, cycle_count)
										VALUES ($1, $2, $3, $4)
										ON CONFLICT (workflow_template_id, task_id, retrieval_id)
										DO UPDATE SET cycle_count = EXCLUDED.cycle_count										
										RETURNING task_id`, workflowTemplate.WorkflowTemplateID, at.TaskID, rd.RetrievalID, at.CycleCount).Scan(&at.TaskID)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
		}
		if at.RetrievalDependencies == nil || len(at.RetrievalDependencies) <= 0 {
			err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_analysis_tasks(workflow_template_id, task_id, retrieval_id, cycle_count)
										VALUES ($1, $2, $3, $4)
										ON CONFLICT (workflow_template_id, task_id, retrieval_id)
										DO UPDATE SET cycle_count = EXCLUDED.cycle_count										
										RETURNING task_id`, workflowTemplate.WorkflowTemplateID, at.TaskID, nil, at.CycleCount).Scan(&at.TaskID)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}
		}

		// Link analysis tasks to eval functions
		for _, ef := range at.EvalFns {
			if ef.EvalCycleCount < 1 {
				ef.EvalCycleCount = 1
			}
			if ef.EvalStrID != nil && *ef.EvalStrID != "" {
				eid, eerr := strconv.Atoi(*ef.EvalStrID)
				if eerr != nil {
					log.Err(eerr).Msg("failed to parse int")
					return eerr
				}
				ef.EvalID = aws.Int(eid)
			}

			var taskEvalID int
			ts := chronos.Chronos{}
			err = tx.QueryRow(ctx, `INSERT INTO public.ai_workflow_template_eval_task_relationships(task_eval_id, workflow_template_id, task_id, cycle_count, eval_id)
                                    VALUES ($1, $2, $3, $4, $5)
                                    ON CONFLICT (workflow_template_id, task_id, eval_id)
                                    DO UPDATE SET 
                                        cycle_count = EXCLUDED.cycle_count
                                    RETURNING task_eval_id`,
				ts.UnixTimeStampNow(), workflowTemplate.WorkflowTemplateID, at.TaskID, ef.EvalCycleCount, ef.EvalID).Scan(&taskEvalID)
			if err != nil {
				log.Err(err).Interface("ef", ef).Msg("failed to insert eval task relationship for analysis task")
				return err
			}
		}
	}
	for _, aggTask := range tasks.AggTasks {
		// Link aggregation tasks to eval functions

		for ai, aggt := range aggTask.Tasks {
			if aggt.TaskID == 0 && aggt.TaskStrID != "" {
				aggt.TaskID, err = strconv.Atoi(aggt.TaskStrID)
				if err != nil {
					log.Err(err).Msg("failed to parse int")
					return err
				}
				aggTask.Tasks[ai].TaskID = aggt.TaskID
			}
			err = tx.QueryRow(ctx, `INSERT INTO ai_workflow_template_agg_tasks(agg_task_id, workflow_template_id, analysis_task_id, cycle_count)
											VALUES ($1, $2, $3, $4)
											ON CONFLICT (workflow_template_id, agg_task_id, analysis_task_id)
											DO UPDATE SET cycle_count = EXCLUDED.cycle_count
											RETURNING analysis_task_id`,
				aggTask.AggId, workflowTemplate.WorkflowTemplateID, aggt.TaskID, aggTask.CycleCount).Scan(&aggt.TaskID)
			if err != nil {
				log.Err(err).Msg("failed to insert workflow component")
				return err
			}

			for _, ef := range aggTask.EvalFns {
				if ef.EvalCycleCount < 1 {
					ef.EvalCycleCount = 1
				}
				if ef.EvalStrID != nil && *ef.EvalStrID != "" {
					eid, eerr := strconv.Atoi(*ef.EvalStrID)
					if eerr != nil {
						log.Err(eerr).Msg("failed to parse int")
						return eerr
					}
					ef.EvalID = aws.Int(eid)
				}
				var taskEvalID int
				ts := chronos.Chronos{}
				err = tx.QueryRow(ctx, `INSERT INTO public.ai_workflow_template_eval_task_relationships(task_eval_id, workflow_template_id, task_id, cycle_count, eval_id)
                                    VALUES ($1, $2, $3, $4, $5)
                                    ON CONFLICT (workflow_template_id, task_id, eval_id)
                                    DO UPDATE SET 
                                        cycle_count = EXCLUDED.cycle_count
                                    RETURNING task_eval_id`,
					ts.UnixTimeStampNow(), workflowTemplate.WorkflowTemplateID, aggTask.AggId, ef.EvalCycleCount, ef.EvalID).Scan(&taskEvalID)
				if err != nil {
					log.Err(err).Interface("ef", ef).Msg("failed to insert or update eval task relationship for aggregation task")
					return err
				}
			}
			// Link aggregation tasks to eval functions
			for _, ef := range aggt.EvalFns {
				if ef.EvalStrID != nil && *ef.EvalStrID != "" {
					eid, eerr := strconv.Atoi(*ef.EvalStrID)
					if eerr != nil {
						log.Err(eerr).Msg("failed to parse int")
						return eerr
					}
					ef.EvalID = aws.Int(eid)
				}
				var taskEvalID int
				ts := chronos.Chronos{}
				err = tx.QueryRow(ctx, `INSERT INTO public.ai_workflow_template_eval_task_relationships(task_eval_id, workflow_template_id, task_id, cycle_count, eval_id)
                                    VALUES ($1, $2, $3, $4, $5)
                                    ON CONFLICT (workflow_template_id, task_id, eval_id)
                                    DO UPDATE SET 
                                        cycle_count = EXCLUDED.cycle_count
                                    RETURNING task_eval_id`,
					ts.UnixTimeStampNow(), workflowTemplate.WorkflowTemplateID, aggt.TaskID, aggt.CycleCount, ef.EvalID).Scan(&taskEvalID)
				if err != nil {
					log.Err(err).Interface("ef", ef).Msg("failed to insert or update eval task relationship for analysis-aggregation task")
					return err
				}
			}
		}
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		log.Err(err).Msg("failed to commit transaction")
		return err
	}
	return nil
}
