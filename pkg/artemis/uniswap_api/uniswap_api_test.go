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
	pairs, err := GetTokenPairsWithVolume(ctx)
	s.Require().Nil(err)
	s.Assert().NotEmpty(pairs)
}

func TestUniswapApiTestSuite(t *testing.T) {
	suite.Run(t, new(UniswapAPITestSuite))
}
