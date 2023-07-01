package artemis_realtime_trading

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
)

/*
  adding in other filters here
	  - filter by token
	  - filter by profit
	  - filter by risk score
	  - adds sourcing of new blocks
*/

var TokenMap map[string]artemis_autogen_bases.Erc20TokenInfo

func InitTokenFilter(ctx context.Context) {
	_, tm, terr := artemis_validator_service_groups_models.SelectERC20Tokens(ctx)
	if terr != nil {
		panic(terr)
	}
	TokenMap = tm
}

func (a *ActiveTrading) EntryTxFilter(ctx context.Context, tx *types.Transaction) error {
	if tx.To() == nil {
		return errors.New("ActiveTrading: EntryTxFilter, tx.To() is nil")
	}
	to := tx.To().String()
	tmTradingEnabled := TokenMap[tx.To().String()].TradingEnabled
	if tmTradingEnabled == nil {
		log.Info().Msgf("ActiveTrading: EntryTxFilter, erc20 at address %s not registered", to)
		chainId := tx.ChainId().Int64()
		err := artemis_validator_service_groups_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           to,
			ProtocolNetworkID: int(chainId),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
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
	tradingEnabled := false
	tmTradingEnabled := TokenMap[to].TradingEnabled
	if tmTradingEnabled == nil {
		return errors.New("ActiveTrading: EntryTxFilter, erc20 at address not registered")
	}
	tradingEnabled = *TokenMap[to].TradingEnabled
	if !tradingEnabled {
		return errors.New("ActiveTrading: EntryTxFilter, trading not enabled for this token")
	}
	return nil
}
