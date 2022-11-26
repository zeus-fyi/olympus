package web3_client

import "context"

func (s *Web3ClientTestSuite) TestWeb3SendEther() {
	ctx := context.Background()
	b, err := s.GoerliWeb3User.GetCurrentBalance(ctx)

	s.Require().Nil(err)
	s.Assert().NotNil(b)

}
