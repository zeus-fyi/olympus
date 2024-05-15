package artemis_orchestrations

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func SelectAiSystemOrchestrationsUI(ctx context.Context, ou org_users.OrgUser, rid int) ([]OrchestrationsAnalysis, error) {
	var ojs []OrchestrationsAnalysis
	q := sql_query_templates.QueryParams{}
	args := []interface{}{ou.OrgID}

	var limit string
	queryByRunID := ""
	if rid > 0 {
		queryByRunID = " AND ar.workflow_run_id = $2"
		args = append(args, rid)
	} else {
		limit = " LIMIT 100"
	}
	/*
	   ALTER TABLE public.ai_workflow_runs
	       ADD COLUMN total_api_requests int8 NOT NULL DEFAULT 0;

	   ALTER TABLE public.ai_workflow_runs
	       ADD COLUMN total_csv_cells int8 NOT NULL DEFAULT 0;
	*/
	// uses main for unique id, so type == real name for related workflow
	q.RawQuery = `WITH cte_a AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name AS group_name,
							o.type AS orchestration_type,
							o.active,
							o.org_id,
							ar.total_csv_cells,
							ar.total_api_requests
						FROM 
							public.ai_workflow_runs AS ar 
					  	JOIN
							public.orchestrations AS o ON o.orchestration_id = ar.orchestration_id
						WHERE 
							o.org_id = $1 AND ar.is_archived = false` + queryByRunID + ` 
						ORDER BY
							o.orchestration_id DESC
						` + limit + `
					), cte_0 AS (
						SELECT
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_res.workflow_result_id,
							ai_res.response_id,
							ai_res.source_task_id,
							ai_res.iteration_count,
							ai_res.chunk_offset,
							ai_res.running_cycle_number,
							ai_res.skip_analysis,
							ai_res.search_window_unix_start,
							ai_res.search_window_unix_end,
							ai_res.metadata
						FROM 
							cte_a o 
						JOIN
							public.ai_workflow_analysis_results ai_res ON ai_res.orchestration_id = o.orchestration_id
						GROUP BY 								
							o.orchestration_id,
							o.orchestration_name,
							o.group_name,
							o.orchestration_type,
							o.active,
							ai_res.workflow_result_id,
							ai_res.response_id,
							ai_res.source_task_id,
							ai_res.iteration_count,
							ai_res.chunk_offset,
							ai_res.running_cycle_number,
							ai_res.skip_analysis,
							ai_res.search_window_unix_start,
							ai_res.search_window_unix_end,
							ai_res.metadata
						), cte_00 AS (
							SELECT 
								ai_res.orchestration_id,
								MAX(ai_res.running_cycle_number) AS max_run_cycle,
								SUM(comp_resp.total_tokens) AS total_workflow_token_usage
							FROM cte_0 ai_res
							JOIN 
								public.completion_responses AS comp_resp ON comp_resp.response_id = ai_res.response_id
							GROUP BY
								ai_res.orchestration_id
					)
							SELECT 
								ca.orchestration_id,
								ca.orchestration_id::text,
								ca.orchestration_name,
								ca.group_name,
								ca.orchestration_type,
								ca.active,
								ca.total_csv_cells,
								ca.total_api_requests,
							  	COALESCE(c00.max_run_cycle, 0),
							  	COALESCE(c00.total_workflow_token_usage, 0) AS total_workflow_token_usage
							 FROM cte_a ca 
							  LEFT JOIN cte_00 c00 ON c00.orchestration_id = ca.orchestration_id
							 GROUP BY ca.orchestration_id,
								ca.orchestration_id::text,
								ca.orchestration_name,
								ca.group_name,
								ca.orchestration_type,
								ca.total_csv_cells,
								ca.total_api_requests,
								ca.active, c00.max_run_cycle, c00.total_workflow_token_usage
							  ORDER BY ca.orchestration_id DESC;`
	var ojsRunsActions []OrchestrationsAnalysis
	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, args...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationsAnalysis{}
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationStrID, &oj.OrchestrationName, &oj.GroupName,
			&oj.Type, &oj.Active, &oj.TotalCsvCells, &oj.TotalApiRequests, &oj.RunCycles, &oj.TotalWorkflowTokenUsage)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return nil, rowErr
		}
		if oj.Active {
			mad := len(oj.AggregatedData)
			if mad > oj.TotalCsvCells {
				// alert?
				mad = oj.TotalCsvCells
			}
			rad := len(oj.AggregatedRetrievalResults)
			if rad > oj.TotalApiRequests {
				// todo alert?
				rad = oj.TotalApiRequests
			}
			oj.Progress = (float64(mad+rad) / float64(oj.TotalApiRequests+oj.TotalApiRequests)) * float64(100)
		} else {
			oj.Progress = 100
		}
		ojsRunsActions = append(ojsRunsActions, oj)
	}
	for _, oj := range ojs {
		if !oj.Active && oj.RunCycles == 0 && oj.TotalWorkflowTokenUsage == 0 {
			continue
		}
	}
	sortRunsByID(ojsRunsActions)
	return ojsRunsActions, err
}
