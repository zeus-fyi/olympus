package ext_clusters

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
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

func InsertOrUpdateExtClusterConfigs(ctx context.Context, ou org_users.OrgUser, extClusterConfigs []ExtClusterConfig) error {
	q := sql_query_templates.NewQueryParam("InsertOrUpdateExtClusterConfigs", "nodes", "where", 1000, []string{})
	cte := sql_query_templates.CTE{Name: "InsertOrUpdateExtClusterConfigs"}
	cte.SubCTEs = []sql_query_templates.SubCTE{}
	cte.Params = []interface{}{}
	ts := chronos.Chronos{}
	for _, cc := range extClusterConfigs {
		if cc.ExtConfigStrID == "" {
			cc.ExtConfigID = ts.UnixTimeStampNow()
		} else {
			excID, err := strconv.Atoi(cc.ExtConfigStrID)
			if err != nil {
				log.Err(err).Msg("")
			}
			cc.ExtConfigID = excID
		}
		queryName := fmt.Sprintf("cc_insert_%s", cc.ExtConfigStrID)
		scte := sql_query_templates.NewSubInsertCTE(queryName)
		scte.TableName = "ext_cluster_configs"
		cte.OnConflicts = []string{"org_id", "cloud_provider", "region", "context"}
		cte.OnConflictsUpdateColumns = []string{"context_alias", "env", "region"}
		scte.Columns = []string{"ext_config_id", "org_id", "cloud_provider", "region", "context", "context_alias", "env"}
		pgValues := apps.RowValues{
			cc.ExtConfigID,
			ou.OrgID,
			cc.CloudProvider,
			cc.Region,
			cc.Context,
			cc.ContextAlias,
			cc.Env,
		}
		scte.Values = []apps.RowValues{pgValues}
		cte.SubCTEs = append(cte.SubCTEs, scte)
	}
	q.RawQuery = cte.GenerateChainedCTE()
	r, err := apps.Pg.Exec(ctx, q.RawQuery, cte.Params...)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("Configs: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
