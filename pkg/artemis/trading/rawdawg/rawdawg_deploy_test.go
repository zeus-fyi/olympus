package artemis_rawdawg_contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (s *ArtemisTradingContractsTestSuite) TestDeployRawdawgContract() {
	sessionID := fmt.Sprintf("%s-%s", "local-network-session", uuid.New().String())
	//sessionID := fmt.Sprintf("%s-%s", "local-network-session", "12b5d9ce-29dd-4f95-8e89-fed4aef2193d")
	w3a := CreateLocalUser(ctx, s.Tc.ProductionLocalTemporalBearerToken, sessionID)
	defer func(sessionID string) {
		err := w3a.EndAnvilSession()
		s.Require().Nil(err)
	}(sessionID)
	rawDawgPayload, bc, err := artemis_oly_contract_abis.LoadLocalRawdawgAbiPayload()
	s.Require().Nil(err)
	s.Require().NotNil(rawDawgPayload)
	rawDawgPayload.Params = []interface{}{}

	err = w3a.SuggestAndSetGasPriceAndLimitForTx(ctx, rawDawgPayload, common.HexToAddress(rawDawgPayload.ToAddress.Hex()))
	s.Require().Nil(err)
	s.Require().NotZero(rawDawgPayload.GasLimit)
	s.Require().NotEmpty(rawDawgPayload.GasFeeCap)
	s.Require().NotEmpty(rawDawgPayload.GasTipCap)
	rawDawgPayload.GasLimit *= 100

	tx, err := w3a.DeployContract(ctx, bc, *rawDawgPayload)
	s.Require().Nil(err)
	s.Assert().NotNil(tx)

	rx, err := w3a.WaitForReceipt(ctx, tx.Hash())
	s.Assert().Nil(err)
	s.Assert().NotNil(rx)

	s.Require().NotEmpty(rx.ContractAddress)
	fmt.Println(rx.ContractAddress.String())
}
