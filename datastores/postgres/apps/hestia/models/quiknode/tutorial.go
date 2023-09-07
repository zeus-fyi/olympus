package hestia_quicknode_models

import (
	"context"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func ToggleTutorialSetting(ctx context.Context, qnID string) (bool, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `WITH cte_tutorial AS (
					SELECT public_key AS quicknode_id
					FROM orgs o 
					INNER JOIN org_users ou ON ou.org_id = o.org_id
					INNER JOIN users_keys usk ON usk.user_id = ou.user_id
					WHERE o.org_id = $1 AND public_key_name = 'quickNodeMarketPlace' AND public_key_verified = true
	   	          )
				  UPDATE quicknode_marketplace_customer
				  SET tutorial_on = NOT tutorial_on
				  WHERE quicknode_id IN (SELECT quicknode_id FROM cte_tutorial)
				  RETURNING tutorial_on`

	var tutorialOn bool
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, qnID).Scan(&tutorialOn)
	if err != nil {
		return tutorialOn, misc.ReturnIfErr(err, q.LogHeader("InsertProvisionedQuickNodeService"))
	}
	return tutorialOn, misc.ReturnIfErr(err, q.LogHeader("InsertProvisionedQuickNodeService"))
}
