package artemis_mev_transcations

import (
	"context"
	"time"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/rs/zerolog/log"
	web3_types "github.com/zeus-fyi/gochain/web3/types"
	"github.com/zeus-fyi/gochain/web3/web3_actions"
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
	return []interface{}{d.SendEther, d.SubmitSignedTx, d.WaitForTxReceipt,
		d.GetMempoolTxs, d.DecodeMempoolTxs, d.SimulateAndValidateBundle, d.SubmitFlashbotsBundle}
}

func (d *ArtemisMevActivities) SendEther(ctx context.Context, payload web3_actions.SendEtherPayload) (common.Hash, error) {
	send, err := d.Send(ctx, payload)
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("tx", send).Interface("payload", payload).Msg("ArtemisEthereumBroadcastTxActivities: Send failed")
		return send.Hash, err
	}
	return send.Hash, err
}

func (d *ArtemisMevActivities) SubmitSignedTx(ctx context.Context, signedTx *types.Transaction) (*web3_types.Transaction, error) {
	ctx, cancelFn := context.WithTimeout(ctx, submitSignedTxTimeout)
	defer cancelFn()
	txData, err := d.Web3Actions.SubmitSignedTxAndReturnTxData(ctx, signedTx)
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("signedTx", signedTx).Interface("txData", txData).Msg("ArtemisEthereumBroadcastTxActivities: SubmitSignedTx failed or timed out")
		return nil, err
	}
	return txData, err
}

func (d *ArtemisMevActivities) WaitForTxReceipt(ctx context.Context, hash common.Hash) (*web3_types.Receipt, error) {
	ctx, cancelFn := context.WithTimeout(ctx, waitForTxRxTimeout)
	defer cancelFn()
	rx, err := d.WaitForReceipt(ctx, hash)
	if err != nil {
		log.Err(err).Str("network", d.Network).Str("nodeURL", d.NodeURL).Interface("txHash", hash).Interface("rx", rx).Msg("ArtemisEthereumBroadcastTxActivities: WaitForTxReceipt failed or timed out")
		return nil, err
	}
	return rx, err
}
