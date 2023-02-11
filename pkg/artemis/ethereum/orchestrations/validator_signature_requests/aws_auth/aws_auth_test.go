package artemis_hydra_orchestrations_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
)

type ArtemisHydraSecretsManagerTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ArtemisHydraSecretsManagerTestSuite) SetupTest() {
	s.InitLocalConfigs()
	auth := aws_secrets.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestCreateSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID + 1
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceURL:        s.Tc.AwsLamdbaTestURL,
		ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
			SecretName:   "testLambdaExternalSecret",
			AccessKey:    s.Tc.AwsAccessKeyLambdaExt,
			AccessSecret: s.Tc.AwsSecretKeyLambdaExt,
		}},
	}
	b, err := json.Marshal(v)
	s.Require().Nil(err)
	si := secretsmanager.CreateSecretInput{
		Name:         aws.String(fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID)),
		Description:  aws.String(fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID)),
		SecretBinary: b,
		SecretString: nil,
	}
	_, err = HydraSecretManagerAuthAWS.CreateSecret(ctx, &si)
	s.Require().Nil(err)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestFetchSecret() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID + 1
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceURL:        s.Tc.AwsLamdbaTestURL,
		ServiceAuth:       hestia_req_types.ServiceAuthConfig{},
	}
	si := aws_secrets.SecretInfo{
		Region: "us-west-1",
		Name:   fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID),
		Key:    fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID),
	}
	so, err := GetServiceRoutesAuths(ctx, si)
	s.Require().Nil(err)
	s.Require().NotEmpty(so)
}

func (s *ArtemisHydraSecretsManagerTestSuite) TestFetchServiceRoutesAuths() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID + 1
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	v := hestia_req_types.ServiceRequestWrapper{
		GroupName:         "testGroup",
		ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
		ServiceURL:        s.Tc.AwsLamdbaTestURL,
		ServiceAuth:       hestia_req_types.ServiceAuthConfig{},
	}
	si := aws_secrets.SecretInfo{
		Region: "us-west-1",
		Name:   fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID),
		Key:    fmt.Sprintf("%s-%d-%d", v.GroupName, ou.OrgID, v.ProtocolNetworkID),
	}
	srw, err := GetServiceRoutesAuths(ctx, si)
	s.Require().Nil(err)
	s.Require().NotEmpty(srw)

	s.Assert().Equal(v.GroupName, srw.GroupName)
	s.Assert().Equal(v.ServiceURL, srw.ServiceURL)
	s.Assert().Equal(v.ServiceAuth, srw.ServiceAuth)

}

func TestArtemisHydraSecretsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ArtemisHydraSecretsManagerTestSuite))
}
