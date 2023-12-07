package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowTemplateData struct {
	AnalysisTaskID                int     `json:"analysisTaskID"`
	AnalysisCycleCount            int     `json:"analysisCycleCount"`
	AnalysisPrompt                string  `json:"analysisPrompt"`
	AnalysisModel                 string  `json:"analysisModel"`
	AnalysisTokenOverflowStrategy string  `json:"analysisTokenOverflowStrategy"`
	AnalysisTaskName              string  `json:"analysisTaskName"`
	AnalysisTaskType              string  `json:"analysisTaskType"`
	AnalysisMaxTokensPerTask      int     `json:"analysisMaxTokensPerTask"`
	AggTaskID                     *int    `json:"aggTaskID,omitempty"`
	AggCycleCount                 *int    `json:"aggCycleCount,omitempty"`
	AggTaskName                   *string `json:"aggTaskName,omitempty"`
	AggTaskType                   *string `json:"aggTaskType,omitempty"`
	AggPrompt                     *string `json:"aggPrompt,omitempty"`
	AggModel                      *string `json:"aggModel,omitempty"`
	AggTokenOverflowStrategy      *string `json:"aggTokenOverflowStrategy,omitempty"`
	AggMaxTokensPerTask           *int    `json:"aggMaxTokensPerTask,omitempty"`
	RetrievalID                   *int    `json:"retrievalID,omitempty"`
	RetrievalName                 *string `json:"retrievalName,omitempty"`
	RetrievalGroup                *string `json:"retrievalGroup,omitempty"`
	RetrievalPlatform             *string `json:"retrievalPlatform,omitempty"`
	RetrievalInstructions         []byte  `json:"retrievalInstructions,omitempty"`
}

func SelectWorkflowTemplate(ctx context.Context, ou org_users.OrgUser, workflowName string) ([]WorkflowTemplateData, error) {
	var results []WorkflowTemplateData
	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID, workflowName}
	q.RawQuery = `WITH cte_1 AS (
						SELECT
							awtat.task_id AS analysis_task_id,
							awtat.cycle_count AS analysis_cycle_count,
							ait.task_name AS analysis_task_name,
							ait.task_type AS analysis_task_type,
							ait.prompt AS analysis_prompt, 
							ait.model AS analysis_model,
							ait.token_overflow_strategy AS analysis_token_overflow_strategy,
							ait.max_tokens_per_task AS analysis_max_tokens_per_task,
							awtat.retrieval_id AS retrieval_id
						FROM ai_workflow_template wate
						JOIN public.ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
						JOIN public.ai_task_library ait ON ait.task_id = awtat.task_id
						WHERE wate.org_id = $1 AND wate.workflow_name = $2
						), cte_2 AS (
							SELECT
								awtat.agg_task_id,
								awtat.analysis_task_id,
								awtat.cycle_count AS agg_cycle_count,
								ait.task_name AS agg_task_name,
								ait.task_type AS agg_task_type,
								ait.prompt AS agg_prompt,
								ait.model AS agg_model,
								ait.token_overflow_strategy AS agg_token_overflow_strategy,
								ait.max_tokens_per_task AS agg_max_tokens_per_task
							FROM ai_workflow_template wate
							JOIN public.ai_workflow_template_agg_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
							JOIN public.ai_task_library ait ON ait.task_id = awtat.agg_task_id
							JOIN public.ai_task_library ait1 ON ait1.task_id = awtat.analysis_task_id
							WHERE wate.org_id = $1 AND wate.workflow_name = $2
							), cte_3 AS (
								SELECT  art.retrieval_id,
											art.retrieval_name,
											art.retrieval_group,
											art.retrieval_platform as retrieval_platform,
											art.instructions as retrieval_instructions
								FROM cte_1 c1 
								JOIN ai_retrieval_library art ON art.retrieval_id = c1.retrieval_id
						), aggregate_cte AS (
						SELECT
							cte_1.analysis_task_id,
							cte_1.analysis_cycle_count,
							cte_1.analysis_task_name,
							cte_1.analysis_task_type,
							cte_1.analysis_prompt,
							cte_1.analysis_model,
							cte_1.analysis_token_overflow_strategy,
							cte_1.analysis_max_tokens_per_task,
							cte_2.agg_task_id,
							cte_2.analysis_task_id,
							cte_2.agg_task_name,
							cte_2.agg_task_type,
							cte_2.agg_prompt,
							cte_2.agg_model,
							cte_2.agg_token_overflow_strategy,
							cte_2.agg_max_tokens_per_task,
							cte_2.agg_cycle_count,
							cte_3.retrieval_id,
							cte_3.retrieval_name,
							cte_3.retrieval_group,
							cte_3.retrieval_platform,
							cte_3.retrieval_instructions
						FROM cte_1 
						LEFT JOIN cte_2 ON cte_1.analysis_task_id = cte_2.analysis_task_id
						LEFT JOIN cte_3 ON cte_1.retrieval_id = cte_3.retrieval_id
						) 
						SELECT * FROM aggregate_cte`
	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		log.Err(err).Msg("Error querying SelectWorkflowTemplate")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var aggAnalysisTaskID *int
		var data WorkflowTemplateData
		rowErr := rows.Scan(
			&data.AnalysisTaskID,
			&data.AnalysisCycleCount,
			&data.AnalysisTaskName,
			&data.AnalysisTaskType,
			&data.AnalysisPrompt,
			&data.AnalysisModel,
			&data.AnalysisTokenOverflowStrategy,
			&data.AnalysisMaxTokensPerTask,
			&data.AggTaskID,
			&aggAnalysisTaskID,
			&data.AggTaskName,
			&data.AggTaskType,
			&data.AggPrompt,
			&data.AggModel,
			&data.AggTokenOverflowStrategy,
			&data.AggMaxTokensPerTask,
			&data.AggCycleCount,
			&data.RetrievalID,
			&data.RetrievalName,
			&data.RetrievalGroup,
			&data.RetrievalPlatform,
			&data.RetrievalInstructions,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg("Error scanning row in SelectWorkflowTemplate")
			return nil, rowErr
		}
		results = append(results, data)
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating over rows")
		return nil, err
	}

	return results, nil
}

type AggTaskDb struct {
	AggModel                 string `json:"agg_model"`
	AggPrompt                string `json:"agg_prompt"`
	AggTaskId                int    `json:"agg_task_id"`
	AggTaskName              string `json:"agg_task_name"`
	AggTaskType              string `json:"agg_task_type"`
	AggCycleCount            int    `json:"agg_cycle_count"`
	AggAnalysisTaskId        int    `json:"agg_analysis_task_id"`
	AggMaxTokensPerTask      int    `json:"agg_max_tokens_per_task"`
	AggTokenOverflowStrategy string `json:"agg_token_overflow_strategy"`
}
type AnalysisTaskDB struct {
	AnalysisModel                 string `json:"analysis_model"`
	AnalysisPrompt                string `json:"analysis_prompt"`
	AnalysisTaskId                int    `json:"analysis_task_id"`
	AnalysisTaskName              string `json:"analysis_task_name"`
	AnalysisTaskType              string `json:"analysis_task_type"`
	AnalysisMaxTokensPerTask      int    `json:"analysis_max_tokens_per_task"`
	AnalysisTokenOverflowStrategy string `json:"analysis_token_overflow_strategy"`
	RetrievalDB
}

type RetrievalDB struct {
	RetrievalId           int    `json:"retrieval_id"`
	RetrievalName         string `json:"retrieval_name"`
	RetrievalGroup        string `json:"retrieval_group"`
	RetrievalPlatform     string `json:"retrieval_platform"`
	AnalysisCycleCount    int    `json:"analysis_cycle_count"`
	RetrievalInstructions struct {
		RetrievalPlatform string `json:"retrievalPlatform"`
	} `json:"retrieval_instructions"`
}

func SelectWorkflowTemplates(ctx context.Context, ou org_users.OrgUser) (*Workflows, error) {
	results := &Workflows{
		WorkflowTemplatesMap: make(map[int]WorkflowTemplateMapValue),
	}
	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID}
	q.RawQuery = `WITH cte_1 AS (
					SELECT
						wate.workflow_template_id,
						wate.workflow_name,
						wate.workflow_group,
						awtat.task_id AS analysis_task_id,
						awtat.cycle_count AS analysis_cycle_count,
						ait.task_name AS analysis_task_name,
						ait.task_type AS analysis_task_type,
						ait.prompt AS analysis_prompt, 
						ait.model AS analysis_model,
						ait.token_overflow_strategy AS analysis_token_overflow_strategy,
						ait.max_tokens_per_task AS analysis_max_tokens_per_task,
						awtat.retrieval_id AS retrieval_id
					FROM ai_workflow_template wate
					JOIN public.ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
					JOIN public.ai_task_library ait ON ait.task_id = awtat.task_id
					WHERE wate.org_id = $1
				), 
				cte_2 AS (
					SELECT
						awtat.agg_task_id,
						awtat.analysis_task_id AS agg_analysis_task_id,
						awtat.cycle_count AS agg_cycle_count,
						ait.task_name AS agg_task_name,
						ait.task_type AS agg_task_type,
						ait.prompt AS agg_prompt,
						ait.model AS agg_model,
						ait.token_overflow_strategy AS agg_token_overflow_strategy,
						ait.max_tokens_per_task AS agg_max_tokens_per_task
					FROM ai_workflow_template wate
					JOIN public.ai_workflow_template_agg_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
					JOIN public.ai_task_library ait ON ait.task_id = awtat.agg_task_id
					JOIN public.ai_task_library ait1 ON ait1.task_id = awtat.analysis_task_id
					WHERE wate.org_id = $1
				), 
				cte_3 AS (
					SELECT  
						art.retrieval_id,
						art.retrieval_name,
						art.retrieval_group,
						art.retrieval_platform as retrieval_platform,
						art.instructions as retrieval_instructions
					FROM cte_1 c1 
					JOIN ai_retrieval_library art ON art.retrieval_id = c1.retrieval_id
				)
				SELECT
					cte_1.workflow_template_id,
					cte_1.workflow_name,
					cte_1.workflow_group,
					jsonb_object_agg(
						cte_1.analysis_task_id,
						JSON_BUILD_OBJECT(
							'analysis_task_id', cte_1.analysis_task_id,
							'analysis_cycle_count', cte_1.analysis_cycle_count,
							'analysis_task_name', cte_1.analysis_task_name,
							'analysis_task_type', cte_1.analysis_task_type,
							'analysis_prompt', cte_1.analysis_prompt,
							'analysis_model', cte_1.analysis_model,
							'analysis_token_overflow_strategy', cte_1.analysis_token_overflow_strategy,
							'analysis_max_tokens_per_task', cte_1.analysis_max_tokens_per_task,
							'retrieval_id', cte_3.retrieval_id,
							'retrieval_name', cte_3.retrieval_name,
							'retrieval_group', cte_3.retrieval_group,
							'retrieval_platform', cte_3.retrieval_platform,
							'retrieval_instructions', cte_3.retrieval_instructions
						)
					) AS analysis_tasks,
						jsonb_object_agg(
							cte_2.agg_task_id,
							JSON_BUILD_OBJECT(
							'agg_task_id', cte_2.agg_task_id,
							'agg_analysis_task_id', cte_2.agg_analysis_task_id,
							'agg_task_name', cte_2.agg_task_name,
							'agg_task_type', cte_2.agg_task_type,
							'agg_prompt', cte_2.agg_prompt,
							'agg_model', cte_2.agg_model,
							'agg_token_overflow_strategy', cte_2.agg_token_overflow_strategy,
							'agg_max_tokens_per_task', cte_2.agg_max_tokens_per_task,
							'agg_cycle_count', cte_2.agg_cycle_count
						)
					) AS agg_tasks
				FROM cte_1 
				JOIN cte_2 ON cte_1.analysis_task_id = cte_2.agg_analysis_task_id
				JOIN cte_3 ON cte_1.retrieval_id = cte_3.retrieval_id
				GROUP BY cte_1.workflow_template_id, cte_1.workflow_name, cte_1.workflow_group, cte_1.analysis_task_id, cte_2.agg_task_id, cte_2.agg_analysis_task_id, cte_3.retrieval_id`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var taskJSON, aggTasksJSON string
		wt := WorkflowTemplateMapValue{
			AnalysisTasks:      make(map[int]AnalysisTaskDB),
			AnalysisRetrievals: make(map[int]map[int]RetrievalDB),
			AggTasks:           make(map[int]AggTaskDb),
			AggAnalysisTasks:   make(map[int]map[int]AggTaskDb),
		}
		err = rows.Scan(
			&wt.WorkflowTemplateID,
			&wt.WorkflowName,
			&wt.WorkflowGroup,
			&taskJSON,
			&aggTasksJSON,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		err = json.Unmarshal([]byte(taskJSON), &wt.AnalysisTasks)
		if err != nil {
			log.Printf("Error unmarshalling analysis tasks JSON: %v", err)
			return nil, err
		}

		err = json.Unmarshal([]byte(aggTasksJSON), &wt.AggTasks)
		if err != nil {
			log.Printf("Error unmarshalling aggregate tasks JSON: %v", err)
			return nil, err
		}
		if _, ok := results.WorkflowTemplatesMap[wt.WorkflowTemplateID]; !ok {
			results.WorkflowTemplatesMap[wt.WorkflowTemplateID] = wt
		}
		tmp := results.WorkflowTemplatesMap[wt.WorkflowTemplateID]
		for k, v := range wt.AnalysisTasks {
			tmp.AnalysisTasks[k] = v
			if v.RetrievalId > 0 {
				if _, ok := tmp.AnalysisRetrievals[v.AnalysisTaskId]; !ok {
					tmp.AnalysisRetrievals[v.AnalysisTaskId] = make(map[int]RetrievalDB)
				}
				tmp.AnalysisRetrievals[v.AnalysisTaskId][v.RetrievalId] = v.RetrievalDB
			}
		}

		for k, v := range wt.AggTasks {
			tmp.AggTasks[k] = v
			if tmp.AggTasks[k].AggAnalysisTaskId > 0 {
				if _, ok := tmp.AggAnalysisTasks[v.AggTaskId]; !ok {
					tmp.AggAnalysisTasks[v.AggTaskId] = make(map[int]AggTaskDb)
				}
				tmp.AggAnalysisTasks[v.AggTaskId][v.AggAnalysisTaskId] = v
			}
		}
		results.WorkflowTemplatesMap[wt.WorkflowTemplateID] = tmp
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating over rows")
		return nil, err
	}

	return results, nil
}

type Workflows struct {
	WorkflowTemplatesMap map[int]WorkflowTemplateMapValue `json:"templates"`
}

type WorkflowTemplateMapValue struct {
	WorkflowTemplateID int                         `json:"workflow_template_id"`
	WorkflowName       string                      `json:"workflow_name"`
	WorkflowGroup      string                      `json:"workflow_group"`
	AnalysisTasks      map[int]AnalysisTaskDB      `json:"analysis_tasks"`
	AnalysisRetrievals map[int]map[int]RetrievalDB `json:"analysis_retrievals"`
	AggTasks           map[int]AggTaskDb           `json:"agg_tasks"`
	AggAnalysisTasks   map[int]map[int]AggTaskDb   `json:"agg_analysis_tasks"`
}

type WorkflowTaskRelationships struct {
	AnalysisRetrievals map[int]map[int]bool
	AggregateAnalysis  map[int]map[int]bool
}

func MapDependencies(res []WorkflowTemplateData) WorkflowTaskRelationships {
	analysisRetrievals := make(map[int]map[int]bool)
	aggregateAnalysis := make(map[int]map[int]bool)

	for _, v := range res {
		if _, ok := analysisRetrievals[v.AnalysisTaskID]; !ok {
			analysisRetrievals[v.AnalysisTaskID] = make(map[int]bool)
		}
		if v.RetrievalID != nil {
			if _, ok := analysisRetrievals[v.AnalysisTaskID][*v.RetrievalID]; !ok {
				analysisRetrievals[v.AnalysisTaskID][*v.RetrievalID] = true
			} else {
				fmt.Println("Duplicate retrieval id", v.RetrievalID)
			}
		}

		if v.AggTaskID != nil {
			if _, ok := aggregateAnalysis[*v.AggTaskID]; !ok {
				aggregateAnalysis[*v.AggTaskID] = make(map[int]bool)
			}
			if _, ok := aggregateAnalysis[*v.AggTaskID][v.AnalysisTaskID]; !ok {
				aggregateAnalysis[*v.AggTaskID][v.AnalysisTaskID] = true
			} else {
				fmt.Println("Duplicate agg id", *v.AggTaskID)
			}
		}
	}
	return WorkflowTaskRelationships{
		AnalysisRetrievals: analysisRetrievals,
		AggregateAnalysis:  aggregateAnalysis,
	}
}
func MapDependencies1(res WorkflowTemplateMapValue) WorkflowTaskRelationships {
	analysisRetrievals := make(map[int]map[int]bool)
	for _, analysisTask := range res.AnalysisTasks {
		analysisTaskID := analysisTask.AnalysisTaskId
		if _, ok := analysisRetrievals[analysisTaskID]; !ok {
			analysisRetrievals[analysisTaskID] = make(map[int]bool)
		}
	}
	for analysisTaskID, retrievalMap := range res.AnalysisRetrievals {
		for _, retrieval := range retrievalMap {
			analysisRetrievals[analysisTaskID][retrieval.RetrievalId] = true
		}
	}
	aggregateAnalysis := make(map[int]map[int]bool)
	for _, aggTask := range res.AggTasks {
		aggTaskID := aggTask.AggTaskId
		if _, ok := aggregateAnalysis[aggTaskID]; !ok {
			aggregateAnalysis[aggTaskID] = make(map[int]bool)
		}
	}
	for aggTaskID, analysisMap := range res.AggAnalysisTasks {
		for _, analysis := range analysisMap {
			aggregateAnalysis[aggTaskID][analysis.AggAnalysisTaskId] = true
		}
	}
	return WorkflowTaskRelationships{
		AnalysisRetrievals: analysisRetrievals,
		AggregateAnalysis:  aggregateAnalysis,
	}
}
