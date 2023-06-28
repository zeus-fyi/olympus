package artemis_realtime_trading

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

/*
  adding in other filters here
	  - filter by token
	  - filter by profit
	  - filter by risk score
	  - adds sourcing of new blocks
*/

func (a *ActiveTrading) FilterTx(ctx context.Context, tx *types.Transaction) {

}
