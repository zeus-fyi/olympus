package artemis_validator_signature_service_routing

import (
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_req_types "github.com/zeus-fyi/zeus/pkg/hestia/client/req_types"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

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
	err := GetServiceAuthAndURLs(ctx, cctx)
	s.Require().Nil(err)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestFetchGroupAuths() {
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
	err := GetServiceAuthAndURLs(ctx, cctx)
	s.Require().Nil(err)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestSetAndGetGroupAuthToInMemFS() {
	ou := org_users.OrgUser{}
	ou.OrgID = s.Tc.ProductionLocalTemporalOrgID
	ou.UserID = s.Tc.ProductionLocalTemporalUserID

	serviceAuth := hestia_req_types.ServiceAuthConfig{AuthLamdbaAWS: &hestia_req_types.AuthLamdbaAWS{
		SecretName:   "testLambdaExternalSecret",
		AccessKey:    s.Tc.AwsAccessKeyLambdaExt,
		AccessSecret: s.Tc.AwsSecretKeyLambdaExt,
	}}
	groupName := "group1"
	err := SetGroupAuthInMemFS(ctx, groupName, serviceAuth)
	s.Require().Nil(err)

	auth, err := GetGroupAuthFromInMemFS(ctx, "group1")
	s.Require().Nil(err)
	s.Assert().Equal(serviceAuth, auth)
}
