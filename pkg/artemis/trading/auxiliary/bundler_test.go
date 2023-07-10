package artemis_trading_auxiliary

func (t *ArtemisAuxillaryTestSuite) TestExecV2TradeBundle() {
	ta, cmd, _ := t.TestExecV2TradeCall()
	t.Require().Equal(1, len(ta.OrderedTxs))
	bundle := ta.CreateFlashbotsBundle(cmd, "latest")
	t.Require().NotEmpty(bundle)
	t.Require().Equal(1, len(bundle.Txs))
	t.Require().Equal(0, len(ta.OrderedTxs))

	//_, err = ta.universalRouterExecuteTx(ctx, tx)
	//t.Require().Nil(err)
	//fmt.Println("tx", tx.Hash().String())
}
