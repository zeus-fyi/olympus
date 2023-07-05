package artemis_realtime_trading

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	metrics_trading "github.com/zeus-fyi/olympus/pkg/apollo/ethereum/mev/trading"
	artemis_network_cfgs "github.com/zeus-fyi/olympus/pkg/artemis/configs"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type ActiveTrading struct {
	u *web3_client.UniswapClient
	m metrics_trading.TradingMetrics
}

func NewActiveTradingModule(u *web3_client.UniswapClient, tm metrics_trading.TradingMetrics) ActiveTrading {
	return ActiveTrading{u, tm}
}

func (a *ActiveTrading) IngestTx(ctx context.Context, tx *types.Transaction) error {
	a.m.StageProgressionMetrics.CountPreEntryFilterTx()
	err := a.EntryTxFilter(ctx, tx)
	if err != nil {
		return err
	}
	a.m.StageProgressionMetrics.CountPostEntryFilterTx()
	_, err = a.DecodeTx(ctx, tx)
	if err != nil {
		return err
	}
	a.m.StageProgressionMetrics.CountPostDecodeTx()
	tfSlice, err := a.ProcessTxs(ctx)
	if err != nil {
		return err
	}
	a.m.StageProgressionMetrics.CountPostProcessTx(float64(1))
	err = a.SimTxFilter(ctx, tfSlice)
	if err != nil {
		return err
	}
	a.m.StageProgressionMetrics.CountPostSimFilterTx(float64(1))
	log.Info().Msg("tx passed all filters")
	wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
	wc.Dial()
	bn, berr := wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return berr
	}
	defer wc.Close()
	for _, tf := range tfSlice {
		//if tf.InitialPair.PairContractAddr == "" {
		//	a.m.StageProgressionMetrics.CountPostSimTx(tf.InitialPair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String())
		//} else {
		//	a.m.StageProgressionMetrics.CountPostSimTx(tf.InitialPairV3.PoolAddress, tf.FrontRunTrade.AmountInAddr.String())
		//}
		btf, ber := json.Marshal(tf)
		if ber != nil {
			log.Err(ber).Msg("failed to marshal tx flow")
			return ber
		}
		fromStr := ""
		sender := types.LatestSignerForChainID(tf.Tx.ChainId())
		from, ferr := sender.Sender(tf.Tx)
		if ferr != nil {
			log.Err(ferr).Msg("failed to get sender")
		} else {
			fromStr = from.String()
		}

		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                tx.Hash().String(),
			TxFlowPrediction:  string(btf),
			TxHash:            tx.Hash().String(),
			Nonce:             int(tx.Nonce()),
			From:              fromStr,
			To:                tx.To().String(),
			BlockNumber:       int(bn),
		}
		err = artemis_validator_service_groups_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			log.Err(err).Msg("failed to insert mempool tx")
			return err
		}
	}

	return err
}

func (a *ActiveTrading) ProcessTx(ctx context.Context, tx *types.Transaction) error {
	err := a.SimulateTx(ctx, tx)
	if err != nil {
		return err
	}
	a.SendToBundleStack(ctx, tx)
	return nil
}
