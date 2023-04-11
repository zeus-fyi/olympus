package hestia_compute_resources

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func AddResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, freeTrial bool) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
				  VALUES ($1, $2, $3, $4)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func AddDigitalOceanNodePoolResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64, nodePoolID, nodeContextID string, freeTrial bool) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = ` WITH cte_org_resources AS (
					  INSERT INTO org_resources(org_id, resource_id, quantity, free_trial)
					  VALUES ($1, $2, $3, $6)
					  RETURNING org_resources_id
				  ) INSERT INTO digitalocean_node_pools(org_resources_id, resource_id, node_pool_id, node_context_id)
					VALUES ((SELECT org_resources_id FROM cte_org_resources), $2, $4, $5)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity, nodePoolID, nodeContextID, freeTrial)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func RemoveFreeTrialOrgResources(ctx context.Context, orgID int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `DELETE FROM org_resources
				  WHERE org_id = $1	AND free_trial = true
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}

func UpdateFreeTrialOrgResourcesToPaid(ctx context.Context, orgID int) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE org_resources
				  SET free_trial = false
				  WHERE org_id = $1	AND free_trial = true
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID)
	if err == pgx.ErrNoRows {
		return nil
	}
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}
