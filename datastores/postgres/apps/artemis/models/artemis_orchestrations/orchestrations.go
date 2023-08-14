package artemis_orchestrations

import (
	"context"

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
	artemis_autogen_bases.Orchestrations
	Scheduled artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs
	zeus_common_types.CloudCtxNs
}

func (o *OrchestrationJob) InsertOrchestrations(ctx context.Context) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO orchestrations(org_id, orchestration_name)
				  VALUES ($1, $2)
				  RETURNING orchestration_id;
				  `
	log.Debug().Interface("InsertOrchestrations", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, o.OrchestrationName).Scan(&o.OrchestrationID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func SelectActiveOrchestrationsWithInstructions(ctx context.Context, orgID int, orchestType, groupName string) ([]OrchestrationJob, error) {
	var ojs []OrchestrationJob
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT orchestration_id, orchestration_name, instructions
				  FROM orchestrations
				  WHERE org_id = $1 AND active = true AND type = $2 AND group_name = $3
				  `
	log.Debug().Interface("InsertOrchestrations", q.LogHeader(Orchestrations))

	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, orchestType, groupName)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return ojs, err
	}
	defer rows.Close()
	for rows.Next() {
		oj := OrchestrationJob{}
		rowErr := rows.Scan(&oj.OrchestrationID, &oj.OrchestrationName, &oj.Instructions)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Orchestrations))
			return ojs, rowErr
		}
		ojs = append(ojs, oj)
	}
	return ojs, err
}

func (o *OrchestrationJob) InsertOrchestrationsWithInstructions(ctx context.Context, instructions []byte) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO orchestrations(org_id, orchestration_name, instructions)
				  VALUES ($1, $2, $3)
				  RETURNING orchestration_id;
				  `
	log.Debug().Interface("InsertOrchestrations", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, o.OrgID, o.OrchestrationName, instructions).Scan(&o.OrchestrationID)
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
