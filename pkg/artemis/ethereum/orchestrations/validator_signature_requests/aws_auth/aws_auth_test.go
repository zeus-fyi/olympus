package artemis_hydra_orchestrations_aws_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	aegis_aws_secretmanager "github.com/zeus-fyi/zeus/pkg/aegis/aws/secretmanager"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
)

type ArtemisHydraSecretsManagerTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ArtemisHydraSecretsManagerTestSuite) SetupTest() {
	s.InitLocalConfigs()
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestUpdateSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "agekey",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	b, err := json.Marshal(v)
	s.Require().Nil(err)
	name := fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID))
	si := secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(name),
		SecretBinary: b,
	}
	_, err = HydraSecretManagerAuthAWS.UpdateSecret(ctx, &si)
	s.Require().Nil(err)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestCreateSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "agekey",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	b, err := json.Marshal(v)
	s.Require().Nil(err)
	name := fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID))
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(name),
		Description:  aws.String(fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID))),
		SecretBinary: b,
		SecretString: nil,
	}
	_, err = HydraSecretManagerAuthAWS.CreateSecret(ctx, &si)
	errStr := err.Error()
	errCheckStr := fmt.Sprintf("the secret %s already exists", name)
	if strings.Contains(errStr, errCheckStr) {
		fmt.Println("Secret already exists, skipping")
	} else {
		s.Require().Nil(err)
	}
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestFetchSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "agekey",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	si := aegis_aws_secretmanager.SecretInfo{
		Region: "us-west-1",
		Name:   fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)),
	}
	so, err := GetServiceRoutesAuths(ctx, si)
	s.Require().Nil(err)
	s.Require().NotEmpty(so)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestFetchServiceRoutesAuths() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			ServiceURL: s.Tc.AwsLamdbaTestURL,
			SecretName: "testLambdaExternalSecret",
			AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
			SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	si := aegis_aws_secretmanager.SecretInfo{
		Region: "us-west-1",
		Name:   fmt.Sprintf("%s-%d-%s", v.GroupName, ou.OrgID, hestia_req_types.ProtocolNetworkIDToString(v.ProtocolNetworkID)),
	}
	srw, err := GetServiceRoutesAuths(ctx, si)
	s.Require().Nil(err)
	s.Require().NotEmpty(srw)

	s.Assert().Equal(v.GroupName, srw.GroupName)
	s.Assert().Equal(v.ServiceAuth, srw.ServiceAuth)
}

func TestArtemisHydraSecretsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisHydraSecretsManagerTestSuite))
}
