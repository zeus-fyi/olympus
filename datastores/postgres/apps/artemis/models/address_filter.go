package artemis_validator_service_groups_models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func SelectSourceAddresses(ctx context.Context, protocolID int) ([]string, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT address
				  FROM eth_mev_address_filter
				  WHERE protocol_network_id = $1
				  `
	log.Debug().Interface("SelectSourceAddresses", q.LogHeader("SelectSourceAddresses"))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, protocolID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectSourceAddresses")); returnErr != nil {
		return nil, err
	}
	var addrbook []string
	defer rows.Close()
	for rows.Next() {
		var addr string
		rowErr := rows.Scan(
			&addr,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectSourceAddresses"))
			return nil, rowErr
		}
		addrbook = append(addrbook, addr)
	}
	return addrbook, misc.ReturnIfErr(err, q.LogHeader("SelectSourceAddresses"))
}
