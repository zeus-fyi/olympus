package create_topology

import (
	"context"

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
			INSERT INTO topologies_org_cloud_ctx_ns(org_id, cloud_provider, context, region, namespace)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (cloud_provider, context, region, namespace) DO NOTHING
			RETURNING cloud_ctx_ns_id`
	return q
}

func (c *CreateTopologiesOrgCloudCtxNs) InsertTopologyAccessCloudCtxNs(ctx context.Context) error {
	q := c.GetInsertTopologyOrgCtxQueryParams()
	log.Debug().Interface("InsertTopologyAccessCloudCtxNs:", q.LogHeader(Sn))
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, c.OrgID, c.CloudProvider, c.Context, c.Region, c.Namespace).Scan(&c.CloudCtxNsID)
	if err.Error() == "no rows in result set" {
		return nil
	}
	return misc.ReturnIfErr(err, q.LogHeader(Sn))
}
