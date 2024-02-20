package artemis_orchestrations

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type WorkflowStageReference struct {
	InputID        int             `json:"inputID"`
	InputStrID     string          `json:"inputStrID"`
	WorkflowRunID  int             `json:"workflowRunID"`
	ChildWfID      string          `json:"childWfID"`
	RunCycle       int             `json:"runCycle"`
	IterationCount int             `json:"iterationCount"`
	ChunkOffset    int             `json:"chunk"`
	InputData      json.RawMessage `json:"inputData"`
	Logs           []string        `json:"logs"`
	LogsStr        string          `json:"-"`
}

func InsertWorkflowStageReference(ctx context.Context, wfStageIO *WorkflowStageReference) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
			INSERT INTO public.ai_workflow_stage_references (input_id, workflow_run_id, child_wf_id, run_cycle, input_data, logs)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (input_id) DO UPDATE
			SET input_data = EXCLUDED.input_data,
				logs = EXCLUDED.logs
			RETURNING input_id, input_id::text;
`
	ts := chronos.Chronos{}
	if wfStageIO.InputID == 0 {
		wfStageIO.InputID = ts.UnixTimeStampNow()
	}
	var jsonb pgtype.JSONB
	jsonb.Bytes = sanitizeBytesUTF8(wfStageIO.InputData)
	jsonb.Status = IsNull(wfStageIO.InputData)

	wfStageIO.LogsStr = strings.Join(wfStageIO.Logs, ",")
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, wfStageIO.InputID, wfStageIO.WorkflowRunID, wfStageIO.ChildWfID, wfStageIO.RunCycle, jsonb, wfStageIO.LogsStr).Scan(&wfStageIO.InputID, &wfStageIO.InputStrID)
	if err != nil {
		log.Err(err).Msg("failed to execute query for InsertWorkflowStageReference")
		return err
	}
	return nil
}

func SelectWorkflowStageReference(ctx context.Context, inputID int) (WorkflowStageReference, error) {
	var sr WorkflowStageReference
	// Update the SQL query to select by input_id
	q := `SELECT input_id, input_id::text, workflow_run_id, child_wf_id, run_cycle, input_data, logs FROM public.ai_workflow_stage_references WHERE input_id = $1`
	// Executing the query
	err := apps.Pg.QueryRowWArgs(ctx, q, inputID).Scan(&sr.InputID, &sr.InputStrID, &sr.WorkflowRunID, &sr.ChildWfID, &sr.RunCycle, &sr.InputData, &sr.LogsStr)
	if err != nil {
		log.Err(err).Msg("failed to execute query for SelectWorkflowStageReference")
		return sr, err
	}
	sr.Logs = strings.Split(sr.LogsStr, ",")
	return sr, nil
}
