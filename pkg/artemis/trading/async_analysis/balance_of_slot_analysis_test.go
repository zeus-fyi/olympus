package async_analysis

import (
	"fmt"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ArtemisRealTimeTradingTestSuite) TestFindERC20BalanceOfSlotNumber() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	tokens, err := artemis_validator_service_groups_models.SelectERC20TokensWithoutBalanceOfSlotNums(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(tokens)

	for _, token := range tokens {
		s.ca.SmartContractAddr = token.Address
		fmt.Println("token.Address", token.Address)
		//err = s.ca.FindERC20BalanceOfSlotNumber(ctx)
		//s.Assert().Nil(err)
		//time.Sleep(100 * time.Millisecond)
	}
}
