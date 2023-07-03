package uniswap_api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type UniswapAPITestSuite struct {
	test_suites_base.TestSuite
}

func (s *UniswapAPITestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *UniswapAPITestSuite) TestGetTokenPairsWithVolume() {
	pairs, err := GetTokenPairsWithVolume(ctx, 10, 1000000, 50000) // 10 pairs, 1M volume, 50k liquidity
	s.Require().Nil(err)
	s.Assert().NotEmpty(pairs)
}

func (s *UniswapAPITestSuite) TestGetPairsForToken() {
	pairs, err := GetPairsForToken(ctx, 100, "0x6b175474e89094c44da98b954eedeac495271d0f") // DAI token
	s.Require().Nil(err)
	s.Assert().NotEmpty(pairs)
}

func TestUniswapApiTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapAPITestSuite))
}
