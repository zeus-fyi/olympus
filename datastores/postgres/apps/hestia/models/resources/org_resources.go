package hestia_compute_resources

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func AddResourcesToOrg(ctx context.Context, orgID, resourceID int, quantity float64) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO org_resources(org_id, resource_id, quantity)
				  VALUES ($1, $2, $3)
				  `
	_, err := apps.Pg.Exec(ctx, q.RawQuery, orgID, resourceID, quantity)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(Sn)); returnErr != nil {
		return returnErr
	}
	return err
}
