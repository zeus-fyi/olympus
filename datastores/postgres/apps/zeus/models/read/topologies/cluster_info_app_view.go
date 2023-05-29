package read_topologies

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func readClusterAppViewSummaryQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})

	query := `SELECT MAX(td.topology_id) as topology_id, topology_system_component_name, tctx.cloud_ctx_ns_id, kns.cloud_provider, kns.region, kns.context, kns.namespace, tctx.namespace_alias
				FROM topologies_deployed td
				JOIN topologies_kns kns ON kns.topology_id = td.topology_id
				JOIN org_users_topologies out ON out.topology_id = td.topology_id
				JOIN topology_infrastructure_components tip ON tip.topology_id = td.topology_id
				JOIN topology_skeleton_base_components sb ON sb.topology_skeleton_base_id = tip.topology_skeleton_base_id
				JOIN topology_base_components bc ON bc.topology_base_component_id = sb.topology_base_component_id
				JOIN topology_system_components tsys ON tsys.topology_system_component_id = bc.topology_system_component_id
				JOIN topologies_org_cloud_ctx_ns tctx ON
					tctx.namespace = kns.namespace AND
					tctx.cloud_provider = kns.cloud_provider AND
					tctx.context = kns.context AND
					tctx.region = kns.region
				WHERE out.org_id = $1 AND td.topology_status != 'DestroyDeployComplete'
				GROUP BY topology_system_component_name, kns.cloud_provider, kns.context, kns.region, kns.namespace,  tctx.namespace_alias, tctx.cloud_ctx_ns_id
				ORDER BY topology_system_component_name, kns.cloud_provider, kns.context, kns.region, kns.namespace,  tctx.namespace_alias, tctx.cloud_ctx_ns_id
				LIMIT 1000`
	q.RawQuery = query
	return q
}

func readClusterAppViewSummaryQueryForApp() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})
	query := `SELECT MAX(td.topology_id) as topology_id, topology_system_component_name, tctx.cloud_ctx_ns_id, kns.cloud_provider, kns.region, kns.context, kns.namespace, tctx.namespace_alias
				FROM topologies_deployed td
				JOIN topologies_kns kns ON kns.topology_id = td.topology_id
				JOIN org_users_topologies out ON out.topology_id = td.topology_id
				JOIN topology_infrastructure_components tip ON tip.topology_id = td.topology_id
				JOIN topology_skeleton_base_components sb ON sb.topology_skeleton_base_id = tip.topology_skeleton_base_id
				JOIN topology_base_components bc ON bc.topology_base_component_id = sb.topology_base_component_id
				JOIN topology_system_components tsys ON tsys.topology_system_component_id = bc.topology_system_component_id
				JOIN topologies_org_cloud_ctx_ns tctx ON
					tctx.namespace = kns.namespace AND
					tctx.cloud_provider = kns.cloud_provider AND
					tctx.context = kns.context AND
					tctx.region = kns.region
				WHERE out.org_id = $1 AND td.topology_status != 'DestroyDeployComplete' AND tsys.topology_system_component_name = $2
				GROUP BY topology_system_component_name, kns.cloud_provider, kns.context, kns.region, kns.namespace,  tctx.namespace_alias, tctx.cloud_ctx_ns_id
				ORDER BY topology_system_component_name, kns.cloud_provider, kns.context, kns.region, kns.namespace,  tctx.namespace_alias, tctx.cloud_ctx_ns_id
				LIMIT 1000`
	q.RawQuery = query
	return q
}

type ClusterAppView struct {
	ClusterClassName string `json:"clusterClassName"`
	autogen_bases.TopologiesOrgCloudCtxNs
}

func SelectClusterAppView(ctx context.Context, orgID int) ([]ClusterAppView, error) {
	var TopologiesOrgCloudCtxNs []ClusterAppView
	q := readClusterAppViewSummaryQuery()
	log.Debug().Interface("SelectClusterAppView:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tpCtx := ClusterAppView{}
		var noOp int
		rowErr := rows.Scan(
			&noOp, &tpCtx.ClusterClassName, &tpCtx.CloudCtxNsID, &tpCtx.CloudProvider, &tpCtx.Region, &tpCtx.Context, &tpCtx.Namespace, &tpCtx.NamespaceAlias,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return TopologiesOrgCloudCtxNs, rowErr
		}
		TopologiesOrgCloudCtxNs = append(TopologiesOrgCloudCtxNs, tpCtx)
	}
	return TopologiesOrgCloudCtxNs, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func SelectClusterSingleAppView(ctx context.Context, orgID int, appName string) ([]ClusterAppView, error) {
	var TopologiesOrgCloudCtxNs []ClusterAppView
	q := readClusterAppViewSummaryQueryForApp()
	log.Debug().Interface("SelectClusterSingleAppView:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, appName)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tpCtx := ClusterAppView{}
		var noOp int
		rowErr := rows.Scan(
			&noOp, &tpCtx.ClusterClassName, &tpCtx.CloudCtxNsID, &tpCtx.CloudProvider, &tpCtx.Region, &tpCtx.Context, &tpCtx.Namespace, &tpCtx.NamespaceAlias,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return TopologiesOrgCloudCtxNs, rowErr
		}
		TopologiesOrgCloudCtxNs = append(TopologiesOrgCloudCtxNs, tpCtx)
	}
	return TopologiesOrgCloudCtxNs, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
