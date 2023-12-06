package artemis_orchestrations

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

const (
	Orchestrations = "Orchestrations"
)

type OrchestrationJob struct {
	artemis_autogen_bases.Orchestrations `json:"orchestrations,omitempty"`
	Scheduled                            artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs `json:"scheduled,omitempty"`
	zeus_common_types.CloudCtxNs         `json:"cloud_ctx_ns,omitempty"`
}

const (
	internalOrgID = 7138983863666903883
)

func NewInternalActiveTemporalOrchestrationJobTemplate(orchName, groupName, orchType string) OrchestrationJob {
	return NewActiveTemporalOrchestrationJobTemplate(internalOrgID, orchName, groupName, orchType)
}

func NewActiveTemporalOrchestrationJobTemplate(orgID int, orchName, groupName, orchType string) OrchestrationJob {
	return OrchestrationJob{
		Orchestrations: artemis_autogen_bases.Orchestrations{
			OrgID:             orgID,
			Active:            true,
			GroupName:         groupName,
			Type:              orchType,
			OrchestrationName: orchName,
		},
		Scheduled:  artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs{},
		CloudCtxNs: zeus_common_types.CloudCtxNs{},
	}
}

func NewActiveTemporalOrchestrationJobTemplateWithInstructions(orgID int, orchName, groupName, orchType, instructions string) OrchestrationJob {
	oj := NewActiveTemporalOrchestrationJobTemplate(orgID, orchName, groupName, orchType)
	oj.Instructions = instructions
	return oj
}

func InsertOrchestration(ctx context.Context, oj OrchestrationJob, b []byte) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO orchestrations(org_id, orchestration_name, group_name, type, instructions)
				  VALUES ($1, $2, $3, $4, $5)
				  ON CONFLICT (org_id, orchestration_name) 
				  DO UPDATE SET instructions = EXCLUDED.instructions
				  RETURNING orchestration_id;`

	var id int
	log.Debug().Interface("InsertOrchestrations", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, oj.OrgID, oj.OrchestrationName, oj.GroupName, oj.Type, &pgtype.JSONB{Bytes: sanitizeBytesUTF8(b), Status: IsNull(b)}).Scan(&id)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return 0, err
	}
	return id, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func IsNull(b []byte) pgtype.Status {
	if b == nil {
		return pgtype.Null
	}
	return pgtype.Present
}

func sanitizeBytesUTF8(b []byte) []byte {
	bs := bytes.ReplaceAll(b, []byte{0}, []byte{})
	return bs
}

func SelectActiveOrchestrationsWithInstructionsUsingTimeWindow(ctx context.Context, orgID int, orchestType, groupName string, updatedAtWindowThreshold time.Duration) ([]OrchestrationJob, error) {
	var ojs []OrchestrationJob
	q := sql_query_templates.QueryParams{}
	thresholdTime := time.Now().UTC().Add(-updatedAtWindowThreshold) // Calculate the threshold time
	q.RawQuery = `SELECT orchestration_id, orchestration_name, instructions, type, group_name, org_id
				  FROM orchestrations
				  WHERE org_id = $1 AND active = true AND type = $2 AND group_name = $3 AND updated_at < $4
				  `
	log.Debug().Interface("SelectActiveOrchestrationsWithInstructionsUsingTimeWindow", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, orchestType, groupName, thresholdTime)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return ojs, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationJob{}
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.Instructions, &oj.Type, &oj.GroupName, &oj.OrgID)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return ojs, rowErr
		}
		ojs = append(ojs, oj)
	}
	return ojs, err
}

func SelectOrchestrationByName(ctx context.Context, orgID int, name string) (OrchestrationJob, error) {
	var oj OrchestrationJob
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT orchestration_id, orchestration_name, instructions, type, group_name
				  FROM orchestrations
				  WHERE org_id = $1 AND orchestration_name = $2 
				  `
	log.Debug().Interface("SelectOrchestrationByName", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, name).Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.Instructions, &oj.Type, &oj.GroupName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return oj, err
	}
	return oj, err
}

func SelectSystemOrchestrationsWithInstructionsByGroup(ctx context.Context, orgID int, groupName string) ([]OrchestrationJob, error) {
	var ojs []OrchestrationJob
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT orchestration_id, orchestration_name, instructions, type, group_name, org_id
				  FROM orchestrations
				  WHERE org_id = $1 AND active = true AND group_name = $2
				  `
	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, groupName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return ojs, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationJob{}
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.Instructions, &oj.Type, &oj.GroupName, &oj.OrgID)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return ojs, rowErr
		}
		ojs = append(ojs, oj)
	}
	return ojs, err
}

func (o *OrchestrationJob) UpsertOrchestrationWithInstructions(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO orchestrations(org_id, orchestration_name, instructions, type, group_name, active)
				  VALUES ($1, $2, $3, $4, $5, $6)
				  ON CONFLICT (org_id, orchestration_name)
				  DO UPDATE SET 
					  instructions = EXCLUDED.instructions,
					  type = EXCLUDED.type,
					  group_name = EXCLUDED.group_name,
					  active = EXCLUDED.active
				  RETURNING orchestration_id;
				  `
	log.Debug().Interface("InsertOrchestrationsWithInstructions", q.LogHeader(Orchestrations))
	if len(o.Instructions) == 0 {
		o.Instructions = "{}"
	}

	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, o.OrchestrationName, o.Instructions, o.Type, o.GroupName, o.Active).Scan(&o.OrchestrationID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func (o *OrchestrationJob) UpdateOrchestrationActiveStatus(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE orchestrations
				  SET active = $3
				  WHERE org_id = $1 AND orchestration_name = $2
				  RETURNING orchestration_id;
				  `
	log.Debug().Interface("UpdateOrchestrationActiveStatus", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, o.OrchestrationName, o.Active).Scan(&o.OrchestrationID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func (o *OrchestrationJob) InsertOrchestrationsScheduledToCloudCtxNsUsingName(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				WITH cte_get_cloud_ctx AS (
					SELECT cloud_ctx_ns_id, org_id, cloud_provider, region, context, namespace
					FROM topologies_org_cloud_ctx_ns
					WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4
					LIMIT 1
				), cte_get_orchestration_id AS (
					SELECT orchestration_id
					FROM orchestrations	
					WHERE orchestration_name = $5 AND org_id = (SELECT org_id FROM cte_get_cloud_ctx)
				  )
				  INSERT INTO orchestrations_scheduled_to_cloud_ctx_ns(orchestration_id, cloud_ctx_ns_id)
				  VALUES ((SELECT orchestration_id FROM cte_get_orchestration_id), (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx))
				  RETURNING orchestration_schedule_id
				  `
	log.Debug().Interface("InsertOrchestrationsScheduledToCloudCtxNs", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.CloudProvider, o.Region, o.Context, o.Namespace,
		o.OrchestrationName).Scan(&o.Scheduled.OrchestrationScheduleID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func (o *OrchestrationJob) InsertOrchestrationsScheduledToCloudCtxNs(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO orchestrations_scheduled_to_cloud_ctx_ns(orchestration_id, cloud_ctx_ns_id)
				  VALUES ($1, $2)
				  `
	log.Debug().Interface("InsertOrchestrationsScheduledToCloudCtxNs", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrchestrationID, o.Scheduled.CloudCtxNsID).Scan(&o.Scheduled.OrchestrationScheduleID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func (o *OrchestrationJob) UpdateOrchestrationsScheduledToCloudCtxNs(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				WITH cte_get_cloud_ctx AS (
					SELECT cloud_ctx_ns_id, org_id, cloud_provider, region, context, namespace
					FROM topologies_org_cloud_ctx_ns
					WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4
					LIMIT 1
				)
				UPDATE orchestrations_scheduled_to_cloud_ctx_ns
				SET status = $5
				FROM orchestrations
				WHERE orchestrations_scheduled_to_cloud_ctx_ns.orchestration_id = orchestrations.orchestration_id
				AND orchestrations.org_id = (SELECT org_id FROM cte_get_cloud_ctx)
				AND orchestrations.orchestration_name = $6
				AND orchestrations_scheduled_to_cloud_ctx_ns.cloud_ctx_ns_id = (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx)
				RETURNING orchestrations_scheduled_to_cloud_ctx_ns.orchestration_schedule_id;`
	log.Debug().Interface("UpdateOrchestrationsScheduledToCloudCtxNs", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.CloudProvider, o.Region, o.Context, o.Namespace,
		o.Scheduled.Status, o.OrchestrationName).Scan(&o.Scheduled.OrchestrationScheduleID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func (o *OrchestrationJob) SelectOrchestrationsAtCloudCtxNsWithStatus(ctx context.Context) (bool, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				WITH cte_get_cloud_ctx AS (
					SELECT cloud_ctx_ns_id, org_id, cloud_provider, region, context, namespace
					FROM topologies_org_cloud_ctx_ns
					WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4
					LIMIT 1
				)
				SELECT true
				FROM orchestrations_scheduled_to_cloud_ctx_ns os
				INNER JOIN orchestrations o
				ON os.orchestration_id = o.orchestration_id
				WHERE o.org_id = (SELECT org_id FROM cte_get_cloud_ctx)
				AND os.cloud_ctx_ns_id = (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx)
				AND os.status = $5
				AND o.orchestration_name = $6;`
	var orchestrationTodo bool
	log.Debug().Interface("SelectOrchestrations", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.CloudProvider, o.Region, o.Context, o.Namespace, o.Scheduled.Status, o.OrchestrationName).Scan(&orchestrationTodo)
	if err == pgx.ErrNoRows {
		// no rows were found by the query
		return false, nil
	} else if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		// an error occurred during the query execution
		return orchestrationTodo, err
	}
	// the query returned a row
	return orchestrationTodo, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

type AggregatedData struct {
	WorkflowResultID      int             `json:"workflowResultId"`
	ResponseID            int             `json:"responseId"`
	SourceTaskID          int             `json:"sourceTaskId"`
	TaskName              string          `json:"taskName"`
	TaskType              string          `json:"taskType"`
	Model                 string          `json:"model"`
	RunningCycleNumber    int             `json:"runningCycleNumber"`
	SearchWindowUnixStart int             `json:"searchWindowUnixStart"`
	SearchWindowUnixEnd   int             `json:"searchWindowUnixEnd"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`
	CompletionChoices     json.RawMessage `json:"completionChoices,omitempty"`
	PromptTokens          int             `json:"promptTokens"`
	CompletionTokens      int             `json:"completionTokens"`
	TotalTokens           int             `json:"totalTokens"`
}

type OrchestrationsAnalysis struct {
	TotalWorkflowTokenUsage int              `db:"total_workflow_token_usage" json:"totalWorkflowTokenUsage"`
	RunCycles               int              `db:"max_run_cycle" json:"runCycles"`
	AggregatedData          []AggregatedData `db:"aggregated_data" json:"aggregatedData"`

	artemis_autogen_bases.Orchestrations `json:"orchestrations,omitempty"`
}

func SelectAiSystemOrchestrations(ctx context.Context, orgID int) ([]OrchestrationsAnalysis, error) {
	var ojs []OrchestrationsAnalysis
	q := sql_query_templates.QueryParams{}

	// uses main for unique id, so type == real name for related workflow
	q.RawQuery = `SELECT orchestrations_id,
						 orchestrations.orchestration_name as orch_name,
						 orchestrations.group_name as orch_group_name,
						 orchestrations."type" as orch_type,
						 orchestrations.active as active,
						MAX(running_cycle_number) as max_run_cycle,
						SUM(cr.total_tokens) as total_workflow_token_usage,
						JSON_AGG(
							JSON_BUILD_OBJECT(
								'workflowResultId', workflow_result_id, 
								'responseId', ar.response_id, 
								'sourceTaskId', source_task_id,
								'taskName', ait.task_name,
								'taskType', ait.task_type,
								'model', ait.model,
								'runningCycleNumber', running_cycle_number, 
								'searchWindowUnixStart', search_window_unix_start, 
								'searchWindowUnixEnd', search_window_unix_end, 
								'metadata', metadata,
								'completionChoices', cr.completion_choices, 
								'promptTokens', cr.prompt_tokens, 
								'completionTokens', cr.completion_tokens, 
								'totalTokens', cr.total_tokens
							) ORDER BY running_cycle_number DESC, ar.response_id DESC
						) AS aggregated_data
				FROM 
					ai_workflow_analysis_results ar
				JOIN 
					ai_task_library ait ON ait.task_id = ar.source_task_id
				JOIN 
					completion_responses cr ON cr.response_id = ar.response_id
				JOIN 
					orchestrations ON orchestrations.orchestration_id = ar.orchestrations_id
				WHERE orchestrations.org_id = $1 
					AND (
						EXISTS (
							SELECT 1
							FROM ai_workflow_template 
							WHERE workflow_name = orchestrations.type
						) 
						OR EXISTS (
							SELECT 1
							FROM ai_workflow_template 
							WHERE workflow_group = orchestrations.group_name
						)
					)
				GROUP BY 
					orchestrations_id, orchestration_name, group_name, orch_type, active
				ORDER BY orchestrations_id DESC
				`

	log.Debug().Interface("SelectSystemOrchestrationsWithInstructionsByGroup", q.LogHeader(Orchestrations))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationsAnalysis{}
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.GroupName, &oj.Type, &oj.Active, &oj.RunCycles, &oj.TotalWorkflowTokenUsage, &oj.AggregatedData)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return nil, rowErr
		}
		ojs = append(ojs, oj)
	}
	return ojs, err
}
