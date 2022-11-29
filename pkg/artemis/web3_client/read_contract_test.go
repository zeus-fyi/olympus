package web3_client

import "github.com/gochain/gochain/v4/common"

const LinkGoerliContractAddr = "0x326C977E6efc84E512bB9C30f76E30c160eD06FB"

var LinkGoerliContractAddress = common.StringToAddress(LinkGoerliContractAddr)

func (s *Web3ClientTestSuite) TestReadERC20TokenBalance() {
	b, err := s.GoerliWeb3User.ReadERC20TokenBalance(ctx, LinkGoerliContractAddress, s.GoerliWeb3User.PublicKey())
	s.Require().Nil(err)
	s.Assert().NotZero(b)
}

func (s *Web3ClientTestSuite) TestReadERC20ContractDecimals() {
	dec, err := s.GoerliWeb3User.GetContractDecimals(ctx, LinkGoerliContractAddress)
	s.Require().Nil(err)
	s.Assert().Equal(int32(18), dec)
}
