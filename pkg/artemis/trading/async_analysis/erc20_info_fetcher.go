package async_analysis

import (
	"context"

	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

// FindERC20TokenMetadataInfo finds the erc20 token decimals, symbol, name and updates the database
func (c *ContractAnalysis) FindERC20TokenMetadataInfo(ctx context.Context) error {
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
	decimals, err := c.u.Web3Client.GetContractDecimals(ctx, c.SmartContractAddr)
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: GetContractDecimals")
		return err
	}
	token.Decimals = &decimals
	name, err = c.u.Web3Client.GetContractName(ctx, c.SmartContractAddr)
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: GetContractName")
		return err
	}
	token.Name = &name
	sym, err = c.u.Web3Client.GetContractSymbol(ctx, c.SmartContractAddr)
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: GetContractSymbol")
		return err
	}
	token.Symbol = &sym
	err = artemis_validator_service_groups_models.UpdateERC20TokenInfo(ctx, token)
	if err != nil {
		log.Err(err).Msg("ContractAnalysis: InsertERC20TokenInfo")
		return err
	}
	return nil
}
