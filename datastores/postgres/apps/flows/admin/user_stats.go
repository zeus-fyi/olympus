package flows_admin

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
)

type UserFlowStats struct {
	Email     string `json:"email"`
	FlowCount int    `json:"flowCount"`
}

func SelectUserFlowStats(ctx context.Context, ou org_users.OrgUser) ([]UserFlowStats, error) {
	q := `	SELECT email, o.org_id, COUNT(*)
			FROM ai_workflow_runs wr 
			JOIN orchestrations o ON o.orchestration_id = wr.orchestration_id
			JOIN org_users ou ON ou.org_id = o.org_id
			JOIN users u ON u.user_id = ou.user_id
			GROUP BY email, o.org_id
			ORDER BY email DESC`

	rows, rerr := apps.Pg.Query(ctx, q)
	if rerr != nil {
		return nil, rerr
	}
	defer rows.Close()

	var configs []UserFlowStats
	for rows.Next() {
		var c UserFlowStats
		err := rows.Scan(&c.Email, &c.FlowCount)
		if err != nil {
			log.Err(err).Msg("SelectAuthedClusterConfigsByOrgID")
			return nil, err
		}

		configs = append(configs, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}
