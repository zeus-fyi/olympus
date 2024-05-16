package flows_admin

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
)

type UserFlowStats struct {
	OrgID     int    `json:"orgID"`
	Email     string `json:"email"`
	FlowCount int    `json:"flowCount"`
}

var skipMap = map[string]bool{
	"1712629107490733000": true,
	"1685378241971196000": true,
	"7138983863666903883": true,
	"support@zeus.fyi":    true,
}

func SelectUserFlowStats(ctx context.Context) ([]UserFlowStats, error) {
	q := `	SELECT email, o.org_id, COUNT(*)
			FROM ai_workflow_runs wr 
			JOIN orchestrations o ON o.orchestration_id = wr.orchestration_id
			JOIN org_users ou ON ou.org_id = o.org_id
			JOIN users u ON u.user_id = ou.user_id
			GROUP BY email, o.org_id
			ORDER BY o.org_id DESC`
	rows, rerr := apps.Pg.Query(ctx, q)
	if rerr != nil {
		return nil, rerr
	}
	defer rows.Close()
	var configs []UserFlowStats
	for rows.Next() {
		var c UserFlowStats
		err := rows.Scan(&c.Email, &c.OrgID, &c.FlowCount)
		if err != nil {
			log.Err(err).Msg("SelectUserFlowStats")
			return nil, err
		}
		if _, ok := skipMap[fmt.Sprintf("%d", c.OrgID)]; !ok {
			if _, ok2 := skipMap[c.Email]; !ok2 {
				configs = append(configs, c)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return configs, nil
}
