package create_topology

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/zeus/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
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
			INSERT INTO topologies_org_cloud_ctx_ns(org_id, cloud_provider, context, region, namespace, namespace_alias)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (cloud_provider, context, region, namespace) DO NOTHING
			RETURNING cloud_ctx_ns_id`
	return q
}

func (c *CreateTopologiesOrgCloudCtxNs) InsertTopologyAccessCloudCtxNs(ctx context.Context) error {
	q := c.GetInsertTopologyOrgCtxQueryParams()
	log.Debug().Interface("InsertTopologyAccessCloudCtxNs:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, c.OrgID, c.CloudProvider, c.Context, c.Region, c.Namespace, c.NamespaceAlias).Scan(&c.CloudCtxNsID)
	if err.Error() == "no rows in result set" {
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}

func (c *CreateTopologiesOrgCloudCtxNs) GetDeleteTopologyOrgCtxQueryParams() sql_query_templates.QueryParams {
	q := sql_query_templates.NewQueryParam("InsertTopologyAccessCloudCtxNs", "topologies_org_cloud_ctx_ns", "where", 1000, []string{})
	q.TableName = c.GetTableName()
	q.Columns = c.GetTableColumns()
	q.Values = []apps.RowValues{c.GetRowValues("default")}

	q.RawQuery = `
			DELETE FROM topologies_org_cloud_ctx_ns
			WHERE org_id = $1 AND cloud_provider = $2 AND context = $3 AND region = $4 AND namespace = $5
			`
	return q
}
func (c *CreateTopologiesOrgCloudCtxNs) DeleteTopologyAccessCloudCtxNs(ctx context.Context) error {
	q := c.GetDeleteTopologyOrgCtxQueryParams()
	log.Debug().Interface("DeleteTopologyAccessCloudCtxNs:", q.LogHeader(Sn))
	_, err := apps.Pg.Exec(ctx, q.RawQuery, c.OrgID, c.CloudProvider, c.Context, c.Region, c.Namespace)
	if err == pgx.ErrNoRows {
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
