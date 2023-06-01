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
	q.RawQuery = `INSERT INTO erc20_token_info(address, protocol_network_id, balance_of_slot_num)
				  VALUES ($1, $2, $3)
				  ON CONFLICT (address) DO NOTHING;`

	protocolIDNum := 1
	if tx.ProtocolNetworkID != 0 {
		protocolIDNum = tx.ProtocolNetworkID
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, tx.Address, protocolIDNum, tx.BalanceOfSlotNum)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertMempoolTx"))
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
	return tx.BalanceOfSlotNum, misc.ReturnIfErr(err, q.LogHeader("SelectMempoolTx"))
}
