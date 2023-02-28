package artemis_validator_service_groups_models

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

const (
	ArtemisScheduledDelivery = "ArtemisScheduledEthServices"
	GweiGasFees              = 100000000
	GweiThirtyTwoEth         = 32000000000
)

func InsertDeliverySchedule(ctx context.Context, sd artemis_autogen_bases.EthScheduledDelivery) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  INSERT INTO eth_scheduled_delivery(public_key, protocol_network_id, delivery_schedule_type, amount, units)
				  VALUES ($1, $2, $3, $4, $5)
				  `
	log.Debug().Interface("InsertDeliverySchedule", q.LogHeader(ArtemisScheduledDelivery))
	_, err := apps.Pg.Exec(ctx, q.RawQuery, sd.PublicKey, sd.ProtocolNetworkID, sd.DeliveryScheduleType, sd.Amount, sd.Units)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery)); returnErr != nil {
		return err
	}
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}
