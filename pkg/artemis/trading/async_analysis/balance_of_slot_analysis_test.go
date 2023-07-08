package async_analysis

import (
	"fmt"
	"time"

	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ArtemisRealTimeTradingTestSuite) TestFindERC20BalanceOfSlotNumber() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	uni.Web3Client.IsAnvilNode = true
	shib2Contract := "0x34ba042827996821CFFEB06477D48a2Ff9474483"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB
	tokens, err := artemis_mev_models.SelectERC20TokensWithoutBalanceOfSlotNums(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(tokens)
	s.ca.UserA.IsAnvilNode = true

	for _, token := range tokens {
		err = s.ca.UserA.HardHatResetNetwork(ctx, 17624181)
		if token.BalanceOfSlotNum == -1 {
			continue
		}
		s.ca.SmartContractAddr = token.Address
		fmt.Println("token.Address", token.Address)
		fmt.Println("token.Name", token.Name)
		fmt.Println("token.Symbol", token.Symbol)
		fmt.Println("token.BalanceOfSlotNum", token.BalanceOfSlotNum)
		err = s.ca.FindERC20BalanceOfSlotNumber(ctx)
		s.Assert().Nil(err)
		time.Sleep(100 * time.Millisecond)
	}
}
