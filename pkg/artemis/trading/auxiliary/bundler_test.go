package artemis_trading_auxiliary

import (
	"context"
	"fmt"

	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

func (t *ArtemisAuxillaryTestSuite) TestCreateFbBundle() (*AuxiliaryTradingUtils, MevTxGroup) {
	ta := t.at2
	t.Require().Equal(t.goerliNode, ta.nodeURL())
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(1000)
	cmd := t.testEthToWETH(&ta, toExchAmount)
	// part 1 of bundle

	tx, _, err := universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().Nil(err)
	t.Require().NotEmpty(tx)
	t.Require().Equal(toExchAmount, tx.Value())
	txStart := TxWithMetadata{
		Tx: tx,
	}

	bundle, err := AddTxToBundleGroup(ctx, txStart, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(bundle.MevTxs)
	t.Require().Equal(1, len(bundle.MevTxs))
	t.Require().Equal(1, len(bundle.OrderedTxs))

	// part 2 of bundle
	ctx = context.WithValue(ctx, web3_actions.NonceOffset, 1)
	cmd, permit2Val := t.testExecV2Trade(&ta, hestia_req_types.Goerli)
	tx, _, err = universalRouterCmdToTxBuilder(ctx, *ta.w3c(), cmd)
	t.Require().NotEmpty(tx)
	txMeta := TxWithMetadata{
		Tx: tx,
	}
	if permit2Val != nil {
		txMeta.Permit2Tx = permit2Val.Permit2Tx
	}
	bundle, err = AddTxToBundleGroup(ctx, txMeta, bundle)
	t.Require().Nil(err)
	t.Require().Equal(2, len(bundle.OrderedTxs))
	t.Require().Equal(2, len(bundle.MevTxs))
	t.Require().NotNil(bundle)
	return &ta, *bundle
}

func (t *ArtemisAuxillaryTestSuite) testExecV2TradeFromUser2(ta *AuxiliaryTradingUtils) *web3_client.UniversalRouterExecCmd {
	t.Require().NotEmpty(ta)
	toExchAmount := artemis_eth_units.GweiMultiple(100000)
	uni := web3_client.InitUniswapClient(ctx, t.goerliWeb3User)
	cmd := web3_client.UniversalRouterExecCmd{
		Commands: []web3_client.UniversalRouterExecSubCmd{},
		Deadline: GetDeadline(),
		Payable:  nil,
	}
	ur, err := cmd.EncodeCommands(ctx, nil)
	t.Require().Nil(err)
	t.Require().NotEmpty(ur)

	fmt.Println(t.goerliWeb3User.IsAnvilNode, uni.DebugPrint, toExchAmount.String())
	return nil
}

func (t *ArtemisAuxillaryTestSuite) TestCallBundle() {
	ta, bundle := t.TestCreateFbBundle()
	t.Require().NotEmpty(ta)
	resp, err := CallFlashbotsBundle(ctx, *ta.w3c(), &bundle, nil)
	t.Require().Nil(err)
	t.Require().NotNil(resp)

	fmt.Println("stateBlockNum", resp.StateBlockNumber)
	fmt.Println("bundleHash", resp.BundleHash)
	fmt.Println("gasFees", resp.GasFees)
	fmt.Println("totalGasUsed", resp.TotalGasUsed)
	t.Require().Equal(2, len(resp.Results))

	for _, sr := range resp.Results {
		t.Require().Emptyf(sr.Error, "error in result: %s", sr.Error)
	}
}

func (t *ArtemisAuxillaryTestSuite) TestCallAndSendBundle() {
	ta, bundle := t.TestCreateFbBundle()
	t.Require().NotEmpty(ta)
	resp, err := CallAndSendFlashbotsBundle(ctx, *ta.w3c(), bundle, nil)
	t.Require().Nil(err)
	t.Require().NotNil(resp)
	fmt.Println("bundleHash", resp.BundleHash)
}
