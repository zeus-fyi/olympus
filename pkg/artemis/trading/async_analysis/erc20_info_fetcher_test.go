package async_analysis

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
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
	tokens, err := artemis_mev_models.SelectERC20TokensWithoutMetadata(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(tokens)

	for _, token := range tokens {
		fmt.Println(token.Address)
		s.ca.SmartContractAddr = token.Address
		err = s.ca.FindERC20TokenMetadataInfo(ctx)
		s.Assert().Nil(err)
		time.Sleep(100 * time.Millisecond)
	}
}
