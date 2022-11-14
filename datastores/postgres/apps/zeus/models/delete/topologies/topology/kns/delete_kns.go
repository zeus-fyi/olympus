package delete_kns

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/topologies/definitions/kns"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const Sn = "DeleteTopologyKubeCtxNs"

func getDeleteKnsQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteKns"

	q.RawQuery = `
	WITH cte_delete_kns AS (
		DELETE FROM topologies_kns 
		WHERE topology_id = $1 AND cloud_provider = $2 AND region = $3 AND context = $4 AND namespace = $5 AND env = $6
	) SELECT true
	`
	return q
}

func DeleteKns(ctx context.Context, kns *kns.TopologyKubeCtxNs) error {
	q := getDeleteKnsQuery()
	log.Debug().Interface("DeleteQuery:", q.LogHeader(Sn))
	success := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, kns.TopologyID, kns.CloudProvider, kns.Region, kns.Context, kns.Namespace, kns.Env).Scan(&success)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
