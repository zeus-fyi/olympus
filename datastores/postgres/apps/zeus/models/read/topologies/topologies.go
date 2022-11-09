package read_topologies

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

type Topologies struct {
	topology.Topology
}

type ReadTopologiesMetadata struct {
	TopologyID       int            `db:"topology_id" json:"topology_id"`
	TopologyName     string         `db:"topology_name" json:"topology_name"`
	ChartName        string         `db:"chart_name" json:"chart_name"`
	ChartVersion     string         `db:"chart_version" json:"chart_version"`
	ChartDescription sql.NullString `db:"chart_description" json:"chart_description"`
}

type ReadTopologiesMetadataGroup struct {
	Slice []ReadTopologiesMetadata
}

func NewReadTopologiesMetadataGroup() ReadTopologiesMetadataGroup {
	return ReadTopologiesMetadataGroup{Slice: []ReadTopologiesMetadata{}}
}

const Sn = "ReadTopologiesMetadataGroup"

func (t *ReadTopologiesMetadataGroup) defaultQ() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("ReadState", "read_topology_metadata", "where", 1000, []string{})
	query := `SELECT t.topology_id, t.name, cp.chart_name, cp.chart_version, cp.chart_description
			  FROM org_users_topologies out
			  INNER JOIN topologies t ON t.topology_id = out.topology_id
			  LEFT JOIN topology_infrastructure_components ti ON ti.topology_id = out.topology_id
			  INNER JOIN chart_packages cp ON cp.chart_package_id = ti.chart_package_id
			  WHERE out.org_id = $1 AND out.user_id = $2
			  ORDER BY t.topology_id DESC
			  `
	q.RawQuery = query
	return q
}
func (t *ReadTopologiesMetadataGroup) SelectTopologiesMetadata(ctx context.Context, ou org_users.OrgUser) error {
	t.Slice = []ReadTopologiesMetadata{}

	q := t.defaultQ()
	log.Debug().Interface("SelectTopologiesMetadata:", q.LogHeader(Sn))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, ou.OrgID, ou.UserID)
	if err != nil {
		log.Err(err).Msg(q.LogHeader(Sn))
		return err
	}
	defer rows.Close()
	for rows.Next() {
		tm := ReadTopologiesMetadata{}
		rowErr := rows.Scan(
			&tm.TopologyID, &tm.TopologyName, &tm.ChartName, &tm.ChartVersion, &tm.ChartDescription,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(Sn))
			return rowErr
		}
		t.Slice = append(t.Slice, tm)
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
