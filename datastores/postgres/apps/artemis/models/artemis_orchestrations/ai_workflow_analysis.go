package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type AIWorkflowAnalysisResult struct {
	WorkflowResultID      int             `json:"workflowResultID"`
	OrchestrationID       int             `json:"orchestrationID"`
	ResponseID            int             `json:"responseID"`
	SourceTaskID          int             `json:"sourceTaskID"`
	IterationCount        int             `json:"iterationCount"`
	ChunkOffset           int             `json:"chunkOffset"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd"`
	SkipAnalysis          bool            `json:"skipAnalysis"`
	TaskOffset            int             `json:"taskOffset"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	CompletionChoices     json.RawMessage `json:"completionChoices,omitempty"`
}

func InsertAiWorkflowAnalysisResult(ctx context.Context, wr *AIWorkflowAnalysisResult) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO ai_workflow_analysis_results(orchestration_id, response_id, source_task_id, chunk_offset, iteration_count, 
                                         running_cycle_number, search_window_unix_start, search_window_unix_end, skip_analysis, metadata, task_offset)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
                  ON CONFLICT (orchestration_id, response_id, source_task_id, running_cycle_number, iteration_count, chunk_offset)
                  DO UPDATE SET 
                      task_offset = EXCLUDED.task_offset,
                      orchestration_id = EXCLUDED.orchestration_id,
                      response_id = EXCLUDED.response_id,
                      source_task_id = EXCLUDED.source_task_id,
                      running_cycle_number = EXCLUDED.running_cycle_number,
                      search_window_unix_start = EXCLUDED.search_window_unix_start,
                      search_window_unix_end = EXCLUDED.search_window_unix_end,
                      skip_analysis = EXCLUDED.skip_analysis,
                      metadata = EXCLUDED.metadata
                  RETURNING workflow_result_id;`
	md := &pgtype.JSONB{Bytes: sanitizeBytesUTF8(wr.Metadata), Status: IsNull(wr.Metadata)}
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, wr.OrchestrationID, wr.ResponseID, wr.SourceTaskID,
		wr.ChunkOffset, wr.IterationCount, wr.RunningCycleNumber,
		wr.SearchWindowUnixStart, wr.SearchWindowUnixEnd, wr.SkipAnalysis,
		md, wr.TaskOffset).Scan(&wr.WorkflowResultID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("AIWorkflowAnalysisResults")); returnErr != nil {
		log.Err(returnErr).Interface("wr", wr).Msg(q.LogHeader("AIWorkflowAnalysisResults"))
		return err
	}
	return nil
}

func SelectAiWorkflowAnalysisResults(ctx context.Context, w Window, ojIds, sourceTaskIds []int) ([]AIWorkflowAnalysisResult, error) {
	q := sql_query_templates.QueryParams{}
	// Then, select rows using the search window and source task IDs
	q.RawQuery = `SELECT ar.workflow_result_id, ar.orchestration_id, ar.response_id, ar.source_task_id,ar.chunk_offset, ar.iteration_count,
       					 ar.running_cycle_number, ar.search_window_unix_start, ar.search_window_unix_end, ar.metadata, cr.completion_choices, ar.task_offset
                  FROM ai_workflow_analysis_results ar
                  JOIN completion_responses cr ON cr.response_id = ar.response_id	
                  WHERE ar.skip_analysis = false AND ar.search_window_unix_start >= $1 AND ar.search_window_unix_end <= $2
                    AND ar.source_task_id = ANY($3) AND ar.orchestration_id = ANY($4);`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, w.UnixStartTime, w.UnixEndTime, pq.Array(sourceTaskIds), pq.Array(ojIds))
	if err != nil {
		log.Err(err).Msg(q.LogHeader("AIWorkflowAnalysisResults"))
		return nil, err
	}
	defer rows.Close()

	var results []AIWorkflowAnalysisResult
	for rows.Next() {
		var result AIWorkflowAnalysisResult
		err = rows.Scan(&result.WorkflowResultID, &result.OrchestrationID, &result.ResponseID, &result.SourceTaskID,
			&result.ChunkOffset, &result.IterationCount, &result.RunningCycleNumber,
			&result.SearchWindowUnixStart, &result.SearchWindowUnixEnd, &result.Metadata, &result.CompletionChoices, &result.TaskOffset)
		if err != nil {
			log.Err(err).Msg(q.LogHeader("AIWorkflowAnalysisResults"))
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func SelectAiWorkflowAnalysisResultsIds(ctx context.Context, w Window, ojIds, sourceTaskIds []int, to int) ([]AIWorkflowAnalysisResult, error) {
	q := sql_query_templates.QueryParams{}
	// Then, select rows using the search window and source task IDs
	q.RawQuery = `SELECT ar.workflow_result_id, ar.orchestration_id, ar.response_id, ar.source_task_id,ar.chunk_offset, ar.iteration_count,
       					 ar.running_cycle_number, ar.search_window_unix_start, ar.search_window_unix_end
                  FROM ai_workflow_analysis_results ar
                  JOIN completion_responses cr ON cr.response_id = ar.response_id	
                  WHERE ar.skip_analysis = false AND ar.search_window_unix_start >= $1 AND ar.search_window_unix_end <= $2
                    AND ar.source_task_id = ANY($3) AND ar.orchestration_id = ANY($4) AND ar.chunk_offset = $5 AND ar.task_offset = $5;`

	rows, err := apps.Pg.Query(ctx, q.RawQuery, w.UnixStartTime, w.UnixEndTime, pq.Array(sourceTaskIds), pq.Array(ojIds), to)
	if err != nil {
		log.Err(err).Msg(q.LogHeader("AIWorkflowAnalysisResults"))
		return nil, err
	}
	defer rows.Close()

	var results []AIWorkflowAnalysisResult
	for rows.Next() {
		var result AIWorkflowAnalysisResult
		err = rows.Scan(&result.WorkflowResultID, &result.OrchestrationID, &result.ResponseID, &result.SourceTaskID,
			&result.ChunkOffset, &result.IterationCount, &result.RunningCycleNumber,
			&result.SearchWindowUnixStart, &result.SearchWindowUnixEnd)
		if err != nil {
			log.Err(err).Msg(q.LogHeader("AIWorkflowAnalysisResults"))
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// todo more efficient way to do this

func GenerateContentTextFromOpenAIResp(wrs []AIWorkflowAnalysisResult) (string, error) {
	var temp string

	for _, wr := range wrs {
		if len(wr.CompletionChoices) > 0 {
			// todo, this isn't working
			var cc []openai.ChatCompletionChoice
			err := json.Unmarshal(wr.CompletionChoices, &cc)
			if err != nil {
				log.Err(err).Msg("ChatCompletionChoice Unmarshal Error")
				return "", err
			}
			for _, c := range cc {
				temp += c.Message.Content + "\n"
			}
		}
		if len(wr.Metadata) > 0 && string(wr.Metadata) != "null" {
			temp += string(wr.Metadata) + "\n"
		}
	}

	// todo refactor after fixing the above
	if (len(temp) <= 0 || temp == "\n") && len(wrs) > 0 {
		tc := ToolCompletions{}
		err := json.Unmarshal(wrs[0].CompletionChoices, &tc)
		if err != nil {
			log.Err(err).Msg("ToolCompletions Unmarshal Error")
			var tmp = []ToolCompletions{}
			err = json.Unmarshal(wrs[0].CompletionChoices, &tmp)
			if err != nil {
				log.Err(err).Msg("ToolCompletions UnmarshalSlice Error")
			}
			for _, tcIn := range tmp {
				for _, c := range tcIn.Message.ToolCalls {
					temp += c.Function.Arguments + "\n"
				}
			}
			err = nil
		}
		for _, c := range tc.Message.ToolCalls {
			temp += c.Function.Arguments + "\n"
		}

	}
	return temp, nil
}

type ToolCompletions struct {
	Index   int `json:"index"`
	Message struct {
		Role      string `json:"role"`
		Content   string `json:"content"`
		ToolCalls []struct {
			Id       string `json:"id"`
			Type     string `json:"type"`
			Function struct {
				Name      string `json:"name"`
				Arguments string `json:"arguments"`
			} `json:"function"`
		} `json:"tool_calls"`
	} `json:"message"`
	FinishReason string `json:"finish_reason"`
}
