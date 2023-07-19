package artemis_trading_auxiliary

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) testMockSandwichBundle() (*AuxiliaryTradingUtils, MevTxGroup) {
	toExchAmount := artemis_eth_units.GweiMultiple(10000)
	//toExchAmount := artemis_eth_units.GweiMultiple(1000)
	ta := t.at2
	cmd := t.testEthToWETH(&ta, toExchAmount)
	// part 1 of bundle
	ctx = CreateFrontRunCtx(ctx)
	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(toExchAmount, tx.Value())
	txWithMetadata := TxWithMetadata{
		Tx: tx,
	}
	bundle, err := AddTxToBundleGroup(ctx, txWithMetadata, nil)
	t.Require().Nil(err)
	t.Require().Equal(1, len(bundle.MevTxs))
	t.Require().Equal(1, len(bundle.OrderedTxs))
	ctx = context.Background()
	fmt.Println("frontRun: txGasLimit", tx.Gas())
	fmt.Println("frontRun: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("frontRun: txGasTipCap", tx.GasTipCap().String())
	// middle of bundle

	user := t.at1
	fmt.Println("userTrader", user.tradersAccount().PublicKey())
	cmd = t.testEthToWETH(&user, toExchAmount)
	ctx = CreateUserTradeCtx(ctx)
	tx, _, err = universalRouterCmdToTxBuilder(ctx, *user.w3c(), cmd)
	t.Require().NotEmpty(tx)
	fmt.Println("userTrade: txGasLimit", tx.Gas())
	fmt.Println("userTrade: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("userTrade: txGasTipCap", tx.GasTipCap().String())
	txWithMetadata = TxWithMetadata{
		Tx: tx,
	}
	bundle, err = AddTxToBundleGroup(ctx, txWithMetadata, bundle)
	t.Require().Nil(err)
	signer := types.LatestSignerForChainID(artemis_eth_units.NewBigInt(hestia_req_types.EthereumGoerliProtocolNetworkID))
	sender, err := signer.Sender(tx)
	t.Require().Nil(err)
	t.Require().Equal(user.w3a().Address().String(), sender.String())
	t.Require().Equal(user.w3c().Address().String(), sender.String())

	t.Require().NotEqual(t.at1.tradersAccount().PublicKey(), t.at2.tradersAccount().PublicKey())
	ctx = context.Background()
	// part 3 of bundle
	cmd, pt := t.testExecV2Trade(&ta, hestia_req_types.Goerli)
	ctx = CreateBackRunCtx(ctx, *ta.w3c())
	fmt.Println("mainTraderAddr", ta.w3a().Address().String())
	tx, _, err = universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	fmt.Println("backRun: txGasLimit", tx.Gas())
	fmt.Println("backRun: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("backRun: txGasTipCap", tx.GasTipCap().String())
	t.Require().Equal(2, len(bundle.OrderedTxs))

	txWithMetadata = TxWithMetadata{
		Tx:        tx,
		Permit2Tx: pt.Permit2Tx,
	}
	bundle, err = AddTxToBundleGroup(ctx, txWithMetadata, bundle)
	t.Require().Nil(err)
	t.Require().Equal(3, len(bundle.MevTxs))
	t.Require().Equal(3, len(bundle.OrderedTxs))
	t.Require().NotNil(bundle)
	return &ta, *bundle
}
func (t *ArtemisAuxillaryTestSuite) TestSandwichCallBundle() {
	ta, bundle := t.testMockSandwichBundle()
	t.Require().NotEmpty(ta)
	resp, err := CallFlashbotsBundle(ctx, *ta.w3c(), &bundle)
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

//func (t *ArtemisAuxillaryTestSuite) TestSandwichCallAndSendBundle() {
//	ta, bundle := t.testMockSandwichBundle()
//	t.Require().NotEmpty(ta)
//	resp, err := ta.CallAndSendFlashbotsBundle(ctx, bundle)
//	t.Require().Nil(err)
//	t.Require().NotNil(resp)
//}
