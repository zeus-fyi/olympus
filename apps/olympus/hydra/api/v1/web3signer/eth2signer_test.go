package hydra_eth2_web3signer

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
)

type HydraEth2Web3SignerTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *HydraEth2Web3SignerTestSuite) TestQueueResponse() {
	newUUID := uuid.New()
	sr := SignRequest{
		UUID:        newUUID,
		Pubkey:      "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		Type:        AGGREGATION_SLOT,
		SigningRoot: "gadgdsgdsg",
	}
	mockSig := "dsadfgdasf"
	SignatureResponsesCache.Set(newUUID.String(), mockSig, cache.DefaultExpiration)
	resp := ReturnSignedMessage(ctx, sr)
	t.Assert().Equal(mockSig, resp.Signature)
}

func (t *HydraEth2Web3SignerTestSuite) TestAggregationSlot() {
	ags := t.GenerateMockAggregationSlotSigningRequest()
	pubkey := "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288"
	ags.Type = AGGREGATION_SLOT
	resp, err := t.PostRequest(ctx, Eth2SignRequestWithPubkey(pubkey), ags)
	t.Require().Nil(err)
	fmt.Println(resp)
}

func (t *HydraEth2Web3SignerTestSuite) TestAttestation() {
	t.InitLocalConfigs()
	t.E.POST(Eth2SignRoute, HydraEth2SignRequestHandler)
	start := make(chan struct{}, 1)
	go func() {
		close(start)
		_ = t.E.Start(":9000")
	}()

	<-start
	defer t.E.Shutdown(ctx)

	att := t.GenerateMockAttestationSigningRequest()
	pubkey := "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288"
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
