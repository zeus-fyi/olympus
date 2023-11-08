package artemis_trading_auxiliary

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestSandwichCallBundle() {
	ta, bundle := t.testMockSandwichBundle()
	t.Require().NotEmpty(ta)
	resp, err := CallFlashbotsBundle(ctx, *ta.w3c(), &bundle, nil)
	t.Require().Nil(err)
	t.Require().NotNil(resp)

	fmt.Println("stateBlockNum", resp.StateBlockNumber)
	fmt.Println("bundleHash", resp.BundleHash)
	fmt.Println("gasFees", resp.GasFees)
	fmt.Println("totalGasUsed", resp.TotalGasUsed)
	t.Require().Equal(3, len(resp.Results))

	for _, sr := range resp.Results {
		t.Require().Emptyf(sr.Error, "error in result: %s", sr.Error)
	}
}

func (t *ArtemisAuxillaryTestSuite) testMockSandwichBundle() (*AuxiliaryTradingUtils, MevTxGroup) {
	toExchAmount := artemis_eth_units.GweiMultiple(10000)
	//toExchAmount := artemis_eth_units.GweiMultiple(1000)
	ta := t.at2
	cmd := t.testEthToWETH(&ta, toExchAmount)

	w3c := *ta.w3c()
	// part 1 of bundle
	frontRunCtx := CreateFrontRunCtx(context.Background())
	scInfoFrontRun, err := universalRouterCmdToUnsignedTxPayload(frontRunCtx, cmd)
	t.Require().Nil(err)
	err = w3c.SuggestAndSetGasPriceAndLimitForTx(ctx, scInfoFrontRun, common.HexToAddress(scInfoFrontRun.SmartContractAddr))
	t.Require().Nil(err)
	scInfoFrontRun.GasTipCap = artemis_eth_units.NewBigInt(0)
	frontRunTx, err := w3c.GetSignedTxToCallFunctionWithData(ctx, scInfoFrontRun, scInfoFrontRun.Data)
	t.Require().Nil(err)
	t.Require().NotEmpty(frontRunTx)
	t.Require().Equal(toExchAmount, frontRunTx.Value())
	fmt.Println("frontRun: txGasLimit", frontRunTx.Gas())
	fmt.Println("frontRun: txGasFeeCap", frontRunTx.GasFeeCap().String())
	fmt.Println("frontRun: txGasTipCap", frontRunTx.GasTipCap().String())
	t.Require().Equal("0", frontRunTx.GasTipCap().String())

	frontRunTxWithMetadata := TxWithMetadata{
		TradeType: FrontRun,
		Tx:        frontRunTx,
		ScPayload: scInfoFrontRun,
	}
	ctx = context.Background()
	bundle, err := AddTxToBundleGroup(ctx, frontRunTxWithMetadata, nil)
	t.Require().Nil(err)
	t.Require().Equal(1, len(bundle.MevTxs))
	t.Require().Equal(1, len(bundle.OrderedTxs))

	// middle of bundle
	user := t.at1
	fmt.Println("userTrader", user.tradersAccount().PublicKey())
	cmd = t.testEthToWETH(&user, toExchAmount)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *user.w3c(), cmd)
	t.Require().NotEmpty(tx)
	userTxWithMetadata := TxWithMetadata{
		TradeType: UserTradeTest,
		Tx:        tx,
	}
	bundle, err = AddTxToBundleGroup(ctx, userTxWithMetadata, bundle)
	t.Require().Nil(err)
	t.Require().Equal(2, len(bundle.MevTxs))
	t.Require().Equal(2, len(bundle.OrderedTxs))
	adjustedTx := bundle.MevTxs[1]
	fmt.Println("userTrade: txGasLimit", adjustedTx.GasLimit.Int64)
	fmt.Println("userTrade: txGasFeeCap", adjustedTx.GasFeeCap.Int64)
	fmt.Println("userTrade: txGasTipCap", adjustedTx.GasTipCap.Int64)
	t.Require().Equal(artemis_eth_units.NewBigInt(10).Int64(), adjustedTx.GasTipCap.Int64)

	signer := types.LatestSignerForChainID(artemis_eth_units.NewBigInt(hestia_req_types.EthereumGoerliProtocolNetworkID))
	sender, err := signer.Sender(tx)
	t.Require().Nil(err)
	t.Require().Equal(user.w3a().Address().String(), sender.String())
	t.Require().Equal(user.w3c().Address().String(), sender.String())

	t.Require().NotEqual(t.at1.tradersAccount().PublicKey(), t.at2.tradersAccount().PublicKey())
	ctx = context.Background()
	// part 3 of bundle
	cmd, pt := t.testExecV2Trade(&ta, hestia_req_types.Goerli)

	scInfoSand, err := universalRouterCmdToUnsignedTxPayload(ctx, cmd)
	t.Require().Nil(err)

	fmt.Println("frontRunGas", frontRunTxWithMetadata.ScPayload.GasLimit)
	backRunGasLimit := frontRunTxWithMetadata.ScPayload.GasLimit * 2
	scInfoSand.GasLimit = backRunGasLimit
	fmt.Println("backRunGas", scInfoSand.GasLimit)
	scInfoSand.GasTipCap = artemis_eth_units.MulBigIntWithInt(frontRunTxWithMetadata.ScPayload.GasFeeCap, 2)
	scInfoSand.GasFeeCap = artemis_eth_units.MulBigIntWithInt(frontRunTxWithMetadata.ScPayload.GasFeeCap, 2)
	ctx = CreateBackRunCtx(ctx)
	fmt.Println("mainTraderAddr", ta.w3a().Address().String())

	sandwichSignedTx, err := ta.w3a().GetSignedTxToCallFunctionWithData(ctx, scInfoSand, scInfoSand.Data)
	t.Require().Nil(err)
	t.Require().NotNil(sandwichSignedTx)

	t.Require().Equal(backRunGasLimit, sandwichSignedTx.Gas())
	fmt.Println("backRun: txGasFeeCap", sandwichSignedTx.GasFeeCap().String())
	fmt.Println("backRun: txGasTipCap", sandwichSignedTx.GasTipCap().String())
	t.Require().Equal(2, len(bundle.OrderedTxs))

	txWithMetadata := TxWithMetadata{
		TradeType: BackRun,
		Tx:        sandwichSignedTx,
		Permit2Tx: pt.Permit2Tx,
	}
	bundle, err = AddTxToBundleGroup(ctx, txWithMetadata, bundle)
	t.Require().Nil(err)
	t.Require().Equal(3, len(bundle.MevTxs))
	t.Require().Equal(3, len(bundle.OrderedTxs))
	t.Require().NotNil(bundle)
	return &ta, *bundle
}

//func (t *ArtemisAuxillaryTestSuite) TestSandwichCallAndSendBundle() {
//	ta, bundle := t.testMockSandwichBundle()
//	t.Require().NotEmpty(ta)
//	resp, err := ta.CallAndSendFlashbotsBundle(ctx, bundle)
//	t.Require().Nil(err)
//	t.Require().NotNil(resp)
//}
