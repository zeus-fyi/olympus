package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func (s *Web3ClientTestSuite) TestGetLatestBlockTxs() {
	txs, err := s.MainnetWeb3User.GetBlockTxs(ctx)
	s.Require().Nil(err)

	fmt.Println("txs", len(txs))
	for _, tx := range txs {
		fmt.Println("tx", tx.Hash().String())
	}
}

func (s *Web3ClientTestSuite) TestGetTxByHash() {
	hashStr := "0xb841ae58afb7c6e0e7c321e2d151d93599dfd826ac3835f3c7cd8c029b6fd9a7"
	tx, _, err := s.MainnetWeb3User.GetTxByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)
}
