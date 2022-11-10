package create_topology_deployment_status

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "TopologyKubeCtxNs"

func (s *DeploymentStatus) insertTopologyStatus() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("insertTopologyStatus", "topologies_deployed", "where", 1000, []string{})
	query := `INSERT INTO topologies_deployed(topology_id, topology_status) 
			  VALUES ($1, $2) 
			  ON CONFLICT (topology_id, topology_status) DO UPDATE SET topology_status = EXCLUDED.topology_status
			  RETURNING updated_at`
	q.RawQuery = query
	q.CTEQuery.Params = []interface{}{s.TopologyID, s.TopologyStatus}
	return q
}

func (s *DeploymentStatus) InsertStatus(ctx context.Context) error {
	q := s.insertTopologyStatus()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, q.CTEQuery.Params...).Scan(&s.UpdatedAt)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
