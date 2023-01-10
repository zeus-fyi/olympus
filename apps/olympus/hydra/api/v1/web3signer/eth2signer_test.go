package hydra_eth2_web3signer

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	consensys_eth2_openapi "github.com/zeus-fyi/olympus/hydra/api/v1/web3signer/models"
)

type HydraEth2Web3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *HydraEth2Web3SignerTestSuite) TestEth2Proxy() {
	t.InitLocalConfigs()
	t.E.POST(Eth2SignRoute, Eth2SignRequest)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9000")
	}()

	<-start
	ctx := context.Background()
	defer t.E.Shutdown(ctx)

	att := consensys_eth2_openapi.AttestationSigning{}
	err := faker.FakeData(&att)

	t.Require().Nil(err)
	pubkey := "0x93247f2209abcacf57b75a51dafae777f9dd38bc7053d1af526f220a7489a6d3a2753e5f3e8b1cfe39b56f43611df74a"
	resp, err := t.PostRequest(ctx, Eth2SignRequestWithPubkey(pubkey), att)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func Eth2SignRequestWithPubkey(pubkey string) string {
	return fmt.Sprintf("/api/v1/eth2/sign/%s", pubkey)
}

func TestHydraEth2Web3SignerTestSuite(t *testing.T) {
	suite.Run(t, new(HydraEth2Web3SignerTestSuite))
}
