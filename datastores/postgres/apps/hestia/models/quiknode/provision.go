package hestia_quicknode_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	hestia_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/autogen"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertProvisionedQuickNodeService(ctx context.Context, ps hestia_autogen_bases.ProvisionedQuicknodeServices) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO provisioned_quicknode_services(quicknode_id, endpoint_id, http_url, network, plan, active, org_id, wss_url, chain)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				  ON CONFLICT (quicknode_id) 
				  DO UPDATE SET 
				  updated_at = EXCLUDED.updated_at,
				  endpoint_id = EXCLUDED.endpoint_id,
				  http_url = EXCLUDED.http_url,
				  network = EXCLUDED.network,
				  plan = EXCLUDED.plan,
				  active = EXCLUDED.active,
				  org_id = EXCLUDED.org_id,
				  wss_url = EXCLUDED.wss_url,
				  chain = EXCLUDED.chain;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, ps.QuickNodeID, ps.EndpointID, ps.HttpURL, ps.Network, ps.Plan, ps.Active, ps.OrgID, ps.WssURL, ps.Chain)
	if err == pgx.ErrNoRows {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertProvisionedQuickNodeService"))
}

func UpdateProvisionedQuickNodeService(ctx context.Context, ps hestia_autogen_bases.ProvisionedQuicknodeServices) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET endpoint_id = $1, http_url = $2, network = $3, plan = $4, wss_url = $5, chain = $6
				  WHERE org_id = $1 AND quicknode_id = $2
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ps.OrgID, ps.EndpointID, ps.HttpURL, ps.Network, ps.Plan, ps.WssURL, ps.Chain).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("UpdateProvisionedQuickNodeService"))
}

func DeactivateProvisionedQuickNodeService(ctx context.Context, quickNodeID, endpointID string, ou org_users.OrgUser) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE provisioned_quicknode_services
    			  SET active = false
				  WHERE org_id = $1 AND quicknode_id = $2
			      RETURNING quicknode_id;
				  `
	qnID := ""
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, ou.OrgID, quickNodeID, endpointID).Scan(&qnID)
	if err != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader("DeactivateProvisionedQuickNodeService"))
}
