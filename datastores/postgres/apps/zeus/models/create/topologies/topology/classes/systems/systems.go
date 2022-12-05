package create_systems

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "System"

func InsertSystemQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertSystemDefinition"
	q.RawQuery = `INSERT INTO topology_system_components (org_id, topology_class_type_id, topology_system_component_name)
			      VALUES ($1, $2, $3)
			      RETURNING topology_system_component_id`
	return q
}

func InsertSystem(ctx context.Context, system *systems.Systems) error {
	q := InsertSystemQ()
	log.Debug().Interface("InsertSystem:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, system.OrgID, system.TopologyClassTypeID, system.TopologySystemComponentName).Scan(&system.TopologySystemComponentID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
