package web3_client

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_test_cache "github.com/zeus-fyi/olympus/pkg/artemis/trading/test_suite/test_cache"
	artemis_oly_contract_abis "github.com/zeus-fyi/olympus/pkg/artemis/web3_client/contract_abis"
)

func (s *Web3ClientTestSuite) TestMulticall3() {
	wc := artemis_test_cache.LiveTestNetwork
	wc.Dial()
	defer wc.Close()

	scAddr := artemis_trading_constants.WETH9ContractAddress
	userAddr := "0x6B44ba0a126a2A1a8aa6cD1AdeeD002e141Bcd44"
	payload := web3_actions.SendContractTxPayload{
		SmartContractAddr: scAddr,
		SendEtherPayload:  web3_actions.SendEtherPayload{},
		ContractABI:       artemis_oly_contract_abis.MustLoadERC20Abi(),
		MethodName:        "balanceOf",
		Params:            []interface{}{userAddr},
	}
	bal, err := wc.CallConstantFunction(ctx, &payload)
	s.Assert().Nil(err)
	s.Assert().NotNil(bal)
	fmt.Println(bal[0].(*big.Int).String())

	userAddr2 := "0xc261aA0F3fe1Ce69972F65397D9d062511efe9cA"
	m3 := []MultiCallElement{{
		Name: "balanceOf",
		Call: Call{
			Target:       common.HexToAddress(artemis_trading_constants.WETH9ContractAddress),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       artemis_oly_contract_abis.MustLoadERC20Abi(),
		DecodedInputs: []interface{}{common.HexToAddress(userAddr)},
	}, {
		Name: "balanceOf",
		Call: Call{
			Target:       common.HexToAddress(artemis_trading_constants.LidoSEthAddr),
			AllowFailure: false,
			Data:         nil,
		},
		AbiFile:       artemis_oly_contract_abis.MustLoadERC20Abi(),
		DecodedInputs: []interface{}{common.HexToAddress(userAddr2)},
	}}
	payload, err = CreateMulticall3Payload(ctx, m3)
	s.Assert().Nil(err)
	resp, err := wc.CallConstantFunction(ctx, &payload)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)
	for _, r := range resp {
		encData := r.([]struct {
			Success    bool    "json:\"success\""
			ReturnData []uint8 "json:\"returnData\""
		})[0]
		fmt.Println(encData.Success)
		bi := new(big.Int).SetBytes(encData.ReturnData)
		fmt.Println(bi.String())
	}

	m := Multicall3{
		Calls:   m3,
		Results: nil,
	}

	mcResp, err := m.PackAndCall(ctx, wc)
	s.Assert().Nil(err)
	s.Assert().NotNil(mcResp)
	s.Require().Len(mcResp, 2)
	for _, mcr := range mcResp {
		s.Assert().True(mcr.Success)
		fmt.Println(mcr.Success)
		fmt.Println(mcr.DecodedReturnData)
	}
}
