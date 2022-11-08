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
	return q
}

func (s *DeploymentStatus) InsertStatus(ctx context.Context) error {
	q := s.defaultQ()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	r, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("State: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
