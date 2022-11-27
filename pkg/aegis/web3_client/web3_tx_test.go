package web3_client

import (
	"github.com/zeus-fyi/gochain/web3/web3_actions"
)

func (s *Web3ClientTestSuite) TestWeb3SendEther() {
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)

	params := web3_actions.SendEtherPayload{
		Amount:         Finney,
		ToAddress:      s.GoerliWeb3User2.Address(),
		GasPriceLimits: web3_actions.GasPriceLimits{},
	}
	send, err := s.GoerliWeb3User.Send(ctx, params)
	s.Require().Nil(err)
	s.Assert().NotNil(send)
}
