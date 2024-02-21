package create_topology

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type CreateTopologiesOrgCloudCtxNs struct {
	autogen_bases.TopologiesOrgCloudCtxNs
}

func NewCreateTopologiesOrgCloudCtxNs(orgID int, kns zeus_common_types.CloudCtxNs) CreateTopologiesOrgCloudCtxNs {
	ct := CreateTopologiesOrgCloudCtxNs{autogen_bases.TopologiesOrgCloudCtxNs{
		OrgID:         orgID,
		CloudProvider: kns.CloudProvider,
		Context:       kns.Context,
		Region:        kns.Region,
		Namespace:     kns.Namespace,
	}}
	return ct
}

func (c *CreateTopologiesOrgCloudCtxNs) GetInsertTopologyOrgCtxQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("InsertTopologyAccessCloudCtxNs", "topologies_org_cloud_ctx_ns", "where", 1000, []string{})
	q.TableName = c.GetTableName()
	q.Columns = c.GetTableColumns()
	q.Values = []apps.RowValues{c.GetRowValues("default")}

	q.RawQuery = `
			INSERT INTO topologies_org_cloud_ctx_ns(org_id, cloud_provider, region, context, namespace, namespace_alias)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (cloud_provider, context, region, namespace) DO NOTHING
			RETURNING cloud_ctx_ns_id`
	return q
}

func (c *CreateTopologiesOrgCloudCtxNs) InsertTopologyAccessCloudCtxNs(ctx context.Context, orgID int, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	q := c.GetInsertTopologyOrgCtxQueryParams()
	log.Debug().Interface("InsertTopologyAccessCloudCtxNs:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, orgID, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context, cloudCtxNs.Namespace, cloudCtxNs.Alias).Scan(&c.CloudCtxNsID)
	if err == pgx.ErrNoRows {
		log.Info().Msg("InsertTopologyAccessCloudCtxNs: no rows to insert")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func GetDeleteTopologyOrgCtxQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("DeleteTopologyAccessCloudCtxNs", "topologies_org_cloud_ctx_ns", "where", 1000, []string{})

	q.RawQuery = `
		WITH cte_get_cloud_ctx AS (
			SELECT cloud_ctx_ns_id
			FROM topologies_org_cloud_ctx_ns
			WHERE org_id = $1 AND cloud_provider = $2 AND region = $3 AND context = $4 AND namespace = $5
			LIMIT 1
		), cte_remove_cloud_ctx_res AS (
			DELETE
			FROM org_resources_cloud_ctx
			WHERE cloud_ctx_ns_id IN (SELECT cloud_ctx_ns_id FROM cte_get_cloud_ctx)
		)
		DELETE
		FROM topologies_org_cloud_ctx_ns
		WHERE org_id = $1 AND cloud_provider = $2 AND region = $3 AND context = $4 AND namespace = $5;`
	return q
}
func DeleteTopologyAccessCloudCtxNs(ctx context.Context, orgID int, cloudCtxNs zeus_common_types.CloudCtxNs) error {
	q := GetDeleteTopologyOrgCtxQueryParams()
	log.Info().Interface("cloudCtxNs", cloudCtxNs).Interface("orgID:", orgID).Interface("DeleteTopologyAccessCloudCtxNs:", q.LogHeader(Sn))
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, cloudCtxNs.CloudProvider, cloudCtxNs.Region, cloudCtxNs.Context, cloudCtxNs.Namespace)
	if err == pgx.ErrNoRows {
		log.Ctx(ctx).Info().Msg("DeleteTopologyAccessCloudCtxNs: no rows to delete")
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
