package artemis_trading_types

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
	web3_actions "github.com/zeus-fyi/gochain/web3/client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type TradingTypesTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
}

func (s *TradingTypesTestSuite) TestTxConverter() {
	hashStr := "0xb841ae58afb7c6e0e7c321e2d151d93599dfd826ac3835f3c7cd8c029b6fd9a7"
	wc := web3_actions.NewWeb3ActionsClient(s.Tc.MainnetNodeUrl)
	wc.Dial()
	defer wc.Close()
	tx, _, err := wc.C.TransactionByHash(ctx, common.HexToHash(hashStr))
	s.Require().Nil(err)
	s.Require().NotNil(tx)

	newTx := JSONTx{}
	err = newTx.UnmarshalTx(tx)
	s.Require().Nil(err)
	s.Require().NotEmpty(newTx)
	fmt.Println("newTx", newTx)

	conTx, err := newTx.ConvertToTx()
	s.Require().Nil(err)
	s.Require().NotEmpty(conTx)
	chainID := conTx.ChainId()
	fmt.Println("chainID", chainID)
	fmt.Println("conTx", conTx.Hash().Hex())

	fromStr := ""
	sender := types.LatestSignerForChainID(chainID)
	from, ferr := sender.Sender(conTx)
	s.Require().Nil(ferr)
	fromStr = from.Hex()
	fmt.Println("fromStr", fromStr)

}

func (s *TradingTypesTestSuite) SetupTest() {
	s.InitLocalConfigs()

}

func TestTradingTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TradingTypesTestSuite))
}
