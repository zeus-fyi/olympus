package hestia_compute_resources

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	do_types "github.com/zeus-fyi/olympus/pkg/hestia/digitalocean/types"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func AddResourcesToOrgAndCtx(ctx context.Context, orgID, resourceID int, quantity float64, freeTrial bool, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $4)
				  	  RETURNING org_resource_id
				  ), cte_get_cloud_ctx AS (
					SELECT cloud_ctx_ns_id 
					FROM topologies_org_cloud_ctx_ns
					WHERE cloud_provider = $5 AND region = $6 AND context = $7 AND namespace = $8 AND org_id = $1
					LIMIT 1
				  ) INSERT INTO org_resources_cloud_ctx(org_resource_id, cloud_ctx_ns_id)
					VALUES ((SELECT org_resource_id FROM cte_org_resources), (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx))
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, freeTrial, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context, cloudCtxNs.Namespace)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func AddGkeNodePoolResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, nodePoolID, nodeContextID string, freeTrial bool) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = ` WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $6)
					  RETURNING org_resource_id
				  ) INSERT INTO gke_node_pools(org_resource_id, resource_id, node_pool_id, node_context_id)
					VALUES ((SELECT org_resource_id FROM cte_org_resources), $2, $4, $5)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, nodePoolID, nodeContextID, freeTrial)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func AddEksNodePoolResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, nodePoolID, nodeContextID string, freeTrial bool, clusterCfgStrID string) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = ` WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $6)
					  RETURNING org_resource_id
				  ) INSERT INTO eks_node_pools(org_resource_id, resource_id, node_pool_id, node_context_id, ext_config_id)
					VALUES ((SELECT org_resource_id FROM cte_org_resources), $2, $4, $5, $7)
				  `

	extConfigID := sql.NullInt64{Valid: false}
	if len(clusterCfgStrID) > 0 {
		cid, err := strconv.Atoi(clusterCfgStrID)
		if err != nil {
			log.Err(err).Msg("AddEksNodePoolResourcesToOrg: strconv.Atoi(cctx.ClusterCfg) error")
			return err
		}
		if cid > 0 {
			extConfigID = sql.NullInt64{Int64: int64(cid), Valid: true}
		}
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, nodePoolID, nodeContextID, freeTrial, extConfigID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func AddOvhNodePoolResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, nodePoolID, nodeContextID string, freeTrial bool) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = ` WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $6)
					  RETURNING org_resource_id
				  ) INSERT INTO ovh_node_pools(org_resource_id, resource_id, node_pool_id, node_context_id)
					VALUES ((SELECT org_resource_id FROM cte_org_resources), $2, $4, $5)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, nodePoolID, nodeContextID, freeTrial)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func AddDigitalOceanNodePoolResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, nodePoolID, nodeContextID string, freeTrial bool) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = ` WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $6)
					  RETURNING org_resource_id
				  ) INSERT INTO digitalocean_node_pools(org_resource_id, resource_id, node_pool_id, node_context_id)
					VALUES ((SELECT org_resource_id FROM cte_org_resources), $2, $4, $5)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, nodePoolID, nodeContextID, freeTrial)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func SelectFreeTrialDigitalOceanNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM digitalocean_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND free_trial = true AND begin_service <= CURRENT_TIMESTAMP - INTERVAL '1 hour'
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func GkeSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM gke_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND free_trial = true AND begin_service <= CURRENT_TIMESTAMP - INTERVAL '1 hour'
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func EksSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM eks_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND free_trial = true AND begin_service <= CURRENT_TIMESTAMP - INTERVAL '1 hour'
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func OvhSelectFreeTrialNodes(ctx context.Context, orgID int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM ovh_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND free_trial = true AND begin_service <= CURRENT_TIMESTAMP - INTERVAL '1 hour'
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func RemoveFreeTrialOrgResources(ctx context.Context, orgID int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_org_free_trial_resources AS (
					SELECT org_resource_id FROM org_resources WHERE org_id = $1 AND free_trial = true
			      ), cte_digitalocean_node_pools AS (
					DELETE FROM digitalocean_node_pools
					WHERE org_resource_id IN (SELECT org_resource_id FROM cte_org_free_trial_resources)
				  ), cte_gke_node_pools AS (
					DELETE FROM gke_node_pools
					WHERE org_resource_id IN (SELECT org_resource_id FROM cte_org_free_trial_resources)
				  ), cte_ovh_node_pools AS (
					DELETE FROM ovh_node_pools
					WHERE org_resource_id IN (SELECT org_resource_id FROM cte_org_free_trial_resources)
				  ), cte_eks_node_pools AS (
					DELETE FROM eks_node_pools
					WHERE org_resource_id IN (SELECT org_resource_id FROM cte_org_free_trial_resources)
				  ), cte_org_resource_ctx_id_delete AS (
  					DELETE FROM org_resources_cloud_ctx WHERE org_resource_id IN (SELECT org_resource_id FROM cte_org_free_trial_resources)	
				  )
				  DELETE FROM org_resources
				  WHERE org_id = $1	AND free_trial = true
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func UpdateFreeTrialOrgResourcesToPaid(ctx context.Context, orgID int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE org_resources
				  SET free_trial = false
				  WHERE org_id = $1	AND free_trial = true
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func DoesOrgHaveOngoingFreeTrial(ctx context.Context, orgID int) (bool, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT COUNT(*) > 0 FROM org_resources WHERE org_id = $1 AND free_trial = true AND end_service IS NULL AND quantity > 0
				  `
	var isFreeTrialOngoing bool
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID).Scan(&isFreeTrialOngoing)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return true, returnErr
	}
	return isFreeTrialOngoing, err
}

func GkeSelectNodeResources(ctx context.Context, orgID int, orgResourceIDs []int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM gke_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND org_resource_id = ANY($2::bigint[]) AND free_trial = false
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, pq.Array(orgResourceIDs))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func EksSelectNodeResources(ctx context.Context, orgID int, orgResourceIDs []int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM eks_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND org_resource_id = ANY($2::bigint[]) AND free_trial = false
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, pq.Array(orgResourceIDs))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func OvhSelectNodeResources(ctx context.Context, orgID int, orgResourceIDs []int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM ovh_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND org_resource_id = ANY($2::bigint[]) AND free_trial = false
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, pq.Array(orgResourceIDs))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func SelectNodeResources(ctx context.Context, orgID int, orgResourceIDs []int) ([]do_types.DigitalOceanNodePoolRequestStatus, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT node_pool_id, node_context_id
 				  FROM digitalocean_node_pools
 				  JOIN org_resources USING (org_resource_id)
				  WHERE org_id = $1 AND org_resource_id = ANY($2::bigint[]) AND free_trial = false
				  `
	rows, err := apps.Pg.Query(ctx, q.RawQuery, orgID, pq.Array(orgResourceIDs))
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()
	var nodePools []do_types.DigitalOceanNodePoolRequestStatus
	for rows.Next() {
		np := do_types.DigitalOceanNodePoolRequestStatus{}
		err = rows.Scan(&np.NodePoolID, &np.ClusterID)
		if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
			return nil, returnErr
		}
		nodePools = append(nodePools, np)
	}
	return nodePools, err
}

func UpdateEndServiceOrgResources(ctx context.Context, orgID int, orgResourceIDs []int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE org_resources
				  SET end_service = NOW()
				  WHERE org_id = $1	AND org_resource_id = ANY($2::bigint[]) AND end_service IS NULL
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, pq.Array(orgResourceIDs))
	if err == pgx.ErrNoRows {
		log.Ctx(ctx).Info().Msg("No org resources to update end service")
		return nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}
