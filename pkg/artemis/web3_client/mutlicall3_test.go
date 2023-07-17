package web3_client

import (
	"fmt"
	"math/big"

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

	//scAddr := artemis_trading_constants.Multicall3Address
	//userAddr := "0x6B44ba0a126a2A1a8aa6cD1AdeeD002e141Bcd44"
	//payload := web3_actions.SendContractTxPayload{
	//	SmartContractAddr: scAddr,
	//	SendEtherPayload:  web3_actions.SendEtherPayload{},
	//	ContractABI:       artemis_oly_contract_abis.MustLoadERC20Abi(),
	//	MethodName:        "balanceOf",
	//	Params:            []interface{}{userAddr},
	//}
}
