package create_kns

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_req_types"
)

const Sn = "TopologyKubeCtxNs"

func (k *Kns) InsertKns(ctx context.Context, q sql_query_templates.QueryParams) error {
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	r, err := apps.Pg.Exec(ctx, q.InsertSingleElementQuery())
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	rowsAffected := r.RowsAffected()
	log.Debug().Msgf("TopologyKubeCtxNs: %s, Rows Affected: %d", q.LogHeader(Sn), rowsAffected)
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

/*
func getDeleteKnsQueryByOrgMatchAndCloudCtx() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "DeleteKnsUsingLookup"

	q.RawQuery = `
	WITH cte_get_authed_ctx AS (
		SELECT org_id, cloud_provider, region, context, namespace
		FROM topologies_org_cloud_ctx_ns
		WHERE cloud_provider = $1 AND region = $2 AND context = $3 AND namespace = $4
		LIMIT 1
	), cte_delete_kns AS (
		DELETE FROM topologies_kns
		WHERE cloud_provider = (SELECT cloud_provider FROM cte_get_authed_ctx) AND region = (SELECT region FROM cte_get_authed_ctx) AND namespace = (SELECT namespace FROM cte_get_authed_ctx) AND env = $5
	) SELECT true
	`
	return q
}
*/

// TODO add when no topology id provided
func getInsertKnsQuery() sql_query_templates.QueryParams {
	q := sql_query_templates.QueryParams{}
	q.QueryName = "InsertKns"

	q.RawQuery = `
	WITH cte_insert_kns AS (
		INSERT INTO topologies_kns(topology_id, cloud_provider, region, context, namespace, env)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ON CONSTRAINT kns_pk DO NOTHING
	) SELECT true
	`
	return q
}

func InsertKns(ctx context.Context, kns *zeus_req_types.TopologyDeployRequest) error {
	if kns == nil {
		return nil
	}
	q := getInsertKnsQuery()
	log.Debug().Interface("InsertQuery:", q.LogHeader(Sn))
	success := false
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, kns.TopologyID, kns.CloudProvider, kns.Region, kns.Context, kns.Namespace, kns.Env).Scan(&success)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
