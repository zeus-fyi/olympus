package create_clusters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/conversions/chart_workload"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/bases/skeletons"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/class_types"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/topology/classes/systems"
	create_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases"
	create_infra "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/infra"
	create_skeletons "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/definitions/bases/skeleton"
	create_systems "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/create/topologies/topology/classes/systems"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	zeus_templates "github.com/zeus-fyi/olympus/zeus/api/v1/zeus/topology/infra/create/templates"
)

const Sn = "Cluster"

var (
	checkExistingCluster             = `SELECT topology_system_component_id FROM topology_system_components WHERE topology_system_component_name = $1 AND org_id = $2`
	checkExistingClusterBase         = `SELECT topology_base_component_id FROM topology_base_components WHERE topology_base_name = $1 AND org_id = $2`
	checkExistingClusterSkeletonBase = `SELECT topology_skeleton_base_id FROM topology_skeleton_base_components WHERE topology_skeleton_base_name = $1 AND topology_base_component_id = $2 AND org_id = $3`
)

func InsertCluster(ctx context.Context, tx pgx.Tx, sys *systems.Systems, cbMap zeus_templates.ClusterPreviewWorkloads, ou org_users.OrgUser) (pgx.Tx, error) {
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

		for sbName, skeleton := range cb {
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
			if skeleton.Service == nil && skeleton.ServiceMonitor == nil && skeleton.ConfigMap == nil && skeleton.Deployment == nil && skeleton.StatefulSet == nil {
				continue
			}
			nk := chart_workload.TopologyBaseInfraWorkload{}
			if skeleton.Deployment != nil {
				b, berr := json.Marshal(skeleton.Deployment)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			if skeleton.StatefulSet != nil {
				b, berr := json.Marshal(skeleton.StatefulSet)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			if nk.StatefulSet != nil && nk.Deployment != nil {
				err = errors.New("cannot include both a stateful set and deployment, must only choose one per topology infra chart components")
				return nil, err
			}

			if skeleton.Service != nil {
				b, berr := json.Marshal(skeleton.Service)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			if skeleton.Ingress != nil {
				b, berr := json.Marshal(skeleton.Ingress)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			if skeleton.ConfigMap != nil {
				b, berr := json.Marshal(skeleton.ConfigMap)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			if skeleton.ServiceMonitor != nil {
				b, berr := json.Marshal(skeleton.ServiceMonitor)
				if berr != nil {
					log.Err(berr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, berr
				}
				err = nk.DecodeBytes(b)
				if err != nil {
					log.Err(err).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: TopologyCreateRequestFromUI, CreateChartWorkloadFromTopologyBaseInfraWorkload")
					return nil, err
				}
			}

			cw, cerr := nk.CreateChartWorkloadFromTopologyBaseInfraWorkload()
			if cerr != nil {
				log.Err(cerr).Interface("kubernetesWorkload", nk).Msg("TopologyActionCreateRequest: CreateTopology, CreateChartWorkloadFromTopologyBaseInfraWorkload")
				return nil, cerr
			}
			inf := create_infra.NewCreateInfrastructure()
			inf.ChartWorkload = cw
			inf.ClusterClassName = cbMap.ClusterName
			inf.ComponentBaseName = cbName
			inf.SkeletonBaseName = sbName

			inf.OrgID = ou.OrgID
			inf.UserID = ou.UserID
			inf.Name = sbName
			inf.Chart.ChartName = sbName
			ts := chronos.Chronos{}
			inf.ChartVersion = fmt.Sprintf("%d", ts.UnixTimeStampNow())
			inf.Tag = "latest"
			tx, err = inf.InsertInfraBaseTx(ctx, tx)
			if err != nil {
				pgErr := err.(*pgconn.PgError)
				switch {
				case strings.Contains(pgErr.Error(), "chart_package_unique"):
					err = errors.New("chart name and version already exists")
					return nil, err
				default:
					log.Err(err).Msg("CreateTopologyFromUI: CreateTopology, InsertInfraBase")
					err = errors.New("unable to add chart, verify it is a valid kubernetes workload that's supported")
				}
				log.Err(err).Interface("orgUser", ou).Msg("TopologyActionCreateRequest: CreateTopology, InsertInfraBase")
				return nil, err
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
	}
	return tx, nil
}
