package read_topologies

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func readClusterSummaryQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})

	query := `SELECT tp.cloud_ctx_ns_id, tp.cloud_provider, tp.region, tp.context, tp.namespace, tp.namespace_alias, tp.created_at
			  FROM topologies_org_cloud_ctx_ns tp
			  WHERE tp.org_id = $1 
			  ORDER BY tp.cloud_provider, tp.region, tp.context, tp.namespace_alias, tp.created_at DESC
			  LIMIT 1000
			  `
	q.RawQuery = query
	return q
}

func SelectTopologiesMetadata(ctx context.Context, orgID int) (autogen_bases.TopologiesOrgCloudCtxNsSlice, error) {
	TopologiesOrgCloudCtxNs := autogen_bases.TopologiesOrgCloudCtxNsSlice{}
	q := readClusterSummaryQuery()
	log.Debug().Interface("SelectTopologiesMetadata:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return TopologiesOrgCloudCtxNs, err
	}
	defer rows.Close()
	for rows.Next() {
		tpCtx := autogen_bases.TopologiesOrgCloudCtxNs{}
		rowErr := rows.Scan(
			&tpCtx.CloudCtxNsID, &tpCtx.CloudProvider, &tpCtx.Region, &tpCtx.Context, &tpCtx.Namespace, &tpCtx.NamespaceAlias, &tpCtx.CreatedAt,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return TopologiesOrgCloudCtxNs, rowErr
		}
		TopologiesOrgCloudCtxNs = append(TopologiesOrgCloudCtxNs, tpCtx)
	}
	return TopologiesOrgCloudCtxNs, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
