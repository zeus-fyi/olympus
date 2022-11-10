package create_topology_deployment_status

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "TopologyKubeCtxNs"

func (s *DeploymentStatus) defaultQ() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("InsertState", "topologies_deployed", "where", 1000, []string{})
	q.TableName = s.GetTableName()
	q.Columns = s.GetTableColumns()
	q.Values = []apps.RowValues{s.GetRowValues("default")}

	query := `INSERT INTO topologies_deployed(topology_id, topology_status) 
			  VALUES ($1, $2) 
			  ON CONFLICT ON CONSTRAINT topologies_deployed_topology_id_fkey DO UPDATE SET topology_id = EXCLUDED.topology_id 
			  RETURNING updated_at`
	q.RawQuery = query
	q.CTEQuery.Params = []interface{}{s.TopologyID, s.TopologyStatus}
	return q
}

func (s *DeploymentStatus) InsertStatus(ctx context.Context) error {
	q := s.defaultQ()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, q.CTEQuery.Params...).Scan(&s.UpdatedAt)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
