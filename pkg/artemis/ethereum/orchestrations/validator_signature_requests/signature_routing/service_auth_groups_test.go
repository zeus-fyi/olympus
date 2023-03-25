package artemis_validator_signature_service_routing

import (
	"time"

	artemis_validator_service_groups_models "github.com/zeus-fyi/olympus/datastores/postgres/apps/artemis/models"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchAndSetServiceGroupsAuths() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	s.InitLocalConfigs()
	s.Pg.InitPG(ctx, s.Tc.ProdLocalDbPgconn)

	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking",
		Env:           "production",
	}
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID:     hestia_req_types.EthereumEphemeryProtocolNetworkID,
		OrgID:                 ou.OrgID,
		ValidatorClientNumber: 0,
	}
	vsMetadata, err := GetServiceMetadata(ctx, vsi, cctx)
	s.Require().Nil(err)

	err = FetchAndSetServiceGroupsAuths(ctx, vsMetadata)
	s.Require().Nil(err)

	//expAuth := hestia_req_types.ServiceRequestWrapper{
	//	GroupName:         "testGroup",
	//	ProtocolNetworkID: hestia_req_types.EthereumEphemeryProtocolNetworkID,
	//	ServiceAuth: hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
	//		ServiceURL: s.Tc.AwsLamdbaTestURL,
	//		SecretName: "testLambdaExternalSecret",
	//		AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
	//		SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
	//	}},
	//}
	//auth, err := GetGroupAuthFromInMemFS(ctx, expAuth.GroupName)
	//s.Require().Nil(err)
	//s.Assert().Equal(expAuth.ServiceAuth, auth)

	time.Sleep(5 * time.Second)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchServiceAuthRouteGrouping() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID
	cctx := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "ephemeral-staking", // set with your own namespace
		Env:           "production",
	}
	vsi := artemis_validator_service_groups_models.ValidatorServiceCloudCtxNsProtocol{
		ProtocolNetworkID:     hestia_req_types.EthereumEphemeryProtocolNetworkID,
		OrgID:                 ou.OrgID,
		ValidatorClientNumber: 0,
	}
	err := GetServiceAuthAndURLs(ctx, vsi, cctx)
	s.Require().Nil(err)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestSetAndGetGroupAuthToInMemFS() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	serviceAuth := hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
		SecretName: "testLambdaExternalSecret",
		AccessKey:  s.Tc.AwsAccessKeyLambdaExt,
		SecretKey:  s.Tc.AwsSecretKeyLambdaExt,
	}}
	groupName := "group1"
	err := SetGroupAuthInMemFS(ctx, groupName, serviceAuth)
	s.Require().Nil(err)

	auth, err := GetGroupAuthFromInMemFS(ctx, "group1")
	s.Require().Nil(err)
	s.Assert().Equal(serviceAuth, auth)
}
