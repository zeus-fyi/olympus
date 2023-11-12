package artemis_rawdawg_contract

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func (s *ArtemisTradingContractsTestSuite) TestLocalDeployRawdawgContract() {
	sessionID := fmt.Sprintf("%s-%s", "local-network-session", uuid.New().String())
	//sessionID = fmt.Sprintf("%s-%s", "local-network-session", "12b5d9ce-29dd-4f95-8e89-fed4aef2193d")
	w3a := CreateUser(ctx, "anvil", s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)

	rawdawgAddr, _ := s.testDeployRawdawgContract(w3a)
	s.Require().NotEmpty(rawdawgAddr)
}

func (s *ArtemisTradingContractsTestSuite) testDeployRawdawgContract(w3a web3_actions.Web3Actions) (common.Address, *abi.ABI) {
	rawDawgPayload, bc, err := artemis_oly_contract_abis.LoadLocalRawdawgAbiPayloadV2()
	s.Require().Nil(err)
	s.Require().NotNil(rawDawgPayload)
	rawDawgPayload.Params = []interface{}{}

	err = w3a.SuggestAndSetGasPriceAndLimitForTx(ctx, rawDawgPayload, common.HexToAddress(rawDawgPayload.ToAddress.Hex()))
	s.Require().Nil(err)
	s.Require().NotZero(rawDawgPayload.GasLimit)
	s.Require().NotEmpty(rawDawgPayload.GasFeeCap)
	s.Require().NotEmpty(rawDawgPayload.GasTipCap)
	if w3a.Network == "anvil" {
		rawDawgPayload.GasLimit *= 100
	}
	if w3a.Network == "mainnet" {
		rawDawgPayload.GasFeeCap = artemis_eth_units.MulBigInt(rawDawgPayload.GasFeeCap, big.NewInt(3))
		rawDawgPayload.GasLimit *= 1000
	}
	tx, err := w3a.DeployContract(ctx, bc, *rawDawgPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	time.Sleep(5 * time.Second)
	rx, err := w3a.GetTxReceipt(ctx, tx.Hash().String())
	s.Assert().Nil(err)
	s.Require().NotNil(rx)

	s.Require().NotEmpty(rx.ContractAddress)
	return rx.ContractAddress, rawDawgPayload.ContractABI
}
