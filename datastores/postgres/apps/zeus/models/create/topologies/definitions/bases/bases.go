package create_bases

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "System"

func InsertBaseQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertBaseDefinition"
	q.RawQuery = `INSERT INTO topology_base_components (org_id, topology_class_type_id, topology_system_component_id, topology_base_name)
			      VALUES ($1, $2, $3, $4)
			      RETURNING topology_base_component_id`
	return q
}

func InsertBase(ctx context.Context, base *bases.Base) error {
	q := InsertBaseQ()
	log.Debug().Interface("InsertBase:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, base.OrgID, class_types.BaseClassTypeID, base.TopologySystemComponentID, base.TopologyBaseName).Scan(&base.TopologyBaseComponentID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
