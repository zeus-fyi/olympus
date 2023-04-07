package create_skeletons

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "SkeletonBase"

func insertSkeletonBaseQ() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertSkeletonBaseDefinition"
	q.RawQuery = `INSERT INTO topology_skeleton_base_components (org_id, topology_base_component_id, topology_class_type_id, topology_skeleton_base_name)
			      VALUES ($1, $2, $3, $4)
			      ON CONFLICT DO NOTHING
				  RETURNING topology_skeleton_base_id`
	return q
}
func InsertSkeletonBaseTx(ctx context.Context, skeletonBase *skeletons.SkeletonBase, tx pgx.Tx) (pgx.Tx, error) {
	q := insertSkeletonBaseQ()
	log.Debug().Interface("InsertSkeletonBase:", q.LogHeader(Sn))
	err := tx.QueryRow(ctx, q.RawQuery, skeletonBase.OrgID, skeletonBase.TopologyBaseComponentID, class_types.SkeletonBaseClassTypeID, skeletonBase.TopologySkeletonBaseName).Scan(&skeletonBase.TopologySkeletonBaseID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg(fmt.Sprintf("InsertSkeletonBase: %s", q.LogHeader(Sn)))
		return nil, err
	}
	return tx, misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func InsertSkeletonBase(ctx context.Context, skeletonBase *skeletons.SkeletonBase) error {
	q := insertSkeletonBaseQ()
	log.Debug().Interface("InsertSkeletonBase:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, skeletonBase.OrgID, skeletonBase.TopologyBaseComponentID, class_types.SkeletonBaseClassTypeID, skeletonBase.TopologySkeletonBaseName).Scan(&skeletonBase.TopologySkeletonBaseID)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func insertSkeletonBasesQ(ctx context.Context, orgID int, clusterClassName, componentBaseName string, skeletonBaseNames []string, cte *sql_query_templates.CTE) sql_query_templates.SubCTEs {
	ts := chronos.Chronos{}
	getClusterIDByName := `SELECT topology_base_component_id FROM topology_base_components
                           WHERE org_id = $1 AND topology_base_name = $2 AND topology_system_component_id IN (SELECT topology_system_component_id
                                                                                                              FROM topology_system_components
                                                                                                              WHERE org_id = $3 AND topology_system_component_name = $4) `
	var fetchClusterInfo = sql_query_templates.SubCTE{}
	fetchClusterInfo.QueryName = "cte_select_base_component_id"
	fetchClusterInfo.RawQuery = getClusterIDByName
	cte.Params = []interface{}{orgID, componentBaseName, orgID, clusterClassName}
	basesCTE := sql_query_templates.NewSubInsertCTE(fmt.Sprintf("cte_insertBasesCTE_%d", ts.UnixTimeStampNow()))
	basesCTE.TableName = "topology_skeleton_base_components"
	basesCTE.Columns = []string{"org_id", "topology_class_type_id", "topology_skeleton_base_name", "topology_skeleton_base_id", "topology_base_component_id"}
	basesCTE.ValuesOverride = make(map[int]string)
	for _, n := range skeletonBaseNames {
		basesCTE.AddValues(orgID, class_types.SkeletonBaseClassTypeID, n, ts.UnixTimeStampNow(), "")
		basesCTE.ValuesOverride[4] = "(SELECT topology_base_component_id FROM cte_select_base_component_id)"
	}

	return []sql_query_templates.SubCTE{fetchClusterInfo, basesCTE}
}

func InsertSkeletonBases(ctx context.Context, orgID int, clusterClassName string, componentBaseName string, skeletonBaseNames []string) error {
	cte := sql_query_templates.CTE{}
	q := insertSkeletonBasesQ(ctx, orgID, clusterClassName, componentBaseName, skeletonBaseNames, &cte)
	cte.AppendSubCtes(q)
	cte.OnConflictDoNothing = true
	query := cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, query, cte.Params...)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InsertSkeletonBases")
		return err
	}
	log.Ctx(ctx).Debug().Int64("rowsAffected", r.RowsAffected())
	return err
}

func InsertSkeletonBasesTx(ctx context.Context, orgID int, clusterClassName string, componentBaseName string, skeletonBaseNames []string, tx pgx.Tx) (pgx.Tx, error) {
	cte := sql_query_templates.CTE{}
	q := insertSkeletonBasesQ(ctx, orgID, clusterClassName, componentBaseName, skeletonBaseNames, &cte)
	cte.AppendSubCtes(q)
	cte.OnConflictDoNothing = true
	query := cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, query, cte.Params...)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("InsertSkeletonBases")
		return tx, err
	}
	log.Ctx(ctx).Debug().Int64("rowsAffected", r.RowsAffected())
	return tx, err
}
