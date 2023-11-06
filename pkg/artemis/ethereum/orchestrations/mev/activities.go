package artemis_mev_transcations

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/gochain/web3/accounts"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	waitForTxRxTimeout    = 15 * time.Minute
	submitSignedTxTimeout = 5 * time.Minute
)

type ArtemisMevActivities struct {
	web3_client.Web3Client
}

func NewArtemisMevActivities(w web3_client.Web3Client) ArtemisMevActivities {
	return ArtemisMevActivities{w}
}

type ActivityDefinition interface{}
type ActivitiesSlice []interface{}

func (d *ArtemisMevActivities) GetActivities() ActivitiesSlice {
	return []interface{}{d.SendEther, d.SubmitSignedTx, d.WaitForTxReceipt, d.BlacklistMinedTxs, d.ConvertMempoolTxs,
		d.GetDynamoDBMempoolTxs, d.ProcessMempoolTxs, d.SimulateAndValidateBundle, d.SubmitFlashbotsBundle, d.RemoveProcessedTx,
		d.HistoricalSimulateAndValidateTx, d.FetchERC20TokenInfo, d.FetchERC20TokenBalanceOfStorageSlot,
		d.CalculateTransferTaxFee, d.BlacklistProcessedTxs, d.GetPostgresMempoolTxs, d.GetLookaheadPrices,
		d.EndServerlessSession, d.MonitorTxStatusReceipts, d.InsertOrUpdateTxReceipts,
	}
}

func (d *ArtemisMevActivities) SendEther(ctx context.Context, payload web3_actions.SendEtherPayload) (common.Hash, error) {
	send, err := d.Send(ctx, payload)
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("tx", send).Interface("payload", payload).Msg("ArtemisEthereumBroadcastTxActivities: Send failed")
		return send.Hash(), err
	}
	return send.Hash(), err
}

func (d *ArtemisMevActivities) SubmitSignedTx(ctx context.Context, signedTx *types.Transaction) (*types.Transaction, error) {
	ctx, cancelFn := context.WithTimeout(ctx, submitSignedTxTimeout)
	defer cancelFn()
	txData, err := d.Web3Actions.SubmitSignedTxAndReturnTxData(ctx, signedTx)
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("signedTx", signedTx).Interface("txData", txData).Msg("ArtemisEthereumBroadcastTxActivities: SubmitSignedTx failed or timed out")
		return nil, err
	}
	return txData, err
}

func (d *ArtemisMevActivities) WaitForTxReceipt(ctx context.Context, hash accounts.Hash) (*types.Receipt, error) {
	ctx, cancelFn := context.WithTimeout(ctx, waitForTxRxTimeout)
	defer cancelFn()
	rx, err := d.WaitForReceipt(ctx, common.Hash(hash))
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("txHash", hash).Interface("rx", rx).Msg("ArtemisEthereumBroadcastTxActivities: WaitForTxReceipt failed or timed out")
		return nil, err
	}
	return rx, err
}
