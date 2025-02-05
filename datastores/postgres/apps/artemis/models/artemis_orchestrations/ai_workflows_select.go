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

type Workflows struct {
	WorkflowTemplatesMap  map[int]WorkflowTemplateValue `json:"templates"`
	WorkflowTemplateSlice []WorkflowTemplateValue       `json:"templatesSlice"`
}

type WorkflowTemplateValue struct {
	WorkflowTemplate
	AnalysisTasks         map[int]AnalysisTaskDB      `json:"-"`
	AnalysisRetrievals    map[int]map[int]RetrievalDB `json:"-"`
	AnalysisEvalFns       map[int][]EvalFnDB          `json:"-"` // Mapping task ID to its evaluation functions
	AggTasks              map[int]AggTaskDb           `json:"-"`
	AggAnalysisTasks      map[int]map[int]AggTaskDb   `json:"-"`
	AnalysisTasksSlice    []AnalysisTaskDB            `json:"-"`
	AggAnalysisTasksSlice []AggTaskDb                 `json:"-"`
	AggEvalFns            map[int][]EvalFnDB          `json:"-"` // Mapping aggregated task ID to its evaluation functions
	AggAnalysisEvalFns    map[int]map[int]EvalFnDB    `json:"-"` // Mapping aggregated task ID to its evaluation functions
}

type WorkflowTemplateData struct {
	AnalysisTaskDB
	AnalysisMaxTokensPerTask int        `json:"analysisMaxTokensPerTask"`
	AggTaskID                *int       `json:"aggTaskID,omitempty"`
	AggCycleCount            *int       `json:"aggCycleCount,omitempty"`
	AggTaskName              *string    `json:"aggTaskName,omitempty"`
	AggTaskType              *string    `json:"aggTaskType,omitempty"`
	AggResponseFormat        *string    `json:"aggResponseFormat,omitempty"`
	AggPrompt                *string    `json:"aggPrompt,omitempty"`
	AggModel                 *string    `json:"aggModel,omitempty"`
	AggTemperature           *float64   `json:"aggTemperature,omitempty"`
	AggMarginBuffer          *float64   `json:"aggMarginBuffer,omitempty"`
	AggTokenOverflowStrategy *string    `json:"aggTokenOverflowStrategy,omitempty"`
	AggMaxTokensPerTask      *int       `json:"aggMaxTokensPerTask,omitempty"`
	AggAnalysisEvalFns       []EvalFnDB `json:"aggAnalysisEvalFns,omitempty"`
	AggEvalFns               []EvalFnDB `json:"aggEvalFns,omitempty"`
}

type AggTaskDb struct {
	AggModel                 string     `json:"aggModel"`
	AggPrompt                string     `json:"aggPrompt"`
	AggTaskId                int        `json:"aggTaskId"` // Assuming IDs are large numbers
	AggTaskName              string     `json:"aggTaskName"`
	AggResponseFormat        string     `json:"aggResponseFormat"`
	AggTaskType              string     `json:"aggTaskType"`
	AggTemperature           float64    `json:"aggTemperature"`
	AggMarginBuffer          float64    `json:"aggMarginBuffer"`
	AggCycleCount            int        `json:"aggCycleCount"`
	AggAnalysisTaskId        int        `json:"aggAnalysisTaskId"` // Assuming IDs are large numbers
	AggMaxTokensPerTask      int        `json:"aggMaxTokensPerTask"`
	AggTokenOverflowStrategy string     `json:"aggTokenOverflowStrategy"`
	AnalysisAggCycleCount    int        `json:"analysisAggCycleCount"`
	EvalFns                  []EvalFnDB `json:"evalFns"`            // Single object, not a slice
	AnalysisAggEvalFns       []EvalFnDB `json:"analysisAggEvalFns"` // Single object, not a slice
}

type EvalFnDB struct {
	EvalID         int    `json:"evalID"`
	EvalName       string `json:"evalName"`
	EvalType       string `json:"evalType"`
	EvalModel      string `json:"evalModel"`
	EvalFormat     string `json:"evalFormat"`
	EvalTaskID     int    `json:"evalTaskID"`
	EvalGroupName  string `json:"evalGroupName"`
	EvalCycleCount int    `json:"evalCycleCount"`
}

type AnalysisTaskDB struct {
	AnalysisModel                 string  `json:"analysisModel"`
	AnalysisCycleCount            int     `json:"analysisCycleCount"`
	AnalysisPrompt                string  `json:"analysisPrompt"`
	AnalysisTaskID                int     `json:"analysisTaskID"`
	AnalysisTemperature           float64 `json:"analysisTemperature"`
	AnalysisMarginBuffer          float64 `json:"analysisMarginBuffer"`
	AnalysisTaskName              string  `json:"analysisTaskName"`
	AnalysisResponseFormat        string  `json:"analysisResponseFormat"`
	AnalysisTaskType              string  `json:"analysisTaskType"`
	AnalysisMaxTokensPerTask      int     `json:"analysisMaxTokensPerTask"`
	AnalysisTokenOverflowStrategy string  `json:"analysisTokenOverflowStrategy"`
	RetrievalDB
	AnalysisEvalFns []EvalFnDB `json:"analysisEvalFns,omitempty"`
}

type RetrievalDB struct {
	RetrievalID           *int            `json:"retrievalID"`
	RetrievalName         string          `json:"retrievalName"`
	RetrievalGroup        string          `json:"retrievalGroup"`
	RetrievalPlatform     string          `json:"retrievalPlatform"`
	RetrievalInstructions json.RawMessage `json:"retrievalInstructions"`
}

func SelectWorkflowTemplates(ctx context.Context, ou org_users.OrgUser) (*Workflows, error) {
	return SelectWorkflowTemplateByName(ctx, ou, "")
}

func SelectWorkflowTemplateByName(ctx context.Context, ou org_users.OrgUser, name string) (*Workflows, error) {
	results := &Workflows{
		WorkflowTemplatesMap: make(map[int]WorkflowTemplateValue),
	}

	q := sql_query_templates.QueryParams{}
	params := []interface{}{ou.OrgID}
	additionalCondition := ""

	if len(name) > 0 {
		additionalCondition = " AND wate.workflow_name = $2"
		params = append(params, name)
	}
	q.RawQuery = `WITH cte_0 AS (
						SELECT	wate.workflow_template_id,
								wate.workflow_name,
								wate.workflow_group,
								awtat.task_id AS analysis_task_id,
								awtat.retrieval_id AS retrieval_id,
								awtat.cycle_count AS analysis_cycle_count
						FROM ai_workflow_template wate
						JOIN public.ai_workflow_template_analysis_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
						WHERE wate.org_id = $1 ` + additionalCondition + `
						GROUP BY wate.workflow_template_id, awtat.task_id, awtat.retrieval_id, awtat.cycle_count
					), cte_wf_evals AS (
						SELECT
							c0.workflow_template_id,
							evtr.task_id,
							JSON_BUILD_OBJECT(
								'evalID', ef.eval_id,
								'evalTaskID', evtr.task_id,
								'evalCycleCount', evtr.cycle_count,
								'evalName', ef.eval_name,
								'evalType', ef.eval_type,
								'evalGroupName', ef.eval_group_name,
								'evalModel', ef.eval_model,
								'evalFormat', ef.eval_format
							) AS eval_fns_data
						FROM cte_0 c0 
						JOIN ai_workflow_template_eval_task_relationships evtr ON evtr.workflow_template_id = c0.workflow_template_id
						LEFT JOIN eval_fns ef ON ef.eval_id = evtr.eval_id
					), cte_1 AS (
					SELECT
						cte_0.workflow_template_id,
						cte_0.analysis_task_id,
						JSON_BUILD_OBJECT(
							'analysisTaskID', cte_0.analysis_task_id,
							'analysisCycleCount', cte_0.analysis_cycle_count,
							'analysisTaskName', ait.task_name,
							'analysisTaskType', ait.task_type,
							'analysisPrompt', ait.prompt,
							'analysisModel', ait.model,
							'analysisMarginBuffer', ait.margin_buffer,
							'analysisTemperature', ait.temperature,
							'analysisTokenOverflowStrategy', ait.token_overflow_strategy,
							'analysisMaxTokensPerTask', ait.max_tokens_per_task,
							'analysisResponseFormat', ait.response_format,
							'retrievalID', COALESCE(art.retrieval_id, 0),
							'retrievalName', COALESCE(art.retrieval_name, ''),
							'retrievalGroup', COALESCE(art.retrieval_group, ''),
							'retrievalPlatform', COALESCE(art.retrieval_platform, ''),
							'retrievalInstructions', COALESCE(art.instructions, '{}'::jsonb),
							'analysisEvalFns', COALESCE(JSON_AGG(cte_wf_evals.eval_fns_data), '[]'::json)
						) AS analysis_tasks
					FROM cte_0 
					JOIN public.ai_task_library ait ON ait.task_id = cte_0.analysis_task_id
					LEFT JOIN ai_retrieval_library art ON art.retrieval_id = cte_0.retrieval_id
					LEFT JOIN cte_wf_evals ON cte_wf_evals.workflow_template_id = cte_0.workflow_template_id AND cte_wf_evals.task_id = cte_0.analysis_task_id
					GROUP BY cte_0.workflow_template_id, cte_0.analysis_task_id, ait.task_name, ait.task_type, ait.prompt, ait.model, ait.token_overflow_strategy, ait.temperature, ait.margin_buffer,
						ait.max_tokens_per_task, ait.response_format, art.retrieval_id, art.retrieval_name, art.retrieval_group, art.retrieval_platform, art.instructions, cte_0.analysis_cycle_count
				),
				cte_1a AS (
					SELECT 
						workflow_template_id, 
						jsonb_agg(analysis_tasks) as analysis_tasks_array
					FROM cte_1
					GROUP BY workflow_template_id
		), unique_evals AS (
			SELECT 
				workflow_template_id, 
				task_id, 
				JSONB_AGG(eval_fns_data) AS eval_fns_data_agg
			FROM cte_wf_evals
			GROUP BY workflow_template_id, task_id
		),
		cte_2 AS (
		SELECT
			wate.workflow_template_id,
			ait.task_id as agg_task_id,
			ait1.task_id as analysis_task_id,
			JSON_BUILD_OBJECT(
				'aggTaskId', ait.task_id,
				'aggAnalysisTaskId', ait1.task_id,
				'aggTaskName', ait.task_name,
				'aggTaskType', ait.task_type,
				'aggPrompt', ait.prompt,
				'aggModel', ait.model,
				'aggResponseFormat', ait.response_format,
				'aggTokenOverflowStrategy', ait.token_overflow_strategy,
				'aggMarginBuffer', ait.margin_buffer,
				'aggTemperature', ait.temperature,
				'aggMaxTokensPerTask', ait.max_tokens_per_task,
				'aggCycleCount', awtat.cycle_count,
				'evalFns', COALESCE(ue.eval_fns_data_agg, '[]'::jsonb),
				'analysisAggEvalFns', COALESCE(ue1.eval_fns_data_agg, '[]'::jsonb),
				'analysisAggCycleCount', awtat1.cycle_count
			) AS agg_tasks
		FROM ai_workflow_template wate
		JOIN public.ai_workflow_template_agg_tasks awtat ON awtat.workflow_template_id = wate.workflow_template_id
		JOIN public.ai_task_library ait ON ait.task_id = awtat.agg_task_id
		JOIN public.ai_task_library ait1 ON ait1.task_id = awtat.analysis_task_id
		JOIN public.ai_workflow_template_analysis_tasks awtat1 ON awtat1.workflow_template_id = wate.workflow_template_id AND awtat1.task_id = awtat.analysis_task_id
		LEFT JOIN unique_evals ue ON ue.workflow_template_id = wate.workflow_template_id AND ue.task_id = ait.task_id
		LEFT JOIN unique_evals ue1 ON ue1.workflow_template_id = wate.workflow_template_id AND ue1.task_id = ait1.task_id
   	 	WHERE wate.org_id = $1 ` + additionalCondition + `
    	GROUP BY wate.workflow_template_id, ait.task_id, ait1.task_id, ait.task_name, ait.task_type, ait.prompt, ait.model, ait.temperature, ait.margin_buffer,
             ait.token_overflow_strategy, ait.max_tokens_per_task, ait.response_format, awtat.cycle_count, ue.eval_fns_data_agg, ue1.eval_fns_data_agg, awtat1.cycle_count
				), cte_2a AS (
						SELECT 
							workflow_template_id, 
							jsonb_agg(agg_tasks) as agg_tasks_array
						FROM cte_2
						GROUP BY workflow_template_id
					), cte_2b AS (
						SELECT 
							workflow_template_id, 
							jsonb_agg(agg_tasks_array) as agg_tasks_array
						FROM cte_2a
						GROUP BY workflow_template_id
					)
					SELECT
						wate.workflow_template_id,
						wate.workflow_name,
						wate.workflow_group,
						wate.fundamental_period,
						wate.fundamental_period_time_unit,
						COALESCE(cte_1a.analysis_tasks_array, '[]') as analysis_tasks_array,
						COALESCE(cte_2b.agg_tasks_array, '[]') as agg_tasks_array
					FROM ai_workflow_template wate
					LEFT JOIN cte_1a ON wate.workflow_template_id = cte_1a.workflow_template_id
					LEFT JOIN cte_2b ON cte_2b.workflow_template_id = wate.workflow_template_id
					WHERE wate.org_id = $1 ` + additionalCondition + `
					ORDER BY wate.workflow_template_id DESC`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var taskJSON string
		var aggTasksJSON *string
		wt := WorkflowTemplateValue{
			AnalysisTasks:      make(map[int]AnalysisTaskDB),
			AnalysisRetrievals: make(map[int]map[int]RetrievalDB),
			AnalysisEvalFns:    make(map[int][]EvalFnDB), // Initialize map for analysis eval functions
			AggTasks:           make(map[int]AggTaskDb),
			AggAnalysisTasks:   make(map[int]map[int]AggTaskDb),
			AggAnalysisEvalFns: make(map[int]map[int]EvalFnDB),
			AggEvalFns:         make(map[int][]EvalFnDB), // Initialize map for agg eval functions
		}
		err = rows.Scan(
			&wt.WorkflowTemplateID,
			&wt.WorkflowName,
			&wt.WorkflowGroup,
			&wt.FundamentalPeriod,
			&wt.FundamentalPeriodTimeUnit,
			&taskJSON,
			&aggTasksJSON,
		)
		if err != nil {
			log.Err(err).Msg("Error scanning row in SelectWorkflowTemplate")
			return nil, err
		}
		wt.WorkflowTemplateStrID = fmt.Sprintf("%d", wt.WorkflowTemplateID)
		err = json.Unmarshal([]byte(taskJSON), &wt.AnalysisTasksSlice)
		if err != nil {
			log.Err(err).Msg("Error unmarshalling analysis tasks JSON")
			return nil, err
		}

		if aggTasksJSON != nil {
			var aggTasksPreFlatten [][]AggTaskDb
			err = json.Unmarshal([]byte(*aggTasksJSON), &aggTasksPreFlatten)
			if err != nil {
				log.Err(err).Msg("Error unmarshalling agg tasks JSON")
				return nil, err
			}
			for _, aggTask := range aggTasksPreFlatten {
				wt.AggAnalysisTasksSlice = append(wt.AggAnalysisTasksSlice, aggTask...)
			}
		}
		if wt.AggAnalysisTasks == nil {
			wt.AggAnalysisTasks = make(map[int]map[int]AggTaskDb)
		}
		if wt.AggAnalysisEvalFns == nil {
			wt.AggAnalysisEvalFns = make(map[int]map[int]EvalFnDB)
		}
		for i, v := range wt.AggAnalysisTasksSlice {
			if _, ok := wt.AggAnalysisTasks[v.AggTaskId]; !ok {
				wt.AggAnalysisTasks[v.AggTaskId] = make(map[int]AggTaskDb)
			}
			wt.AggEvalFns[v.AggTaskId] = v.EvalFns
			wt.AggAnalysisTasks[v.AggTaskId][v.AggAnalysisTaskId] = v
			var tmp []EvalFnDB
			seen := make(map[int]bool)
			for _, ef := range wt.AggAnalysisTasksSlice[i].AnalysisAggEvalFns {
				if _, ok := seen[ef.EvalID]; ok {
					continue
				}
				tmp = append(tmp, ef)
				seen[ef.EvalID] = true
				if _, ok := wt.AggAnalysisEvalFns[v.AggTaskId]; !ok {
					wt.AggAnalysisEvalFns[v.AggTaskId] = make(map[int]EvalFnDB)
				}
				wt.AggAnalysisEvalFns[v.AggTaskId][ef.EvalTaskID] = ef
			}
			wt.AggAnalysisTasks[v.AggTaskId][v.AggAnalysisTaskId] = wt.AggAnalysisTasksSlice[i]
		}

		for i, v := range wt.AnalysisTasksSlice {
			if _, ok := wt.AnalysisRetrievals[v.AnalysisTaskID]; !ok {
				wt.AnalysisRetrievals[v.AnalysisTaskID] = make(map[int]RetrievalDB)
			}
			wt.AnalysisTasks[v.AnalysisTaskID] = v
			if v.RetrievalID != nil && *v.RetrievalID > 0 {
				wt.AnalysisRetrievals[v.AnalysisTaskID][*v.RetrievalID] = v.RetrievalDB
			}

			var tmp []EvalFnDB
			seen := make(map[int]bool)
			for _, ef := range wt.AnalysisTasksSlice[i].AnalysisEvalFns {
				if _, ok := seen[ef.EvalID]; ok {
					continue
				}
				tmp = append(tmp, ef)
				seen[ef.EvalID] = true
			}
			wt.AnalysisTasksSlice[i].AnalysisEvalFns = tmp
			wt.AnalysisTasks[v.AnalysisTaskID] = wt.AnalysisTasksSlice[i]
			wt.AnalysisEvalFns[v.AnalysisTaskID] = tmp
		}
		if _, ok := results.WorkflowTemplatesMap[wt.WorkflowTemplateID]; !ok {
			results.WorkflowTemplatesMap[wt.WorkflowTemplateID] = wt
		}
	}
	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Err(err).Msg("error iterating over rows")
		return nil, err
	}

	for i, v := range results.WorkflowTemplatesMap {
		var taskVals []Task
		for _, at := range results.WorkflowTemplatesMap[i].AnalysisTasksSlice {
			ta := Task{
				TaskStrID:         fmt.Sprintf("%d", at.AnalysisTaskID),
				TaskID:            at.AnalysisTaskID,
				TaskName:          at.AnalysisTaskName,
				TaskType:          at.AnalysisTaskType,
				ResponseFormat:    at.AnalysisResponseFormat,
				Model:             at.AnalysisModel,
				Prompt:            at.AnalysisPrompt,
				Temperature:       at.AnalysisTemperature,
				MarginBuffer:      at.AnalysisMarginBuffer,
				CycleCount:        at.AnalysisCycleCount,
				RetrievalName:     at.RetrievalName,
				RetrievalPlatform: at.RetrievalPlatform,
				EvalFnDBs:         at.AnalysisEvalFns,
			}
			tmp := results.WorkflowTemplatesMap[i]
			if _, ok := results.WorkflowTemplatesMap[i].AnalysisTasks[at.AnalysisTaskID]; !ok {
				tmp.AnalysisTasks[at.AnalysisTaskID] = at
			}
			v.AnalysisTasks[at.AnalysisTaskID] = at
			tmp.AnalysisTasks = v.AnalysisTasks
			results.WorkflowTemplatesMap[i] = tmp
			taskVals = append(taskVals, ta)
		}
		for _, aggTask := range results.WorkflowTemplatesMap[i].AggAnalysisTasksSlice {
			rn := ""
			agat := v.AnalysisTasks[aggTask.AggAnalysisTaskId]
			if agat.AnalysisTaskName != "" {
				rn = agat.AnalysisTaskName
			}
			ta := Task{
				TaskStrID:         fmt.Sprintf("%d", aggTask.AggTaskId),
				TaskID:            aggTask.AggTaskId,
				TaskName:          aggTask.AggTaskName,
				TaskType:          aggTask.AggTaskType,
				Temperature:       aggTask.AggTemperature,
				MarginBuffer:      aggTask.AggMarginBuffer,
				ResponseFormat:    aggTask.AggResponseFormat,
				Model:             aggTask.AggModel,
				Prompt:            aggTask.AggPrompt,
				CycleCount:        aggTask.AggCycleCount,
				RetrievalName:     rn,
				RetrievalPlatform: "aggregate-analysis",
				EvalFnDBs:         aggTask.EvalFns,
			}
			taskVals = append(taskVals, ta)
		}
		v.Tasks = taskVals
		results.WorkflowTemplateSlice = append(results.WorkflowTemplateSlice, v)
	}
	return results, nil
}

func MapDependencies(res []WorkflowTemplateData) WorkflowTaskRelationships {
	analysisRetrievals := make(map[int]map[int]bool)
	aggregateAnalysis := make(map[int]map[int]bool)

	for _, v := range res {
		if _, ok := analysisRetrievals[v.AnalysisTaskID]; !ok {
			analysisRetrievals[v.AnalysisTaskID] = make(map[int]bool)
		}
		if v.RetrievalID != nil && *v.RetrievalID != 0 {
			if _, ok := analysisRetrievals[v.AnalysisTaskID][*v.RetrievalID]; !ok {
				analysisRetrievals[v.AnalysisTaskID][*v.RetrievalID] = true
			} else {
				//fmt.Println("Duplicate retrieval id", v.RetrievalID)
			}
		}

		if v.AggTaskID != nil {
			if _, ok := aggregateAnalysis[*v.AggTaskID]; !ok {
				aggregateAnalysis[*v.AggTaskID] = make(map[int]bool)
			}
			if _, ok := aggregateAnalysis[*v.AggTaskID][v.AnalysisTaskID]; !ok {
				aggregateAnalysis[*v.AggTaskID][v.AnalysisTaskID] = true
			} else {
				//fmt.Println("Duplicate agg id", *v.AggTaskID)
			}
		}
	}
	return WorkflowTaskRelationships{
		AnalysisRetrievals: analysisRetrievals,
		AggregateAnalysis:  aggregateAnalysis,
	}
}
func MapDependenciesGrouped(res WorkflowTemplateValue) WorkflowTaskRelationships {
	analysisRetrievals := make(map[int]map[int]bool)
	for _, analysisTask := range res.AnalysisTasks {
		analysisTaskID := analysisTask.AnalysisTaskID
		if _, ok := analysisRetrievals[analysisTaskID]; !ok {
			analysisRetrievals[analysisTaskID] = make(map[int]bool)
		}
	}
	for analysisTaskID, retrievalMap := range res.AnalysisRetrievals {
		for _, retrieval := range retrievalMap {
			if _, ok := analysisRetrievals[analysisTaskID]; !ok {
				analysisRetrievals[analysisTaskID] = make(map[int]bool)
			}
			if retrieval.RetrievalID != nil && *retrieval.RetrievalID != 0 {
				analysisRetrievals[analysisTaskID][*retrieval.RetrievalID] = true
			}
		}
	}
	aggregateAnalysis := make(map[int]map[int]bool)

	for aggTaskID, analysisMap := range res.AggAnalysisTasks {
		for _, analysis := range analysisMap {
			if _, ok := aggregateAnalysis[aggTaskID]; !ok {
				aggregateAnalysis[aggTaskID] = make(map[int]bool)
			}
			aggregateAnalysis[aggTaskID][analysis.AggAnalysisTaskId] = true
		}
	}
	return WorkflowTaskRelationships{
		AnalysisRetrievals: analysisRetrievals,
		AggregateAnalysis:  aggregateAnalysis,
	}
}
