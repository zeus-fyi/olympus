package artemis_validator_service_groups_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertERC20TokenInfo(ctx context.Context, tx artemis_autogen_bases.Erc20TokenInfo) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `INSERT INTO erc20_token_info(address, protocol_network_id, balance_of_slot_num, name, symbol, decimals, transfer_tax_numerator, transfer_tax_denominator, trading_enabled)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				  ON CONFLICT (address) 
				  DO UPDATE SET 
				  protocol_network_id = excluded.protocol_network_id,
				  balance_of_slot_num = excluded.balance_of_slot_num,
				  name = excluded.name,
				  symbol = excluded.symbol,
				  decimals = excluded.decimals,
				  transfer_tax_numerator = excluded.transfer_tax_numerator,
				  transfer_tax_denominator = excluded.transfer_tax_denominator,
				  trading_enabled = excluded.trading_enabled;`

	protocolIDNum := 1
	if tx.ProtocolNetworkID != 0 {
		protocolIDNum = tx.ProtocolNetworkID
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, tx.Address, protocolIDNum, tx.BalanceOfSlotNum)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertERC20TokenInfo"))
}

func SelectERC20TokenInfo(ctx context.Context, tx artemis_autogen_bases.Erc20TokenInfo) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT balance_of_slot_num
				  FROM erc20_token_info
				  WHERE address = $1 AND protocol_network_id = $2;`
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, tx.Address, tx.ProtocolNetworkID).Scan(&tx.BalanceOfSlotNum)
	if err == pgx.ErrNoRows {
		return -1, nil
	}
	return tx.BalanceOfSlotNum, misc.ReturnIfErr(err, q.LogHeader("SelectERC20TokenInfo"))
}
