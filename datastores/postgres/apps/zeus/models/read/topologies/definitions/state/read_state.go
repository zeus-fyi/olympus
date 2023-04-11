package read_topology_deployment_status

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

const Sn = "ReadDeploymentStatusesGroup"

func (s *ReadDeploymentStatusesGroup) defaultQ() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})

	query := `SELECT tp.topology_id, t.name, tp.topology_status, tp.updated_at, 
       		  COALESCE (tkns.cloud_provider, '') AS cloud_provider,
       		  COALESCE (tkns.region, '') AS region,
       		  COALESCE (tkns.context, '') AS context,
       		  COALESCE (tkns.namespace, '') AS namespace,
       		  COALESCE (tkns.env, '')  AS env
			  FROM topologies_deployed tp
			  INNER JOIN topologies t ON t.topology_id = tp.topology_id
			  INNER JOIN org_users_topologies out ON out.topology_id = tp.topology_id
			  LEFT JOIN topologies_kns tkns ON tkns.topology_id = tp.topology_id
			  WHERE tp.topology_id = $1 AND out.org_id = $2 AND out.user_id = $3
			  ORDER BY tp.updated_at DESC
			  LIMIT 1000
			  `
	q.RawQuery = query
	return q
}

func (s *ReadDeploymentStatusesGroup) readTopologiesDeployedStatusToKns() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})

	query := `SELECT MAX(td.topology_id) as topology_id, topology_skeleton_base_name, td.topology_status
				FROM topologies_deployed td
				JOIN topologies_kns kns ON kns.topology_id = td.topology_id
				JOIN org_users_topologies out ON out.topology_id = td.topology_id
				JOIN topology_infrastructure_components tip ON tip.topology_id = td.topology_id
				JOIN topology_skeleton_base_components sb ON sb.topology_skeleton_base_id = tip.topology_skeleton_base_id
				WHERE kns.cloud_provider = $1 AND kns.region = $2 AND kns.context = $3 AND kns.namespace = $4 AND out.org_id = $5
				GROUP BY topology_skeleton_base_name, topology_status
				ORDER BY topology_skeleton_base_name, topology_status
			  	LIMIT 1000
			  `
	q.RawQuery = query
	return q
}

func (s *ReadDeploymentStatusesGroup) readClustersDeployedToKns() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_deployment_status", "where", 1000, []string{})

	query := `WITH cte_get_org_ctx AS (
					SELECT cloud_provider, region, context, namespace, namespace_alias
					FROM topologies_org_cloud_ctx_ns oc 
					WHERE cloud_ctx_ns_id = $1 AND org_id = $2
					LIMIT 1
				) 
				SELECT MAX(td.topology_id) as topology_id, topology_system_component_name, topology_base_name, topology_skeleton_base_name
				FROM topologies_deployed td
				JOIN topologies_kns kns ON kns.topology_id = td.topology_id
				JOIN org_users_topologies out ON out.topology_id = td.topology_id
				JOIN topology_infrastructure_components tip ON tip.topology_id = td.topology_id
				JOIN topology_skeleton_base_components sb ON sb.topology_skeleton_base_id = tip.topology_skeleton_base_id
				JOIN topology_base_components bc ON bc.topology_base_component_id = sb.topology_base_component_id
				JOIN topology_system_components tsys ON tsys.topology_system_component_id = bc.topology_system_component_id
				WHERE kns.cloud_provider = (SELECT cloud_provider FROM cte_get_org_ctx) AND kns.region = (SELECT region FROM cte_get_org_ctx) AND kns.context = (SELECT context FROM cte_get_org_ctx) AND kns.namespace = (SELECT namespace FROM cte_get_org_ctx) AND out.org_id = $2 AND td.topology_status != 'DestroyDeployComplete'
				GROUP BY topology_system_component_name, topology_base_name, topology_skeleton_base_name
				ORDER BY topology_system_component_name, topology_base_name, topology_skeleton_base_name
			  	LIMIT 1000
			  `
	q.RawQuery = query
	return q
}

func (s *ReadDeploymentStatusesGroup) ReadLatestDeployedClusterTopologies(ctx context.Context, cloudCtxNsID int, ou org_users.OrgUser) error {
	s.Slice = []ReadDeploymentStatus{}
	q := s.readClustersDeployedToKns()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cloudCtxNsID, ou.OrgID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return err
	}
	defer rows.Close()
	for rows.Next() {

		namespaceAlias := ""
		rs := ReadDeploymentStatus{}
		rowErr := rows.Scan(
			&rs.TopologyID, &rs.ClusterName, &rs.ComponentBaseName, &rs.SkeletonBaseName, &namespaceAlias,
		)
		rs.UpdatedAt = time.Unix(0, int64(rs.TopologyID))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return rowErr
		}
		s.Slice = append(s.Slice, rs)
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (s *ReadDeploymentStatusesGroup) ReadLatestDeployedStatus(ctx context.Context, cctx zeus_common_types.CloudCtxNs, ou org_users.OrgUser) error {
	s.Slice = []ReadDeploymentStatus{}
	q := s.readTopologiesDeployedStatusToKns()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, cctx.CloudProvider, cctx.Region, cctx.Context, cctx.Namespace, ou.OrgID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rs := ReadDeploymentStatus{}
		rowErr := rows.Scan(
			&rs.TopologyID, &rs.TopologyName, &rs.TopologyStatus,
		)
		rs.UpdatedAt = time.Unix(0, int64(rs.TopologyID))
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return rowErr
		}
		s.Slice = append(s.Slice, rs)
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (s *ReadDeploymentStatusesGroup) ReadStatus(ctx context.Context, topologyID int, ou org_users.OrgUser) error {
	s.Slice = []ReadDeploymentStatus{}
	q := s.defaultQ()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, topologyID, ou.OrgID, ou.UserID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return err
	}
	defer rows.Close()
	for rows.Next() {
		rs := ReadDeploymentStatus{}
		rowErr := rows.Scan(
			&rs.TopologyID, &rs.TopologyName,
			&rs.TopologyStatus, &rs.UpdatedAt,
			&rs.CloudCtxNs.CloudProvider, &rs.CloudCtxNs.Region, &rs.CloudCtxNs.Context, &rs.CloudCtxNs.Namespace, &rs.CloudCtxNs.Env,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return rowErr
		}
		s.Slice = append(s.Slice, rs)
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
