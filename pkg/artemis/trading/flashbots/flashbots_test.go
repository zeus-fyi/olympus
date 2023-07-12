package artemis_flashbots

import (
	"context"
	"testing"

	"github.com/metachris/flashbotsrpc"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/artemis/web3_client"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
)

var ctx = context.Background()

type FlashbotsTestSuite struct {
	test_suites_encryption.EncryptionTestSuite
	fb FlashbotsClient
}

func (s *FlashbotsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	pkHexString := s.Tc.LocalEcsdaTestPkey
	newAccount, err := accounts.ParsePrivateKey(pkHexString)
	s.Assert().Nil(err)
	w3a := web3_client.NewWeb3Client(s.Tc.MainnetNodeUrl, newAccount)
	uni := web3_client.InitUniswapClient(ctx, w3a)
	s.fb = InitFlashbotsClient(ctx, &uni)
}

// TODO: add real payload
func (s *FlashbotsTestSuite) TestFlashbotsSendBundle() {
	br := flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:          nil,
		BlockNumber:  "",
		MinTimestamp: nil,
		MaxTimestamp: nil,
		RevertingTxs: nil,
	}
	resp, err := s.fb.SendBundle(ctx, br)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)
}

func (s *FlashbotsTestSuite) TestGetFlashbotsBlocksV1() {
	resp, err := s.fb.GetFlashbotsBlocksV1(ctx)
	s.Assert().Nil(err)
	s.Assert().NotNil(resp)
}

func TestFlashbotsTestSuite(t *testing.T) {
	suite.Run(t, new(FlashbotsTestSuite))
}
