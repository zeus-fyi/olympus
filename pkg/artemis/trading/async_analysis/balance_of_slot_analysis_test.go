package async_analysis

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_mev_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/mev"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
)

func (s *ArtemisRealTimeTradingTestSuite) TestFindERC20BalanceOfSlotNumber() {
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	uni.Web3Client.IsAnvilNode = false
	uni.Web3Client.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	uni.Web3Client.NodeURL = "https://hardhat.zeus.fyi"
	shib2Contract := "0x3cda61B56278842876e7fDD56123d83DBAFAe16C"
	s.ca = NewERC20ContractAnalysis(&uni, shib2Contract)
	s.ca.UserB = s.UserB

	tokens, err := artemis_mev_models.SelectERC20TokensWithoutBalanceOfSlotNums(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(tokens)
	s.ca.UserA.IsAnvilNode = true

	for _, token := range tokens {
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

func (s *ArtemisRealTimeTradingTestSuite) TestFindERC20BalanceOfSlotExp() {
	uni := web3_client.InitUniswapClient(ctx, s.UserA)
	uni.Web3Client.IsAnvilNode = false
	uni.Web3Client.AddBearerToken(s.Tc.ProductionLocalTemporalBearerToken)
	uni.Web3Client.NodeURL = "https://hardhat.zeus.fyi"
	sc := "0x3cda61B56278842876e7fDD56123d83DBAFAe16C"
	s.ca = NewERC20ContractAnalysis(&uni, sc)

	//tmp, err := artemis_oly_contract_abis.LoadNewERC20AbiPayload()
	//s.Assert().Nil(err)
	//holder := "0xed5d4dabE7FACDb1Da56956215f486069913ec99"
	//payload := web3_actions.SendContractTxPayload{
	//	SmartContractAddr: sc,
	//	SendEtherPayload:  web3_actions.SendEtherPayload{},
	//	ContractFile:      "",
	//	ContractABI:       tmp.ContractABI,
	//	MethodName:        "balanceOf",
	//	Params:            []interface{}{holder},
	//}
	s.ca.u.Web3Client.Dial()
	defer s.ca.u.Web3Client.Close()
	//bc, err := artemis_oly_contract_abis.LoadERC20DeployedByteCode()

	//err = s.ca.u.Web3Client.SetCodeOverride(ctx, sc, bc)
	//s.Assert().Nil(err)

	//
	//err := s.ca.FindERC20BalanceOfSlotNumber(ctx)
	//s.Assert().Nil(err)
	b, err := s.ca.u.Web3Client.ReadERC20TokenBalance(ctx, sc, s.UserA.Address().String())
	log.Info().Msgf("b: %v", b)
	slotHex, err := web3_client.GetSlot(s.ca.UserA.Address().String(), new(big.Int).SetUint64(uint64(0)))
	s.Assert().Nil(err)
	value := new(big.Int).SetUint64(uint64(100))
	newBalance := common.LeftPadBytes(value.Bytes(), 32)
	err = s.UserA.HardhatSetStorageAt(ctx, sc, slotHex, common.BytesToHash(newBalance).Hex())
	s.Assert().Nil(err)
	b, err = s.ca.u.Web3Client.ReadERC20TokenBalance(ctx, sc, s.UserA.Address().String())
	s.Assert().Nil(err)
	log.Info().Msgf("b: %v", b)
}
