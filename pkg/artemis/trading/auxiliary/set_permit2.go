package artemis_trading_auxiliary

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

type AuxiliaryTradingUtils struct {
	web3_client.Web3Client
}

func (a *AuxiliaryTradingUtils) SetPermit2Approval(ctx context.Context, address string) (*types.Transaction, error) {
	tx, err := a.ApprovePermit2(ctx, address)
	if err != nil {
		return tx, err
	}
	return tx, nil
}
