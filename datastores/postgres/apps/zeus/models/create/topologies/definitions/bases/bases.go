package create_bases

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "System"

func InsertBaseQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertBaseDefinition"
	q.RawQuery = `INSERT INTO topology_base_components (org_id, topology_class_type_id, topology_system_component_id, topology_base_name)
			      VALUES ($1, $2, $3, $4)
			      ON CONFLICT DO NOTHING
			      RETURNING topology_base_component_id`
	return q
}

func InsertBaseQV2() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertBaseDefinitionV2"
	// topology_base_component_id
	q.RawQuery = `INSERT INTO topology_base_components (org_id, topology_class_type_id, topology_system_component_id, topology_base_name)
			      VALUES ($1, $2, $3, $4)
			      ON CONFLICT DO NOTHING
			      RETURNING topology_base_component_id`
	return q
}

func InsertBaseTx(ctx context.Context, base *bases.Base, tx pgx.Tx) (pgx.Tx, error) {
	q := InsertBaseQV2()
	log.Debug().Interface("InsertBase:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, base.OrgID, class_types.BaseClassTypeID, base.TopologySystemComponentID, base.TopologyBaseName).Scan(&base.TopologyBaseComponentID)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InsertBase: failed to insert base")
		return tx, fmt.Errorf("InsertBase: failed to insert base: %w", err)
	}
	return tx, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertBase(ctx context.Context, base *bases.Base) error {
	q := InsertBaseQ()
	log.Debug().Interface("InsertBase:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, base.OrgID, class_types.BaseClassTypeID, base.TopologySystemComponentID, base.TopologyBaseName).Scan(&base.TopologyBaseComponentID)
	if err == pgx.ErrNoRows {
		log.Ctx(ctx).Info().Msg("InsertBase: no rows returned, skipping (probably a duplicate row)")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertBasesQ(ctx context.Context, orgID int, clusterClassName string, cte *sql_query_templates.CTE, bs []bases.Base) sql_query_templates.SubCTEs {
	ts := chronos.Chronos{}
	getClusterIDByName := `SELECT topology_system_component_id FROM topology_system_components WHERE org_id = $1 AND topology_system_component_name = $2 AND topology_class_type_id = $3`
	var fetchClusterInfo = sql_query_templates.SubCTE{}
	fetchClusterInfo.QueryName = "cte_select_sys_component_id"
	fetchClusterInfo.RawQuery = getClusterIDByName
	cte.Params = []interface{}{orgID, clusterClassName, class_types.ClusterClassTypeID}
	basesCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insertBasesCTE_%d", ts.UnixTimeStampNow()))
	basesCTE.TableName = "topology_base_components"
	basesCTE.Columns = []string{"org_id", "topology_class_type_id", "topology_base_name", "topology_base_component_id", "topology_system_component_id"}
	basesCTE.ValuesOverride = make(map[int]string)
	for _, baseVal := range bs {
		basesCTE.AddValues(baseVal.OrgID, class_types.BaseClassTypeID, baseVal.TopologyBaseName, ts.UnixTimeStampNow(), "")
		basesCTE.ValuesOverride[4] = "(SELECT topology_system_component_id FROM cte_select_sys_component_id)"
	}
	return []sql_query_templates.SubCTE{fetchClusterInfo, basesCTE}
}

func InsertBases(ctx context.Context, orgID int, clusterClassName string, bs []bases.Base) error {
	cte := sql_query_templates.CTE{}
	q := InsertBasesQ(ctx, orgID, clusterClassName, &cte, bs)
	cte.AppendSubCtes(q)
	cte.OnConflictDoNothing = true
	query := cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, query, cte.Params...)
	if err != nil {
		return err
	}
	log.Ctx(ctx).Debug().Int64("rowsAffected", r.RowsAffected())
	return err
}

func InsertBasesTx(ctx context.Context, orgID int, clusterClassName string, bs []bases.Base, tx pgx.Tx) (pgx.Tx, error) {
	cte := sql_query_templates.CTE{}
	q := InsertBasesQ(ctx, orgID, clusterClassName, &cte, bs)
	cte.AppendSubCtes(q)
	cte.OnConflictDoNothing = true
	query := cte.GenerateChainedCTE()
	r, err := tx.Exec(ctx, query, cte.Params...)
	if err != nil {
		return tx, err
	}
	log.Ctx(ctx).Debug().Int64("rowsAffected", r.RowsAffected())
	return tx, err
}
