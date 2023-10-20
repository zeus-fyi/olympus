package artemis_validator_service_groups_models

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
	zeus_ecdsa "github.com/zeus-fyi/zeus/pkg/aegis/crypto/ecdsa"
	artemis_client "github.com/zeus-fyi/zeus/pkg/artemis/client"
	artemis_req_types "github.com/zeus-fyi/zeus/pkg/artemis/client/req_types"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	resty_base "github.com/zeus-fyi/zeus/zeus/z_client/base"
)

const (
	ArtemisScheduledDelivery = "ArtemisScheduledEthServices"
	GweiGasFees              = 100000000
	GweiThirtyTwoEth         = 32000000000
)

var ArtemisClient = artemis_client.ArtemisClient{
	Account:        zeus_ecdsa.Account{},
	Resty:          resty_base.Resty{},
	ArtemisConfigs: []*artemis_client.ArtemisConfig{&artemis_client.ArtemisEthereumEphemeral},
}

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

	addr := accounts.HexToAddress(sd.PublicKey)
	rr := artemis_req_types.SendEtherPayload{
		TransferArgs: artemis_req_types.TransferArgs{
			Amount:    big.NewInt(1).Mul(signing_automation_ethereum.Gwei, big.NewInt(int64(sd.Amount))),
			ToAddress: addr,
		},
	}
	rx, err := ArtemisClient.SendEther(ctx, rr, artemis_client.ArtemisEthereumEphemeral)
	if err != nil {
		log.Ctx(ctx).Err(err).Msgf("ArtemisClient.SendEther: %s", q.LogHeader(ArtemisScheduledDelivery))
	}
	log.Info().Interface("rx", rx).Msgf("ArtemisClient.SendEther: %s", q.LogHeader(ArtemisScheduledDelivery))
	return misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}

func SelectEphemeryDeliverySchedule(ctx context.Context) ([]artemis_autogen_bases.EthScheduledDelivery, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `
				  SELECT public_key, amount, units
				  FROM eth_scheduled_delivery
				  WHERE protocol_network_id = $1
				  `
	sdeliveries := []artemis_autogen_bases.EthScheduledDelivery{}
	log.Debug().Interface("SelectDeliverySchedule", q.LogHeader(ArtemisScheduledDelivery))
	rows, err := apps.Pg.Query(ctx, q.RawQuery, hestia_req_types.EthereumEphemeryProtocolNetworkID)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery)); returnErr != nil {
		return sdeliveries, err
	}
	defer rows.Close()
	for rows.Next() {
		sd := artemis_autogen_bases.EthScheduledDelivery{}
		rowErr := rows.Scan(
			&sd.PublicKey, &sd.Amount, &sd.Units,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader(ModelName))
			return nil, rowErr
		}
		sdeliveries = append(sdeliveries, sd)
	}
	return sdeliveries, misc.ReturnIfErr(err, q.LogHeader(ArtemisScheduledDelivery))
}
