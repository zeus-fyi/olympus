package create_systems

import (
	"context"

	"github.com/jackc/pgx/v4"
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
			      ON CONFLICT DO NOTHING
			      RETURNING topology_system_component_id`
	return q
}

func InsertSystem(ctx context.Context, system *systems.Systems) error {
	q := InsertSystemQ()
	log.Debug().Interface("InsertSystem:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, system.OrgID, system.TopologyClassTypeID, system.TopologySystemComponentName).Scan(&system.TopologySystemComponentID)
	if err == pgx.ErrNoRows {
		log.Ctx(ctx).Info().Msg("InsertSystem: no rows returned, skipping (probably a duplicate row)")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertSystemTx(ctx context.Context, system *systems.Systems, tx pgx.Tx) (pgx.Tx, error) {
	q := InsertSystemQ()
	log.Debug().Interface("InsertSystem:", q.LogHeader(Sn))
	err := tx.QueryRow(ctx, q.RawQuery, system.OrgID, system.TopologyClassTypeID, system.TopologySystemComponentName).Scan(&system.TopologySystemComponentID)
	if err == pgx.ErrNoRows {
		log.Ctx(ctx).Info().Msg("InsertSystem: no rows returned, skipping (probably a duplicate row)")
		return tx, nil
	}
	return tx, misc.ReturnIfErr(err, q.LogHeader(Sn))
}
