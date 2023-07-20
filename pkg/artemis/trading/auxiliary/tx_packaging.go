package artemis_trading_auxiliary

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/rs/zerolog/log"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_autogen_bases "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/bases/autogen"
	artemis_eth_txs "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models/txs/eth_txs"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

const (
	FrontRun      = "frontRun"
	UserTrade     = "userTrade"
	UserTradeTest = "userTradeTest"
	BackRun       = "backRun"
	TradeType     = "tradeType"

	TradeCfg = "tradeCfg"
	Permit2  = "permit2"
)

func getAdditionalTxConfig(ctx context.Context) string {
	tt := ctx.Value(TradeCfg)
	if tt == nil {
		return ""
	}
	tt = tt.(string)
	switch tt {
	case Permit2:
		return Permit2
	}
	return ""
}

func getTradeTypeFromCtx(ctx context.Context) string {
	tt := ctx.Value(TradeType)
	if tt == nil {
		return ""
	}
	tt = tt.(string)
	switch tt {
	case FrontRun:
		return FrontRun
	case UserTrade:
		return UserTrade
	case BackRun:
		return BackRun
	}
	return ""
}

func CreateFrontRunCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TradeType, FrontRun)
	return ctx
}
func CreateFrontRunCtxWithPermit2(ctx context.Context) context.Context {
	ctx = CreateFrontRunCtx(ctx)
	ctx = context.WithValue(ctx, TradeCfg, Permit2)
	return ctx
}

func CreateBackRunCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TradeType, BackRun)
	ctx = web3_actions.SetNonceOffset(ctx, 1)
	return ctx
}

func CreateBackRunCtxWithPermit2(ctx context.Context, w3c web3_client.Web3Client) context.Context {
	ctx = CreateBackRunCtx(ctx)
	ctx = context.WithValue(ctx, TradeCfg, Permit2)
	return ctx
}

func CreateUserTradeCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TradeType, UserTrade)
	return ctx
}

func packageTxForBundle(ctx context.Context, from string, txWithMetadata TxWithMetadata) (artemis_eth_txs.EthTx, error) {
	mevTx, err := getEthTxByPackageType(ctx, from, txWithMetadata)
	if err != nil {
		log.Err(err).Msg("error getting eth tx by package type")
		return artemis_eth_txs.EthTx{}, err
	}
	return mevTx, nil
}

func getEthTxByPackageType(ctx context.Context, from string, signedTxWithMetadata TxWithMetadata) (artemis_eth_txs.EthTx, error) {
	tt := getAdditionalTxConfig(ctx)
	switch tt {
	case Permit2:
		return packagePermit2Tx(ctx, from, signedTxWithMetadata)
	}
	return packageRegularTx(ctx, from, signedTxWithMetadata)
}

func packagePermit2Tx(ctx context.Context, from string, signedTxWithMetadata TxWithMetadata) (artemis_eth_txs.EthTx, error) {
	permit2 := signedTxWithMetadata.Permit2Tx
	mevTx, err := packageRegularTx(ctx, from, signedTxWithMetadata)
	if err != nil {
		log.Err(err).Msg("packagePermit2Tx: error packaging regular tx")
		return artemis_eth_txs.EthTx{}, err
	}
	mevTx.Permit2Tx = artemis_eth_txs.Permit2Tx{
		Permit2Tx: artemis_autogen_bases.Permit2Tx{
			Nonce:             permit2.Nonce,
			Owner:             permit2.Owner,
			Deadline:          permit2.Deadline,
			Token:             permit2.Token,
			ProtocolNetworkID: permit2.ProtocolNetworkID,
		},
	}
	return mevTx, nil
}

func packageRegularTx(ctx context.Context, from string, signedTxWithMetadata TxWithMetadata) (artemis_eth_txs.EthTx, error) {
	signedTx := signedTxWithMetadata.Tx
	pi := signedTx.ChainId()
	if pi == nil {
		log.Warn().Msg("packageRegularTx: chain id is nil, setting to 1")
		pi = big.NewInt(1)
	}
	gasFeeCap := signedTx.GasFeeCap()
	gasTipCap := signedTx.GasTipCap()
	gasLimit := signedTx.Gas()
	switch signedTxWithMetadata.TradeType {
	case UserTradeTest:
		log.Info().Msg("txGasAdjuster: UserTradeTest gas adjustment")
		gasTipCap = artemis_eth_units.NewBigInt(10)
		gasLimit = gasLimit * 2
	case UserTrade:
		log.Info().Msg("txGasAdjuster: UserTrade gas adjustment")
		gasTipCap = artemis_eth_units.OneTenthGwei
		gasLimit = gasLimit * 2
	}
	ethGas := artemis_autogen_bases.EthTxGas{
		TxHash: signedTx.Hash().String(),
		GasPrice: sql.NullInt64{
			Valid: false,
		},
		GasLimit: sql.NullInt64{
			Int64: int64(gasLimit),
			Valid: true,
		},
		GasTipCap: sql.NullInt64{
			Int64: gasTipCap.Int64(),
			Valid: true,
		},
		GasFeeCap: sql.NullInt64{
			Int64: gasFeeCap.Int64(),
			Valid: true,
		},
	}
	signerType := int(signedTx.Type())
	typeEnum := "0x02"
	if signerType == 1 {
		typeEnum = "0x01"
		ethGas.GasPrice = sql.NullInt64{
			Int64: signedTx.GasPrice().Int64(),
			Valid: true,
		}
		ethGas.GasLimit = sql.NullInt64{
			Int64: int64(gasLimit),
			Valid: true,
		}
		ethGas.GasFeeCap = sql.NullInt64{
			Int64: gasFeeCap.Int64(),
			Valid: false,
		}
		ethGas.GasTipCap = sql.NullInt64{
			Int64: gasTipCap.Int64(),
			Valid: false,
		}
	}
	mevTx := artemis_eth_txs.EthTx{
		EthTx: artemis_autogen_bases.EthTx{
			ProtocolNetworkID: int(pi.Int64()),
			TxHash:            signedTx.Hash().String(),
			Nonce:             int(signedTx.Nonce()),
			From:              from,
			Type:              typeEnum,
		},
		EthTxGas: ethGas,
	}
	return mevTx, nil
}
