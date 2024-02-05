package ext_clusters

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

const Sn = "extClusterConfigs"

type ExtClusterConfig struct {
	ExtConfigStrID string `json:"extConfigStrID"`
	ExtConfigID    int    `json:"extConfigID,omitempty"`
	CloudProvider  string `json:"cloudProvider"`
	Region         string `json:"region"`
	Context        string `json:"context"`
	ContextAlias   string `json:"contextAlias"`
	Env            string `json:"env,omitempty"`
}

func InsertOrUpdateExtClusterConfigs(ctx context.Context, ou org_users.OrgUser, configs []ExtClusterConfig) error {
	// Start a transaction
	tx, err := apps.Pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Prepare the SQL statement for inserting or updating
	stmt := `INSERT INTO public.ext_cluster_configs (ext_config_id, org_id, cloud_provider, region, context, context_alias, env)
             VALUES ($1, $2, $3, $4, $5, $6, $7)
             ON CONFLICT (ext_config_id)
             DO UPDATE SET 
                 region = EXCLUDED.region, 
                 context_alias = EXCLUDED.context_alias, 
                 env = EXCLUDED.env
             WHERE org_id = EXCLUDED.org_id;`

	// Iterate over configs and execute the upsert operation for each
	for i, config := range configs {
		ts := chronos.Chronos{}
		if config.ExtConfigStrID == "" && config.ExtConfigID == 0 {
			configs[i].ExtConfigID = ts.UnixTimeStampNow()
			configs[i].ExtConfigStrID = fmt.Sprintf("%d", configs[i].ExtConfigID)
		}
		_, err = tx.Exec(ctx, stmt, configs[i].ExtConfigID, ou.OrgID, config.CloudProvider, config.Region, config.Context, config.ContextAlias, config.Env)
		if err != nil {
			log.Err(err).Msg("InsertOrUpdateExtClusterConfigs: failed to insert or update ext cluster config")
			// Roll back the transaction in case of error and return
			return err
		}
	}
	// Commit the transaction
	return tx.Commit(ctx)
}
