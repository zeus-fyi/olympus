package read_topology_deployment_status

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
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
			  `
	q.RawQuery = query
	return q
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
			&rs.Kns.TopologyID, &rs.TopologyName,
			&rs.TopologyStatus, &rs.UpdatedAt,
			&rs.Kns.CloudProvider, &rs.Kns.Region, &rs.Kns.Context, &rs.Kns.Namespace, &rs.Kns.Env,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return rowErr
		}
		s.Slice = append(s.Slice, rs)
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
