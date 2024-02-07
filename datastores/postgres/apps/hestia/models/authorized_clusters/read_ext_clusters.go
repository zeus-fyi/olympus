package authorized_clusters

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
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
			return nil, err
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
			return nil, err
		}
		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}
