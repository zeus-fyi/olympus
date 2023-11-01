package eth_validators_service_requests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	signing_automation_ethereum "github.com/zeus-fyi/zeus/pkg/artemis/web3/signing_automation/ethereum"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	filepaths "github.com/zeus-fyi/zeus/pkg/utils/file_io/lib/v0/paths"
	strings_filter "github.com/zeus-fyi/zeus/pkg/utils/strings"
)

type ValidatorServicesActivitesTestSuite struct {
	test_suites_base.TestSuite
}

var ctx = context.Background()

func (s *ValidatorServicesActivitesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ValidatorServicesActivitesTestSuite) TestVerifyValidatorKeyOwnershipAndSigning() {
	va := ArtemisEthereumValidatorsServiceRequestActivities{}
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	srw := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "demoGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "ageEncryptionKey",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	filter := &strings_filter.FilterOpts{StartsWith: "deposit_data", DoesNotInclude: []string{"keystores.tar.gz.age", ".DS_Store", "keystore.zip"}}
	keystoresPath := filepaths.Path{DirIn: "/Users/alex/go/Olympus/Zeus/builds/serverless/keystores"}
	keystoresPath.FilterFiles = filter
	dpSlice, err := signing_automation_ethereum.ParseValidatorDepositSliceJSON(ctx, keystoresPath)
	if err != nil {
		panic(err)
	}
	pubkeys := make(hestia_req_types.ValidatorServiceOrgGroupSlice, len(dpSlice))
	for i, dp := range dpSlice {
		pubkeys[i] = hestia_req_types.ValidatorServiceOrgGroup{Pubkey: "0x" + dp.Pubkey, FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9"}
	}
	wfParams := ValidatorServiceGroupWorkflowRequest{
		OrgID:                         s.Tc.ProductionLocalTemporalOrgID,
		ServiceRequestWrapper:         srw,
		ValidatorServiceOrgGroupSlice: pubkeys,
	}
	verifiedKeys, err := va.VerifyValidatorKeyOwnershipAndSigning(ctx, wfParams)
	s.Require().Nil(err)
	s.Assert().NotEmpty(verifiedKeys)
}

func TestValidatorServicesActivitesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServicesActivitesTestSuite))
}
