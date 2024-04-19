package artemis_orchestrations

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type AIWorkflowRetrievalResult struct {
	WorkflowResultStrID   string          `json:"workflowResultStrID,omitempty"`
	WorkflowResultID      int             `json:"workflowResultID"`
	OrchestrationID       int             `json:"orchestrationID"`
	RetrievalName         string          `json:"retrievalName,omitempty"`
	RetrievalID           int             `json:"retrievalID"`
	IterationCount        int             `json:"iterationCount"`
	ChunkOffset           int             `json:"chunkOffset"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd"`
	SkipRetrieval         bool            `json:"skipRetrieval"`
	Attempts              int             `json:"attempts"`
	Status                string          `json:"status"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
}

func SelectRetrievalResultsIds(ctx context.Context, w Window, ojIds, sourceRetIds []int) ([]AIWorkflowRetrievalResult, error) {
	query := `SELECT ar.workflow_result_id, ar.workflow_result_id::text, ar.orchestration_id, ar.retrieval_id, ar.chunk_offset, ar.iteration_count,
       ar.running_cycle_number, ar.search_window_unix_start, ar.search_window_unix_end, ar.attempts
       FROM ai_workflow_io_results ar
       WHERE ar.skip_retrieval = false AND ar.search_window_unix_start >= $1 AND ar.search_window_unix_end <= $2 OR (status = 'error' AND attempts > 5)
         AND ar.retrieval_id = ANY($3) AND ar.orchestration_id = ANY($4);`
	rows, err := apps.Pg.Query(ctx, query, w.UnixStartTime, w.UnixEndTime, pq.Array(sourceRetIds), pq.Array(ojIds))
	if err != nil {
		log.Printf("Error executing SelectRetrievalResultsIds query: %v", err)
		return nil, err
	}
	defer rows.Close()
	var results []AIWorkflowRetrievalResult
	for rows.Next() {
		var result AIWorkflowRetrievalResult
		err = rows.Scan(&result.WorkflowResultID, &result.WorkflowResultStrID, &result.OrchestrationID, &result.RetrievalID,
			&result.ChunkOffset, &result.IterationCount, &result.RunningCycleNumber,
			&result.SearchWindowUnixStart, &result.SearchWindowUnixEnd, &result.Attempts)
		if err != nil {
			log.Printf("Error scanning SelectRetrievalResultsIds result: %v", err)
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func InsertWorkflowRetrievalResult(ctx context.Context, wr *AIWorkflowRetrievalResult) error {
	q := `INSERT INTO ai_workflow_io_results(orchestration_id, retrieval_id, status, running_cycle_number, iteration_count, 
                                         chunk_offset, search_window_unix_start, search_window_unix_end, skip_retrieval, metadata)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                  ON CONFLICT (orchestration_id, retrieval_id, running_cycle_number, iteration_count, chunk_offset)
                  DO UPDATE SET 
                      running_cycle_number = EXCLUDED.running_cycle_number,
                      search_window_unix_start = EXCLUDED.search_window_unix_start,
                      search_window_unix_end = EXCLUDED.search_window_unix_end,
                      skip_retrieval = EXCLUDED.skip_retrieval,
                      metadata = EXCLUDED.metadata
                  RETURNING workflow_result_id;`
	md := &pgtype.JSONB{Bytes: sanitizeBytesUTF8(wr.Metadata), Status: IsNull(wr.Metadata)}
	err := apps.Pg.QueryRowWArgs(ctx, q, wr.OrchestrationID, wr.RetrievalID, wr.Status, wr.RunningCycleNumber,
		wr.IterationCount, wr.ChunkOffset, wr.SearchWindowUnixStart, wr.SearchWindowUnixEnd,
		wr.SkipRetrieval, md).Scan(&wr.WorkflowResultID)
	if err != nil {
		log.Err(err).Msgf("Error executing InsertWorkflowRetrievalResultError query: %v", err)
		return err
	}
	return nil
}

func InsertWorkflowRetrievalResultError(ctx context.Context, wr *AIWorkflowRetrievalResult) error {
	q := `	INSERT INTO ai_workflow_io_results(
				orchestration_id, retrieval_id, status, running_cycle_number, iteration_count, 
				chunk_offset, search_window_unix_start, search_window_unix_end, skip_retrieval, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT (orchestration_id, retrieval_id, running_cycle_number, iteration_count, chunk_offset)
			DO UPDATE SET 
				running_cycle_number = EXCLUDED.running_cycle_number,
				search_window_unix_start = EXCLUDED.search_window_unix_start,
				search_window_unix_end = EXCLUDED.search_window_unix_end,
				skip_retrieval = EXCLUDED.skip_retrieval,
				metadata = EXCLUDED.metadata,
				status = EXCLUDED.status
			RETURNING workflow_result_id, attempts
		`
	md := &pgtype.JSONB{Bytes: sanitizeBytesUTF8(wr.Metadata), Status: IsNull(wr.Metadata)}
	err := apps.Pg.QueryRowWArgs(ctx, q, wr.OrchestrationID, wr.RetrievalID, wr.Status, wr.RunningCycleNumber,
		wr.IterationCount, wr.ChunkOffset, wr.SearchWindowUnixStart, wr.SearchWindowUnixEnd,
		wr.SkipRetrieval, md).Scan(&wr.WorkflowResultID, &wr.Attempts)
	if err != nil {
		log.Err(err).Msgf("Error executing InsertWorkflowRetrievalResultError query: %v", err)
		return err
	}
	q = `UPDATE ai_workflow_io_results
		 SET attempts = $2
		 WHERE workflow_result_id = $1
		 RETURNING attempts
		`
	err = apps.Pg.QueryRowWArgs(ctx, q, wr.WorkflowResultID, wr.Attempts+1).Scan(&wr.Attempts)
	if err != nil {
		log.Err(err).Msgf("Error executing InsertWorkflowRetrievalResultError query: %v", err)
		return err
	}
	return nil
}
