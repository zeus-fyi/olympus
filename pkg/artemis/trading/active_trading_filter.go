package artemis_realtime_trading

import (
	"context"

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

func (a *ActiveTrading) FilterTx(ctx context.Context, tx *types.Transaction) *types.Transaction {
	if tx.To() == nil {
		return nil
	}
	to := tx.To().String()
	tradingEnabled := false
	tmTradingEnabled := TokenMap[tx.To().String()].TradingEnabled
	if tmTradingEnabled != nil {
		tradingEnabled = *TokenMap[tx.To().String()].TradingEnabled
	} else {
		log.Info().Msgf("ActiveTrading: FilterTx, erc20 at address %s not registered", to)
		chainId := tx.ChainId().Int64()
		err := artemis_validator_service_groups_models.InsertERC20TokenInfo(ctx, artemis_autogen_bases.Erc20TokenInfo{
			Address:           to,
			ProtocolNetworkID: int(chainId),
			BalanceOfSlotNum:  -2, // -1 means balanceOf it wasn't cracked within 100 slots, -2 means cracking hasn't been attempted yet
		})
		if err != nil {
			log.Err(err).Msg("ActiveTrading: FilterTx, InsertERC20TokenInfo")
		}
		return nil
	}
	if !tradingEnabled {
		return nil
	}
	return tx
}
