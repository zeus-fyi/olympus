package web3_client

import "fmt"

func (s *Web3ClientTestSuite) TestGetLatestBlockTxs() {
	txs, err := s.MainnetWeb3User.GetBlockTxs(ctx)
	s.Require().Nil(err)

	fmt.Println("txs", len(txs))
	for _, tx := range txs {
		fmt.Println("tx", tx.Hash().String())
	}
}
