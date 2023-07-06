package artemis_realtime_trading

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
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

func NewActiveTradingModuleWithoutMetrics(u *web3_client.UniswapClient) ActiveTrading {
	return ActiveTrading{u: u}
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
	mevTxs, err := a.DecodeTx(ctx, tx)
	if err != nil {
		return err
	}
	if len(mevTxs) <= 0 {
		return errors.New("DecodeTx: no txs to process")
	}
	a.m.StageProgressionMetrics.CountPostDecodeTx()
	tfSlice, err := a.ProcessTxs(ctx)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("ProcessTxs: no tx flows to simulate")
	}
	a.m.StageProgressionMetrics.CountPostProcessTx(float64(1))
	err = a.SimTxFilter(ctx, tfSlice)
	if err != nil {
		return err
	}
	if len(tfSlice) <= 0 {
		return errors.New("SimTxFilter: no tx flows to simulate")
	}
	a.m.StageProgressionMetrics.CountPostSimFilterTx(float64(1))
	wc := web3_actions.NewWeb3ActionsClient(artemis_network_cfgs.ArtemisEthereumMainnetQuiknodeLive.NodeURL)
	wc.Dial()
	bn, berr := wc.C.BlockNumber(ctx)
	if berr != nil {
		log.Err(berr).Msg("failed to get block number")
		return berr
	}
	wc.Close()
	for _, tradeFlow := range tfSlice {
		//if tf.InitialPair.PairContractAddr == "" {
		//	a.m.StageProgressionMetrics.CountPostSimTx(tf.InitialPair.PairContractAddr, tf.FrontRunTrade.AmountInAddr.String())
		//} else {
		//	a.m.StageProgressionMetrics.CountPostSimTx(tf.InitialPairV3.PoolAddress, tf.FrontRunTrade.AmountInAddr.String())
		//}
		tradeFlow.CurrentBlockNumber = new(big.Int).SetUint64(bn)
		btf, ber := json.Marshal(tradeFlow)
		if ber != nil {
			log.Err(ber).Msg("failed to marshal tx flow")
			return ber
		}
		baseTx, zerr := tradeFlow.Tx.ConvertToTx()
		if zerr != nil {
			log.Err(zerr).Msg("ActiveTrading: EntryTxFilter, ConvertToTx")
			return zerr
		}
		fromStr := ""
		sender := types.LatestSignerForChainID(baseTx.ChainId())
		from, ferr := sender.Sender(baseTx)
		if ferr != nil {
			log.Err(ferr).Msg("failed to get sender")
			return ferr
		} else {
			fromStr = from.String()
		}

		txStr, terr := json.Marshal(tradeFlow.Tx)
		if terr != nil {
			return terr
		}
		txMempool := artemis_autogen_bases.EthMempoolMevTx{
			ProtocolNetworkID: hestia_req_types.EthereumMainnetProtocolNetworkID,
			Tx:                string(txStr),
			TxFlowPrediction:  string(btf),
			TxHash:            tradeFlow.Tx.Hash,
			Nonce:             int(baseTx.Nonce()),
			From:              fromStr,
			To:                tradeFlow.Tx.To,
			BlockNumber:       int(bn),
		}
		err = artemis_mev_models.InsertMempoolTx(ctx, txMempool)
		if err != nil {
			log.Err(err).Msg("failed to insert mempool tx")
			return err
		}
		a.m.StageProgressionMetrics.CountSavedMempoolTx(float64(1))
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
