package artemis_validator_service_groups_models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils/sql_query_templates"
)

func InsertERC20TokenInfo(ctx context.Context, token artemis_autogen_bases.Erc20TokenInfo) error {
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
	if token.ProtocolNetworkID != 0 {
		protocolIDNum = token.ProtocolNetworkID
	}
	_, err := apps.Pg.Exec(ctx, q.RawQuery, token.Address, protocolIDNum, token.BalanceOfSlotNum, token.Name, token.Symbol, token.Decimals, token.TransferTaxNumerator, token.TransferTaxDenominator, token.TradingEnabled)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("InsertERC20TokenInfo"))
}

func UpdateERC20TokenInfo(ctx context.Context, token artemis_autogen_bases.Erc20TokenInfo) error {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `UPDATE erc20_token_info
    			  SET name = $2, symbol = $3, decimals = $4
				  WHERE address = $1;`

	_, err := apps.Pg.Exec(ctx, q.RawQuery, token.Address, token.Name, token.Symbol, token.Decimals)
	if err == pgx.ErrNoRows {
		err = nil
	}
	return misc.ReturnIfErr(err, q.LogHeader("UpdateERC20TokenInfo"))
}

func SelectERC20TokenInfo(ctx context.Context, token artemis_autogen_bases.Erc20TokenInfo) (int, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT balance_of_slot_num
				  FROM erc20_token_info
				  WHERE address = $1 AND protocol_network_id = $2;`
	err := apps.Pg.QueryRowWArgs(ctx, q.RawQuery, token.Address, token.ProtocolNetworkID).Scan(&token.BalanceOfSlotNum)
	if err == pgx.ErrNoRows {
		return -1, nil
	}
	return token.BalanceOfSlotNum, misc.ReturnIfErr(err, q.LogHeader("SelectERC20TokenInfo"))
}

func SelectERC20TokensWithoutMetadata(ctx context.Context) ([]artemis_autogen_bases.Erc20TokenInfo, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT address
				  FROM erc20_token_info
				  WHERE protocol_network_id = $1 AND decimals IS NULL;`

	var tokens []artemis_autogen_bases.Erc20TokenInfo
	rows, err := apps.Pg.Query(ctx, q.RawQuery, 1)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectERC20Tokens")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var token artemis_autogen_bases.Erc20TokenInfo
		rowErr := rows.Scan(
			&token.Address,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectERC20Tokens"))
			return nil, rowErr
		}
		tokens = append(tokens, token)
	}
	return tokens, misc.ReturnIfErr(err, q.LogHeader("SelectERC20Tokens"))
}

func SelectERC20Tokens(ctx context.Context) ([]artemis_autogen_bases.Erc20TokenInfo, error) {
	q := sql_query_templates.QueryParams{}
	q.RawQuery = `SELECT address
				  FROM erc20_token_info
				  WHERE protocol_network_id = $1;`

	var tokens []artemis_autogen_bases.Erc20TokenInfo
	rows, err := apps.Pg.Query(ctx, q.RawQuery, 1)
	if returnErr := misc.ReturnIfErr(err, q.LogHeader("SelectERC20Tokens")); returnErr != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var token artemis_autogen_bases.Erc20TokenInfo
		rowErr := rows.Scan(
			&token.Address,
		)
		if rowErr != nil {
			log.Err(rowErr).Msg(q.LogHeader("SelectERC20Tokens"))
			return nil, rowErr
		}
		tokens = append(tokens, token)
	}
	return tokens, misc.ReturnIfErr(err, q.LogHeader("SelectERC20Tokens"))
}
