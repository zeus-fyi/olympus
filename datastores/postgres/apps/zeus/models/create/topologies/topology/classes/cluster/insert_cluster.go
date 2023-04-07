package create_clusters

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	create_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases"
	create_skeletons "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/skeleton"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

const Sn = "Cluster"

var (
	checkExistingCluster             = `SELECT topology_system_component_id FROM topology_system_components WHERE topology_system_component_name = $1 AND org_id = $2`
	checkExistingClusterBase         = `SELECT topology_base_component_id FROM topology_base_components WHERE topology_base_name = $1 AND org_id = $2`
	checkExistingClusterSkeletonBase = `SELECT topology_skeleton_base_id FROM topology_skeleton_base_components WHERE topology_skeleton_base_name = $1 AND topology_base_component_id = $2 AND org_id = $3`
)

func InsertCluster(ctx context.Context, tx pgx.Tx, sys *systems.Systems, cbMap zeus_templates.ClusterPreviewWorkloads) (pgx.Tx, error) {
	// Check if the cluster already exists
	err := tx.QueryRow(ctx, checkExistingCluster, sys.TopologySystemComponentName, sys.OrgID).Scan(&sys.TopologySystemComponentID)
	if err == pgx.ErrNoRows {
		err = nil
		tx, err = create_systems.InsertSystemTx(ctx, sys, tx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
			return tx, err
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
			return tx, err
		}
		tx, err = apps.Pg.Begin(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to start tx")
			return tx, err
		}
	}
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
		return nil, err
	}
	for cbName, cb := range cbMap.ComponentBases {
		base := bases.Base{TopologyBaseComponents: autogen_bases.TopologyBaseComponents{
			TopologyBaseName:          cbName,
			OrgID:                     sys.OrgID,
			TopologySystemComponentID: sys.TopologySystemComponentID,
			TopologyClassTypeID:       class_types.BaseClassTypeID,
		}}
		err = tx.QueryRow(ctx, checkExistingClusterBase, cbName, sys.OrgID).Scan(&base.TopologyBaseComponentID)
		if err == pgx.ErrNoRows {
			err = nil
			// Insert the cluster base
			tx, err = create_bases.InsertBaseTx(ctx, &base, tx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
				return tx, err
			}
			err = tx.Commit(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
				return tx, err
			}
			tx, err = apps.Pg.Begin(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to start tx")
				return tx, err
			}
		}
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
			return tx, err
		}

		for sbName, _ := range cb {
			sbEntry := skeletons.SkeletonBase{TopologySkeletonBaseComponents: autogen_bases.TopologySkeletonBaseComponents{
				OrgID:                    sys.OrgID,
				TopologyBaseComponentID:  base.TopologyBaseComponentID,
				TopologyClassTypeID:      class_types.SkeletonBaseClassTypeID,
				TopologySkeletonBaseName: sbName,
			}}

			err = tx.QueryRow(ctx, checkExistingClusterSkeletonBase, sbName, base.TopologyBaseComponentID, sys.OrgID).Scan(&sbEntry.TopologySkeletonBaseID)
			if err == pgx.ErrNoRows && sbEntry.TopologySkeletonBaseID <= 0 {
				err = nil
				tx, err = create_skeletons.InsertSkeletonBaseTx(ctx, &sbEntry, tx)
				if sbEntry.TopologySkeletonBaseID > 0 && err == pgx.ErrNoRows {
					continue
				}
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
					return tx, err
				}
				err = tx.Commit(ctx)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
					return tx, err
				}
				tx, err = apps.Pg.Begin(ctx)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to start tx")
					return tx, err
				}
			}
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("InsertCluster: failed to insert system")
				return tx, err
			}
		}
	}
	return tx, nil
}
