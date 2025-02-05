package artemis_rawdawg_contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	artemis_trading_constants "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/constants"
	artemis_eth_units "github.com/zeus-fyi/olympus/pkg/artemis/trading/lib/units"
	artemis_trading_types "github.com/zeus-fyi/olympus/pkg/artemis/trading/types"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

func mockedTrade() *artemis_trading_types.TradeOutcome {
	to := &artemis_trading_types.TradeOutcome{
		AmountIn:      artemis_eth_units.Ether,
		AmountInAddr:  artemis_trading_constants.WETH9ContractAddressAccount,
		AmountOut:     artemis_eth_units.NewBigInt(0),
		AmountOutAddr: artemis_trading_constants.LinkTokenAddressAccount,
		//AmountOutAddr: accounts.HexToAddress("0x8647Ae4E646cd3CE37FdEB4591b0A7928254bb73"),
	}
	return to
}

func (s *ArtemisTradingContractsTestSuite) mockConditions(w3a web3_actions.Web3Actions, to *artemis_trading_types.TradeOutcome) (common.Address, *abi.ABI) {
	rawDawgAddr, abiFile := s.testDeployRawdawgContract(w3a)
	err := w3a.SetBalanceAtSlotNumber(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex(), 3, artemis_eth_units.EtherMultiple(10))
	s.Require().Nil(err)
	nbal := hexutil.Big{}
	bigInt := nbal.ToInt()
	bigInt.Set(artemis_eth_units.EtherMultiple(10000000))
	nbal = hexutil.Big(*bigInt)

	err = w3a.SetBalance(ctx, rawDawgAddr.Hex(), nbal)
	s.Require().Nil(err)

	rawDawgWethBal, err := w3a.ReadERC20TokenBalance(ctx, to.AmountInAddr.Hex(), rawDawgAddr.Hex())
	s.Require().Nil(err)
	s.Require().Equal(artemis_eth_units.EtherMultiple(10), rawDawgWethBal)
	return rawDawgAddr, abiFile
}
