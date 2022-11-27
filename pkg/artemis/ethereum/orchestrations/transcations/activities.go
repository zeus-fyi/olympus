package artemis_ethereum_transcations

import (
	"context"
	"time"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/rs/zerolog/log"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

type ArtemisEthereumBroadcastTxActivities struct {
}
type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *ArtemisEthereumBroadcastTxActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.SendEther, d.SubmitSignedTxAndReturnTxData, d.WaitForReceipt}
}

func (d *ArtemisEthereumBroadcastTxActivities) SendEther(ctx context.Context, payload web3_actions.SendEtherPayload) (common.Hash, error) {
	send, err := ArtemisEthereumBroadcastTxClient.Send(ctx, payload)
	if err != nil {
		log.Err(err).Interface("tx", send).Interface("payload", payload).Msg("ArtemisEthereumBroadcastTxActivities: Send failed")
		return send.Hash, err
	}
	return send.Hash, err
}

func (d *ArtemisEthereumBroadcastTxActivities) SubmitSignedTxAndReturnTxData(ctx context.Context, signedTx *types.Transaction) (*web3_types.Transaction, error) {
	ctx, cancelFn := context.WithTimeout(ctx, 5*time.Minute)
	defer cancelFn()
	txData, err := ArtemisEthereumBroadcastTxClient.SubmitSignedTxAndReturnTxData(ctx, signedTx)
	if err != nil {
		log.Err(err).Interface("signedTx", signedTx).Interface("txData", txData).Msg("ArtemisEthereumBroadcastTxActivities: SubmitSignedTxAndReturnTxData failed or timed out")
		return nil, err
	}
	return txData, err
}

func (d *ArtemisEthereumBroadcastTxActivities) WaitForReceipt(ctx context.Context, hash common.Hash) (*web3_types.Receipt, error) {
	ctx, cancelFn := context.WithTimeout(ctx, 15*time.Minute)
	defer cancelFn()
	rx, err := ArtemisEthereumBroadcastTxClient.WaitForReceipt(ctx, hash)
	if err != nil {
		log.Err(err).Interface("txHash", hash).Interface("rx", rx).Msg("ArtemisEthereumBroadcastTxActivities: WaitForReceipt failed or timed out")
		return nil, err
	}
	return rx, err
}
