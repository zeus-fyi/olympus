package eth_validators_service_requests

import (
	"context"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
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
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "agekey",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	keyOne := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x8a7addbf2857a72736205d861169c643545283a74a1ccb71c95dd2c9652acb89de226ca26d60248c4ef9591d7e010288",
		FeeRecipient: "0xF7Ab1d834Cd0A33691e9A750bD720cb6436cA1B9",
	}
	keyTwo := hestia_req_types.ValidatorServiceOrgGroup{
		Pubkey:       "0x8258f4ec23d5e113f2b62caa40d77d52c2ad9dfd871173a9815f77ef66e02e5a090e8e940477c7df06477c5ceb42bb08",
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
