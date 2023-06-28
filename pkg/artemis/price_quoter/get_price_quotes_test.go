package pricequoter

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type PriceQuoterTestSuite struct {
	test_suites_base.TestSuite
}

func (s *PriceQuoterTestSuite) SetupTest() {
	s.InitLocalConfigs()
}

func (s *PriceQuoterTestSuite) TestSendSwapRequest() {
	params := map[string]string{
		"sellAmount": "1000000000000000000",                        // 1 ETH in Wei
		"buyToken":   "0x6b175474e89094c44da98b954eedeac495271d0f", // DAI
		"sellToken":  "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // ETH
	}
	client := NewClient()
	quote, err := client.sendSwapRequest(ctx, "quote", params)
	s.Require().Nil(err)
	s.Assert().NotEmpty(quote)
	fmt.Println(quote)
}

func (s *PriceQuoterTestSuite) TestGetUSDSwapQuote() {
	testToken := "0x6982508145454Ce325dDbE47a25d4ec3d2311933" // Pepe token
	// amount := big.NewInt(1000000000000000000) // 1 token
	// amount := "10000000000000000000" //

	quote, err := GetUSDSwapQuote(ctx, testToken)
	s.Require().Nil(err)
	s.Assert().NotEmpty(quote)
	fmt.Println("USDC Guaranteed Price for 1 PEPE: ", quote.GuaranteedPrice)

}

func TestPriceQuoterTestSuite(t *testing.T) {
	suite.Run(t, new(PriceQuoterTestSuite))
}
