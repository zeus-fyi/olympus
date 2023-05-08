package artemis_flashbots

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/gochain/web3/accounts"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_encryption"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
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

	s.fb = InitFlashbotsClient(ctx, s.Tc.MainnetNodeUrl, hestia_req_types.Mainnet, newAccount)
}

func (s *FlashbotsTestSuite) TestFlashbots() {
}

func TestFlashbotsTestSuite(t *testing.T) {
	suite.Run(t, new(FlashbotsTestSuite))
}
