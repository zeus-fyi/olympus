package hydra_eth2_web3signer

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	artemis_hydra_orchestrations_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

type HydraSigningAsyncRequestsTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
}

func (t *HydraSigningAsyncRequestsTestSuite) TestPriorityQueue() {
	t.InitLocalConfigs()
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	awsAuthCfg := aegis_aws_auth.AuthAWS{
		AccountNumber: "",
		Region:        "us-west-1",
		AccessKey:     t.Tc.AwsAccessKeySecretManager,
		SecretKey:     t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_auth.InitHydraSecretManagerAuthAWS(ctx, awsAuthCfg)
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
	}
	inMemFsErr := artemis_validator_signature_service_routing.InitRouteMapInMemFS(ctx)
	t.Require().NoError(inMemFsErr)

	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	initErr := artemis_validator_signature_service_routing.GetServiceAuthAndURLs(ctx, vsi, cctx)
	t.Require().NoError(initErr)

	go func() {
		artemis_validator_signature_service_routing.InitAsyncServiceAuthRoutePolling(ctx, vsi, cctx)
	}()
	go InitAsyncMessageQueuesSyncCommitteeQueues(ctx)
	pubkey := "0xb9f787d2f74ce17f22ad8de1bb936c515ee2112460166d8317b3f2c1f81b69bcd6f84b51a2b9e91abcd51afde3e9bdec"

	pubkeyToUUID := make(map[string]string)
	hexMessage, err := aegis_inmemdbs.RandomHex(10)
	t.Require().NoError(err)
	sr := SignRequest{
		UUID:        uuid.UUID{},
		Pubkey:      pubkey,
		Type:        RANDAO_REVEAL,
		SigningRoot: hexMessage,
	}
	batchSigReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}

	totalRequests := 300
	for i := 0; i < totalRequests; i++ {
		uuidRand := uuid.New()
		sr.UUID = uuidRand
		pubkeyToUUID[sr.Pubkey] = uuidRand.String()
		hexMessage, err = aegis_inmemdbs.RandomHex(10)
		t.Require().NoError(err)
		batchSigReqs.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: hexMessage}
		SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Enqueue(sr)
	}

	fmt.Println(SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Size())
	time.Sleep(100 * time.Millisecond)
	fmt.Println(SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Size())
	time.Sleep(100 * time.Millisecond)
	fmt.Println(SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Size())
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(SyncCommitteeContributionAndProofPriorityQueue.PriorityQueue.Size())
}

func (t *HydraSigningRequestsTestSuite) TestAsyncSignRequest() {
	apps.Pg.InitPG(ctx, t.Tc.ProdLocalDbPgconn)
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: t.Tc.AwsAccessKeySecretManager,
		SecretKey: t.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
	}
	err := artemis_validator_signature_service_routing.InitRouteMapInMemFS(ctx)
	t.Require().Nil(err)

	err = artemis_validator_signature_service_routing.GetServiceAuthAndURLs(ctx, vsi, cctx)
	t.Require().Nil(err)
	pubkey := "0xb9f787d2f74ce17f22ad8de1bb936c515ee2112460166d8317b3f2c1f81b69bcd6f84b51a2b9e91abcd51afde3e9bdec"

	pubkeyToUUID := make(map[string]string)
	hexMessage, err := aegis_inmemdbs.RandomHex(10)
	t.Require().NoError(err)
	sr := SignRequest{
		UUID:        uuid.UUID{},
		Pubkey:      pubkey,
		Type:        RANDAO_REVEAL,
		SigningRoot: hexMessage,
	}
	batchSigReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}

	totalTime := time.Duration(0)
	totalRequests := 30
	for i := 0; i < totalRequests; i++ {
		uuidRand := uuid.New()
		sr.UUID = uuidRand
		pubkeyToUUID[sr.Pubkey] = uuidRand.String()
		hexMessage, err = aegis_inmemdbs.RandomHex(10)
		t.Require().NoError(err)
		batchSigReqs.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: hexMessage}
		start := time.Now()
		err = RequestValidatorSignaturesAsync(ctx, batchSigReqs, pubkeyToUUID)
		resp, _ := WaitForSignature(ctx, sr)
		t.Assert().NotEmpty(resp)
		t.Require().Nil(err)
		end := time.Now()         // Get the current time after the function call
		latency := end.Sub(start) //
		fmt.Println("latency: ", latency)
		totalTime += latency
	}
	fmt.Println("avg latency: ", totalTime/time.Duration(totalRequests))
}

func TestHydraSigningAsyncRequestsTestSuite(t *testing.T) {
	suite.Run(t, new(HydraSigningAsyncRequestsTestSuite))
}
