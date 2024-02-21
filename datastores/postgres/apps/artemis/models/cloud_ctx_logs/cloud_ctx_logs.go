package cloud_ctx_logs

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type CloudCtxNsLogs struct {
	LogID           int               `json:"logID"`
	OrchestrationID int               `json:"orchestrationID"`
	CloudCtxNsID    int               `json:"cloudCtxNsID"`
	Status          string            `json:"status"`
	Msg             string            `json:"msg"`
	Ou              org_users.OrgUser `json:"ou"`
	zeus_common_types.CloudCtxNs
}

func InsertCloudCtxNsLog(ctx context.Context, cl *CloudCtxNsLogs) error {
	if cl == nil {
		return nil
	}
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				WITH cte_get_cloud_ctx AS (
					SELECT cloud_ctx_ns_id, org_id, cloud_provider, region, context, namespace
					FROM topologies_org_cloud_ctx_ns
					WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4 AND org_id = $5
					LIMIT 1
				  )
				  INSERT INTO orchestrations_cloud_ctx_ns_logs(orchestration_id, cloud_ctx_ns_id, msg, status)
				  VALUES ($6, (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx), $7, $8)
				  RETURNING log_id;
				  `
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, cl.CloudProvider, cl.Region, cl.Context, cl.Namespace, cl.Ou.OrgID, cl.OrchestrationID, cl.Msg, cl.Status).Scan(&cl.LogID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("CloudCtxNsLogs")); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("CloudCtxNsLogs"))
}

func SelectCloudCtxNsLogs(ctx context.Context, cl CloudCtxNsLogs) ([]string, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
					WITH cte_get_cloud_ctx AS (
						SELECT cloud_ctx_ns_id, org_id, cloud_provider, region, context, namespace
						FROM topologies_org_cloud_ctx_ns
						WHERE org_id = $1 AND cloud_ctx_ns_id = $2
					LIMIT 1
				  )
				  SELECT log_id, status, msg
				  FROM orchestrations_cloud_ctx_ns_logs
				  WHERE cloud_ctx_ns_id = (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx)
				  ORDER BY log_id DESC
				  LIMIT 100;`
	var logs []string
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cl.Ou.OrgID, cl.CloudCtxNsID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("CloudCtxNsLogs")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var logID int
		var status, msg string
		rowErr := rows.Scan(&logID, &status, &msg)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("CloudCtxNsLogs"))
			return nil, rowErr
		}
		logs = append(logs, fmt.Sprintf("%d | %s | %s", logID, status, msg))
	}
	return logs, misc.ReturnIfErr(err, q.LogHeader("CloudCtxNsLogs"))
}
