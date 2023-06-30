package async_analysis

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ArtemisRealTimeTradingTestSuite) TestERC20InfoFetcher() {
	apps.Pg.InitPG(ctx, s.Tc.LocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	err := s.ca.FindERC20TokenMetadataInfo(ctx)
	s.Assert().Nil(err)
}

func (s *ArtemisRealTimeTradingTestSuite) TestERC20InfoFetcherExisting() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	err := s.ca.FindERC20TokenMetadataInfo(ctx)
	s.Assert().Nil(err)

	// SelectERC20Tokens
}
