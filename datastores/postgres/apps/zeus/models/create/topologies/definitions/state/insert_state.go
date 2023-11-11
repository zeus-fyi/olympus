package create_topology_deployment_status

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	topology_deployment_status "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/state"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "TopologyDeploymentStatus"

func insertTopologyStatus() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("insertTopologyStatus", "topologies_deployed", "where", 1000, []string{})
	query := `INSERT INTO topologies_deployed(topology_id, topology_status) 
			  VALUES ($1, $2) 
			  RETURNING deployment_id, updated_at`
	q.RawQuery = query
	return q
}

func InsertOrUpdateStatus(ctx context.Context, status *topology_deployment_status.DeployStatus) error {
	if status == nil {
		return nil
	}
	q := insertTopologyStatus()
	log.Debug().Interface("InsertOrUpdateQuery:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, status.TopologyID, status.TopologyStatus).Scan(&status.DeploymentID, &status.UpdatedAt)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
