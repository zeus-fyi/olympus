package hydra_eth2_web3signer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps"
	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	hydra_base_test "github.com/zeus-fyi/olympus/hydra/api/test"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	artemis_validator_signature_service_routing "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/signature_routing"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type HydraSigningAsyncRequestsTestSuite struct {
	hydra_base_test.HydraBaseTestSuite
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
	pubkey := "0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08"

	pubkeyToUUID := make(map[string]string)

	sr := SignRequest{
		UUID:        uuid.UUID{},
		Pubkey:      pubkey,
		Type:        "BLOCK_V2",
		SigningRoot: "0x133f5cee5a36d56ca3085db561375b7f668b12f8d8e971aac8578557ca37635f",
	}
	batchSigReqs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	batchSigReqs.Map[sr.Pubkey] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: sr.SigningRoot}
	pubkeyToUUID[sr.Pubkey] = sr.UUID.String()

	err = RequestValidatorSignaturesAsync(ctx, batchSigReqs, pubkeyToUUID)
	t.Require().Nil(err)

	resp, _ := WaitForSignature(ctx, sr)
	t.Assert().NotEmpty(resp)
}

func TestHydraSigningAsyncRequestsTestSuite(t *testing.T) {
	suite.Run(t, new(HydraSigningAsyncRequestsTestSuite))
}
