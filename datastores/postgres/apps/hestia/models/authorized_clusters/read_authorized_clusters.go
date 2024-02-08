package authorized_clusters

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

func SelectAuthedClusterConfigsByOrgID(ctx context.Context, ou org_users.OrgUser) ([]K8sClusterConfig, error) {
	q := `SELECT ext_config_id, ext_config_id::text, cloud_provider, region, context, context_alias, env, is_active
		FROM public.authorized_cluster_configs
		WHERE org_id = $1;`

	rows, rerr := apps.Pg.Query(ctx, q, ou.OrgID)
	if rerr != nil {
		return nil, rerr
	}
	defer rows.Close()

	var configs []K8sClusterConfig
	for rows.Next() {
		var c K8sClusterConfig
		err := rows.Scan(&c.ExtConfigID, &c.ExtConfigStrID, &c.CloudProvider, &c.Region, &c.Context, &c.ContextAlias, &c.Env, &c.IsActive)
		if err != nil {
			log.Err(err).Msg("SelectAuthedClusterConfigsByOrgID")
			return nil, err
		}
		if c.Context == "" {
			continue
		}
		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}

func SelectAuthedAndPublicClusterConfigsByOrgID(ctx context.Context, ou org_users.OrgUser) ([]K8sClusterConfig, error) {
	q := `SELECT ext_config_id, ext_config_id::text, cloud_provider, region, context, context_alias, env, is_active, is_public
		FROM public.authorized_cluster_configs
		WHERE org_id = $1 OR (is_public = true AND org_id = 7138983863666903883);`

	rows, rerr := apps.Pg.Query(ctx, q, ou.OrgID)
	if rerr != nil {
		return nil, rerr
	}
	defer rows.Close()

	var configs []K8sClusterConfig
	for rows.Next() {
		var c K8sClusterConfig
		err := rows.Scan(&c.ExtConfigID, &c.ExtConfigStrID, &c.CloudProvider, &c.Region, &c.Context, &c.ContextAlias, &c.Env, &c.IsActive, &c.IsPublic)
		if err != nil {
			log.Err(err).Msg("SelectAuthedClusterConfigsByOrgID")
			return nil, err
		}
		if c.Context == "" {

			continue
		}
		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}

func SelectAuthedClusterByRouteAndOrgID(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs) (*K8sClusterConfig, error) {
	q := `SELECT ext_config_id, ext_config_id::text, cloud_provider, region, context, context_alias, env, is_active, is_public
		FROM public.authorized_cluster_configs
		WHERE org_id = $1 AND is_active = true AND ext_config_id = $2 AND cloud_provider = $3 AND region = $4 AND context = $5`

	ccid, err := strconv.Atoi(cloudCtxNs.ClusterCfgStrID)
	if err != nil {
		log.Err(err).Interface("ou", ou).Interface("cloudCtxNs", cloudCtxNs).Msg("SelectAuthedClusterByRouteAndOrgID")
		return nil, err
	}
	var ccfg K8sClusterConfig
	rerr := apps.Pg.QueryRowWArgs(ctx, q, ou.OrgID, ccid, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context).Scan(
		&ccfg.ExtConfigID, &ccfg.ExtConfigStrID, &ccfg.CloudProvider, &ccfg.Region, &ccfg.Context, &ccfg.ContextAlias, &ccfg.Env, &ccfg.IsActive, &ccfg.IsPublic)
	if rerr != nil {
		log.Err(rerr).Interface("ou", ou).Interface("cloudCtxNs", cloudCtxNs).Msg("SelectAuthedClusterByRouteAndOrgID")
		return nil, rerr
	}

	return &ccfg, nil
}

func SelectAuthedClusterByRouteOnlyAndOrgID(ctx context.Context, ou org_users.OrgUser, cloudCtxNs zeus_common_types.CloudCtxNs) (*K8sClusterConfig, error) {
	q := `SELECT ext_config_id, ext_config_id::text, cloud_provider, region, context, context_alias, env, is_active, is_public
		FROM public.authorized_cluster_configs
		WHERE org_id = $1 AND is_active = true AND cloud_provider = $2 AND region = $3 AND context = $4 AND is_public = false`

	var ccfg K8sClusterConfig
	rerr := apps.Pg.QueryRowWArgs(ctx, q, ou.OrgID, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context).Scan(
		&ccfg.ExtConfigID, &ccfg.ExtConfigStrID, &ccfg.CloudProvider, &ccfg.Region, &ccfg.Context, &ccfg.ContextAlias, &ccfg.Env, &ccfg.IsActive, &ccfg.IsPublic)
	if rerr == pgx.ErrNoRows {
		return nil, nil
	}
	if rerr != nil {
		log.Err(rerr).Interface("ou", ou).Interface("cloudCtxNs", cloudCtxNs).Msg("SelectAuthedClusterByRouteAndOrgID")
		return nil, rerr
	}
	return &ccfg, nil
}
