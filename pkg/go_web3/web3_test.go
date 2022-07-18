package go_web3

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/internal/test_utils/test_suites"
)

type GoWeb3BaseTestSuite struct {
	test_suites.BaseTestSuite
}

func (s *GoWeb3BaseTestSuite) TestGasPriceQuery() {
	gp, err := GetGasPrice(context.Background())

	s.Require().Nil(err)
	s.Require().NotNil(gp)
	fmt.Println(gp)
}

func TestGoWeb3BaseTestSuite(t *testing.T) {
	suite.Run(t, new(GoWeb3BaseTestSuite))
}
