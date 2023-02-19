package eth_validators_service_requests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
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
	keyOne := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x913d41b26a157bc8f539a9f63695b87a066f5086f259673f602a85cf9be0738629e872efd94eda6b08ecfd3c229e875e",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	keyTwo := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0xabaf170036e7cb6674f146f3e3398d45c951e10c8e4f02fc5b062dd91701e5a45554070d0eceeda0d99ac2d11c4543f3",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	pubkeys := hestia_req_types.ValidatorServiceOrgGroupSlice{keyOne, keyTwo}
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
