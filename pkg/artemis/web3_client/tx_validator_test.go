package web3_client

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func (s *Web3ClientTestSuite) TestValidateTxIsNotPending() {
	txHashStr := "0xee60a7b0fa580b3968906cfa5b5d3f2c2066e221acb330e9dfaece7366179831"
	txHash := common.HexToHash(txHashStr)
	fmt.Println("txHash", txHash.String())
	isPendng, err := s.MainnetWeb3UserExternal.ValidateTxIsPending(ctx, txHashStr)
	s.Require().Nil(err)
	s.Assert().False(isPendng)
}

//func (s *Web3ClientTestSuite) TestValidateTxIsPending() {
//	txHashStr := "0xdb5498fc98e16ba005251d6969c361af98d4e47500a3fbae57f119fa8763c289"
//	txHash := common.HexToHash(txHashStr)
//	fmt.Println("txHash", txHash.String())
//	isPendng, err := s.MainnetWeb3UserExternal.ValidateTxIsPending(ctx, txHashStr)
//	s.Require().Nil(err)
//	s.Assert().False(isPendng)
//}
