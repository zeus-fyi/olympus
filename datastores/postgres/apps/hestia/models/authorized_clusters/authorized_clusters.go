package authorized_clusters

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	util "github.com/wealdtech/go-eth2-util"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

const Sn = "extClusterConfigs"

type K8sClusterConfig struct {
	ExtConfigStrID string `json:"extConfigStrID"`
	ExtConfigID    int    `json:"extConfigID,omitempty"`

	zeus_common_types.CloudCtxNs `json:"cloudCtxNs"`
	ContextAlias                 string `json:"contextAlias"`
	IsActive                     bool   `json:"isActive,omitempty"`
	IsPublic                     bool   `json:"isPublic,omitempty"`

	Path              filepaths.Path `json:"-"`
	InMemFsKubeConfig memfs.MemFS    `json:"-"`
}

func (ecc *K8sClusterConfig) GetHashedKey(orgID int) string {
	orgStr := fmt.Sprintf("%d", orgID)
	return fmt.Sprintf("%x", util.Keccak256([]byte(orgStr+ecc.CloudProvider+ecc.Region+ecc.Context)))
}

func InsertOrUpdateK8sClusterConfigs(ctx context.Context, ou org_users.OrgUser, configs []K8sClusterConfig) error {
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare the SQL statement for inserting or updating
	stmt := `INSERT INTO public.authorized_cluster_configs (ext_config_id, org_id, cloud_provider, region, context, context_alias, env, is_active)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
             ON CONFLICT ("org_id", "cloud_provider", "context")
             DO UPDATE SET 
                 region = EXCLUDED.region, 
                 context_alias = EXCLUDED.context_alias, 
                 env = EXCLUDED.env,
				 is_active = EXCLUDED.is_active;`

	// Iterate over configs and execute the upsert operation for each
	for i, config := range configs {
		ts := chronos.Chronos{}
		if config.ExtConfigStrID == "" && config.ExtConfigID == 0 {
			configs[i].ExtConfigID = ts.UnixTimeStampNow()
			configs[i].ExtConfigStrID = fmt.Sprintf("%d", configs[i].ExtConfigID)
		}
		_, err = tx.Exec(ctx, stmt, configs[i].ExtConfigID, ou.OrgID, config.CloudProvider, config.Region, config.Context, config.ContextAlias, config.Env, config.IsActive)
		if err != nil {
			log.Err(err).Msg("InsertOrUpdateExtClusterConfigs: failed to insert or update ext cluster config")
			// Roll back the transaction in case of error and return
			return err
		}
	}
	// Commit the transaction
	return tx.Commit(ctx)
}

func InsertOrUpdateExtClusterConfigsUnique(ctx context.Context, ou org_users.OrgUser, configs []K8sClusterConfig) error {
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare the SQL statement for inserting or updating
	stmt := `INSERT INTO public.authorized_cluster_configs (ext_config_id, org_id, cloud_provider, region, context, context_alias, env, is_active)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
             ON CONFLICT (org_id, cloud_provider, context)
             DO NOTHING 
             RETURNING ext_config_id;`
	// Iterate over configs and execute the upsert operation for each
	for i, config := range configs {
		ts := chronos.Chronos{}
		if config.ExtConfigStrID == "" && config.ExtConfigID == 0 {
			configs[i].ExtConfigID = ts.UnixTimeStampNow()
			configs[i].ExtConfigStrID = fmt.Sprintf("%d", configs[i].ExtConfigID)
		}
		_, err = tx.Exec(ctx, stmt, configs[i].ExtConfigID, ou.OrgID, config.CloudProvider, config.Region, config.Context, config.ContextAlias, config.Env, config.IsActive)
		if err != nil {
			log.Err(err).Msg("InsertOrUpdateExtClusterConfigs: failed to insert or update ext cluster config")
			// Roll back the transaction in case of error and return
			return err
		}
	}
	// Commit the transaction
	return tx.Commit(ctx)
}
