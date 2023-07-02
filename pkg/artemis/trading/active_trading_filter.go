package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_trading_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/cache"
)

/*
  adding in other filters here
	  - filter by token
	  - filter by profit
	  - filter by risk score
	  - adds sourcing of new blocks
*/

func (a *ActiveTrading) EntryTxFilter(ctx context.Context, tx *types.Transaction) error {
	if tx.To() == nil {
		return errors.New("ActiveTrading: EntryTxFilter, tx.To() is nil")
	}
	to := tx.To().String()
	tmTradingEnabled := artemis_trading_cache.TokenMap[tx.To().String()].TradingEnabled
	if tmTradingEnabled == nil {
		tradeEnabled := false
		log.Info().Msgf("ActiveTrading: EntryTxFilter, erc20 at address %s not registered", to)
		chainId := tx.ChainId().Int64()
		err := artemis_validator_service_groups_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           to,
			ProtocolNetworkID: int(chainId),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
			TradingEnabled:    &tradeEnabled,
		})
		if err != nil {
			log.Err(err).Msg("ActiveTrading: EntryTxFilter, InsertERC20TokenInfo")
			return err
		}
		return errors.New("ActiveTrading: EntryTxFilter, erc20 at address not registered")
	}

	return nil
}

func (a *ActiveTrading) SimTxFilter(ctx context.Context, tx *types.Transaction) error {
	to := tx.To().String()
	/*
		tmTradingEnabled := TokenMap[to].TradingEnabled
		if tmTradingEnabled == nil {
			return errors.New("ActiveTrading: EntryTxFilter, erc20 at address not registered")
		}
		tradingEnabled := false
		tradingEnabled = *TokenMap[to].TradingEnabled
		if !tradingEnabled {
			return errors.New("ActiveTrading: EntryTxFilter, trading not enabled for this token")
		}
	*/
	if artemis_trading_cache.TokenMap[to].BalanceOfSlotNum < 0 {
		return errors.New("ActiveTrading: EntryTxFilter, balanceOf not cracked yet")
	}
	num := artemis_trading_cache.TokenMap[to].TransferTaxNumerator
	den := artemis_trading_cache.TokenMap[to].TransferTaxDenominator
	if num == nil || den == nil {
		return errors.New("ActiveTrading: EntryTxFilter, transfer tax not set")
	}
	if *num == 0 || *den == 0 {
		return errors.New("ActiveTrading: EntryTxFilter, transfer tax not set")
	}
	return nil
}
