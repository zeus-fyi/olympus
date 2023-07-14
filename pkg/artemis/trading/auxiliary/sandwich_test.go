package artemis_trading_auxiliary

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) testMockSandwichBundle() *AuxiliaryTradingUtils {
	toExchAmount := artemis_eth_units.GweiMultiple(1000)
	//toExchAmount := artemis_eth_units.GweiMultiple(1000)
	ta := t.at2
	cmd := t.testEthToWETH(&ta, toExchAmount)
	// part 1 of bundle
	ctx = ta.CreateFrontRunCtx(ctx)
	tx, err := ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(toExchAmount, tx.Value())
	t.Require().Equal(1, len(ta.MevTxGroup.OrderedTxs))
	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(1, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.MevTxGroup.OrderedTxs))
	ctx = context.Background()
	fmt.Println("frontRun: txGasLimit", tx.Gas())
	fmt.Println("frontRun: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("frontRun: txGasTipCap", tx.GasTipCap().String())
	// middle of bundle

	user := t.at1
	cmd = t.testEthToWETH(&user, toExchAmount)
	ctx = user.CreateUserTradeCtx(ctx)
	tx, err = user.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().NotEmpty(tx)
	fmt.Println("userTrade: txGasLimit", tx.Gas())
	fmt.Println("userTrade: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("userTrade: txGasTipCap", tx.GasTipCap().String())
	err = ta.AddTxToBundleGroup(ctx, tx)
	t.Require().Nil(err)
	signer := types.LatestSignerForChainID(artemis_eth_units.NewBigInt(hestia_req_types.EthereumGoerliProtocolNetworkID))
	sender, err := signer.Sender(tx)
	t.Require().Nil(err)
	t.Require().Equal(user.w3a().Address().String(), sender.String())
	t.Require().Equal(user.w3c().Address().String(), sender.String())

	ctx = context.Background()
	// part 3 of bundle
	cmd = t.testExecV2Trade(&ta, hestia_req_types.Goerli)
	ctx = ta.CreateBackRunCtx(ctx)
	fmt.Println("mainTraderAddr", ta.w3a().Address().String())
	tx, err = ta.universalRouterCmdToTxBuilder(ctx, cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	fmt.Println("backRun: txGasLimit", tx.Gas())
	fmt.Println("backRun: txGasFeeCap", tx.GasFeeCap().String())
	fmt.Println("backRun: txGasTipCap", tx.GasTipCap().String())
	t.Require().Equal(2, len(ta.MevTxGroup.OrderedTxs))

	err = ta.CreateOrAddToFlashbotsBundle(cmd, "latest")
	t.Require().Nil(err)
	t.Require().NotEmpty(ta.Bundle.Txs)
	t.Require().Equal(3, len(ta.Bundle.Txs))
	t.Require().Equal(0, len(ta.MevTxGroup.OrderedTxs))

	return &ta
}
func (t *ArtemisAuxillaryTestSuite) TestSandwichCallBundle() {
	ta := t.testMockSandwichBundle()
	t.Require().NotEmpty(ta)
	resp, err := ta.callFlashbotsBundle(ctx)
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

func (t *ArtemisAuxillaryTestSuite) TestSandwichCallAndSendBundle() {
	ta := t.testMockSandwichBundle()
	t.Require().NotEmpty(ta)
	resp, err := ta.CallAndSendFlashbotsBundle(ctx)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
}
