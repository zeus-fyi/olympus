package artemis_validator_service_groups_models

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const (
	Orchestrations = "Orchestrations"
)

func InsertOrchestrations(ctx context.Context, o *artemis_autogen_bases.Orchestrations) error {
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

func InsertOrchestrationsScheduledToCloudCtxNs(ctx context.Context, os *artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO orchestrations_scheduled_to_cloud_ctx_ns(orchestration_id, cloud_ctx_ns_id)
				  VALUES ($1, $2)
				  `
	log.Debug().Interface("InsertOrchestrationsScheduledToCloudCtxNs", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, os.OrchestrationID, os.CloudCtxNsID).Scan(&os.OrchestrationScheduleID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func UpdateOrchestrationsScheduledToCloudCtxNs(ctx context.Context, orgID int, name string, os *artemis_autogen_bases.OrchestrationsScheduledToCloudCtxNs) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				UPDATE orchestrations_scheduled_to_cloud_ctx_ns
				SET status = $1
				FROM orchestrations
				WHERE orchestrations_scheduled_to_cloud_ctx_ns.orchestration_id = orchestrations.orchestration_id
				AND orchestrations.org_id = $2
				AND orchestrations.orchestration_name = $3
				AND orchestrations_scheduled_to_cloud_ctx_ns.status = $4
				AND orchestrations_scheduled_to_cloud_ctx_ns.cloud_ctx_ns_id = $5;
				  `
	log.Debug().Interface("UpdateOrchestrationsScheduledToCloudCtxNs", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, os.Status, orgID, name, os.Status, os.OrchestrationID, os.CloudCtxNsID).Scan(&os.OrchestrationScheduleID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}

func SelectOrchestrationsAtCloudCtxNsWithStatus(ctx context.Context, orgID, cloudCtxNsID int, status, name string) (bool, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				SELECT true
				FROM orchestrations_scheduled_to_cloud_ctx_ns os
				INNER JOIN orchestrations o
				ON os.orchestration_id = o.orchestration_id
				WHERE o.org_id = $1
				AND os.cloud_ctx_ns_id = $3
				AND os.status = $2
				AND o.orchestration_name = $4;`
	var orchestrationTodo bool
	log.Debug().Interface("SelectOrchestrations", q.LogHeader(Orchestrations))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, cloudCtxNsID, status, name).Scan(&orchestrationTodo)
	if err == sql.ErrNoRows {
		// no rows were found by the query
		return false, nil
	} else if returnErr := misc.ReturnIfErr(err, q.LogHeader(Orchestrations)); returnErr != nil {
		// an error occurred during the query execution
		return orchestrationTodo, err
	}
	// the query returned a row
	return orchestrationTodo, misc.ReturnIfErr(err, q.LogHeader(Orchestrations))
}
