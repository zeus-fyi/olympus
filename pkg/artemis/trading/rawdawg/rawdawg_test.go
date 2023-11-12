package artemis_rawdawg_contract

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
	"github.com/zeus-fyi/zeus/examples/adaptive_rpc_load_balancer/smart_contract_library"
	web3_actions "github.com/zeus-fyi/zeus/pkg/artemis/web3/client"
)

var ctx = context.Background()

type ArtemisTradingContractsTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
}

func (s *ArtemisTradingContractsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	apps.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)
}

var LoadBalancerAddress = "https://iris.zeus.fyi/v1/router"

func CreateUser(ctx context.Context, network, bearer, sessionID string) web3_actions.Web3Actions {
	acc, err := accounts.CreateAccount()
	if err != nil {
		panic(err)
	}
	w3a := web3_actions.NewWeb3ActionsClientWithAccount(LoadBalancerAddress, acc)
	w3a.AddAnvilSessionLockHeader(sessionID)
	w3a.AddBearerToken(bearer)
	if network == "mainnet" {
		w3a.AddDefaultEthereumMainnetTableHeader()
	}
	w3a.Network = network
	w3a.IsAnvilNode = true
	nvB := (*hexutil.Big)(smart_contract_library.EtherMultiple(10000))
	w3a.Dial()
	defer w3a.Close()
	err = w3a.SetBalance(ctx, w3a.Address().String(), *nvB)
	if err != nil {
		panic(err)
	}
	return w3a
}

//type QuoteExactInputSingleParams struct {
//	TokenIn           accounts.Address `abi:"tokenIn"`
//	TokenOut          accounts.Address `abi:"tokenOut"`
//	Fee               *big.Int         `abi:"fee"`
//	AmountIn          *big.Int         `abi:"amountIn"`
//	SqrtPriceLimitX96 *big.Int         `abi:"sqrtPriceLimitX96"`
//}
//type UniswapAmountOutV3 struct {
//	AmountOut               *big.Int
//	SqrtPriceX96After       *big.Int
//	InitializedTicksCrossed uint32
//	GasEstimate             *big.Int
//}
//
//func (s *ArtemisTradingContractsTestSuite) TestGetPoolV3ExactInputSingleQuoteFromQuoterV2(ctx context.Context, w3a web3_actions.Web3Actions, qp QuoteExactInputSingleParams) (UniswapAmountOutV3, error) {
//	scInfo := &web3_actions.SendContractTxPayload{
//		SmartContractAddr: "0x61fFE014bA17989E743c5F6cB21bF9697530B21e",
//		SendEtherPayload:  web3_actions.SendEtherPayload{},
//		ContractABI:       artemis_oly_contract_abis.MustLoadQuoterV2Abi(),
//		MethodName:        "quoteExactInputSingle",
//		Params:            []interface{}{qp},
//	}
//	qa := UniswapAmountOutV3{}
//	resp, err := w3a.CallConstantFunction(ctx, scInfo)
//	if err != nil {
//		return qa, err
//	}
//	for i, val := range resp {
//		switch i {
//		case 0:
//			qa.AmountOut = val.(*big.Int)
//		case 1:
//			qa.SqrtPriceX96After = val.(*big.Int)
//		case 2:
//			qa.InitializedTicksCrossed = val.(uint32)
//		case 3:
//			qa.GasEstimate = val.(*big.Int)
//		}
//	}
//	return qa, nil
//}

func TestArtemisTradingContractsTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisTradingContractsTestSuite))
}
