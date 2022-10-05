package exec_client_apis

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/base"
)

type GoWeb3BaseTestSuite struct {
	base.TestSuite
}

func (s *GoWeb3BaseTestSuite) TestGasPriceQuery() {
	gp, err := GetGasPrice(context.Background())

	s.Require().Nil(err)
	s.Require().NotNil(gp)
	fmt.Println(gp)
}

func (s *GoWeb3BaseTestSuite) TestTxQuery() {
	checkpointHash := "0xc9c4994800171335d7c36c96e1d919fb4ada5c7de6630b21a3ca2d2478659def"
	txData, err := GetTxData(context.Background(), checkpointHash)

	s.Require().Nil(err)
	s.Require().NotNil(txData)
	fmt.Println(txData)
}

func TestGoWeb3BaseTestSuite(t *testing.T) {
	suite.Run(t, new(GoWeb3BaseTestSuite))
}
