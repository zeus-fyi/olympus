package delete_kns

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

const Sn = "DeleteTopologyKubeCtxNs"

func getDeleteKnsQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteKns"

	q.RawQuery = `
	WITH cte_delete_kns AS (
		DELETE FROM topologies_kns 
		WHERE topology_id = $1 AND cloud_provider = $2 AND region = $3 AND context = $4 AND namespace = $5 AND env = $6
	) SELECT true
	`
	return q
}

func getDeleteKnsQueryByOrgMatchAndCloudCtx() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteKnsUsingLookup"

	q.RawQuery = `
	WITH cte_get_authed_ctx AS (
		SELECT org_id, cloud_provider, region, context, namespace
		FROM topologies_org_cloud_ctx_ns
		WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4
		LIMIT 1
	), cte_delete_kns AS (
		DELETE FROM topologies_kns 
		WHERE cloud_provider = (SELECT cloud_provider FROM cte_get_authed_ctx) AND region = (SELECT region FROM cte_get_authed_ctx) AND namespace = (SELECT namespace FROM cte_get_authed_ctx) AND env = $5
	) SELECT true
	`
	return q
}

func DeleteKnsByOrgAccessAndCloudCtx(ctx context.Context, kns *zeus_req_types.TopologyDeployRequest) error {
	if kns == nil {
		return nil
	}
	q := getDeleteKnsQueryByOrgMatchAndCloudCtx()
	log.Debug().Interface("DeleteQuery:", q.LogHeader(Sn))
	success := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, kns.CloudProvider, kns.Region, kns.Context, kns.Namespace, kns.Env).Scan(&success)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func DeleteKns(ctx context.Context, kns *zeus_req_types.TopologyDeployRequest) error {
	if kns == nil {
		return nil
	}
	q := getDeleteKnsQuery()
	log.Debug().Interface("DeleteQuery:", q.LogHeader(Sn))
	success := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, kns.TopologyID, kns.CloudProvider, kns.Region, kns.Context, kns.Namespace, kns.Env).Scan(&success)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
