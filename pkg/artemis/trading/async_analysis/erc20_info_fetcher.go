package async_analysis

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

func (c *ContractAnalysis) FindERC20TokenInfo(ctx context.Context) error {
	/*
		1. erc20 decimals, symbol, name
		2. storage slot for balanceOf
		3. transfer tax percentage
	*/
	num := 1
	den := 1
	name, sym := "", ""
	decimals := 0
	tradingEnabled := false
	token := artemis_autogen_bases.Erc20TokenInfo{
		Address:                c.SmartContractAddr,
		ProtocolNetworkID:      1,
		BalanceOfSlotNum:       -1,
		Name:                   &name,
		Symbol:                 &sym,
		Decimals:               &decimals,
		TransferTaxNumerator:   &num,
		TransferTaxDenominator: &den,
		TradingEnabled:         &tradingEnabled,
	}
	err := artemis_validator_service_groups_models.InsertERC20TokenInfo(ctx, token)
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: InsertERC20TokenInfo")
		return err
	}
	return nil
}
