package artemis_validator_signature_service_routing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
	"github.com/zeus-fyi/zeus/pkg/zeus/client/zeus_common_types"
)

type ValidatorServiceAuthRoutesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ValidatorServiceAuthRoutesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	auth := aws_secrets.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestServiceGroupingHelper() {
	srs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}

	keyOne := "0x1"
	srs.Map[keyOne] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "one"}
	srs.Map["0x2"] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "two"}
	srs.Map["0x3"] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "one"}
	InitRouteMapInMemFS(ctx)
	p := filepaths.Path{
		DirIn: ".",
		FnIn:  keyOne,
	}
	err := RouteMapInMemFS.MakeFileIn(&p, []byte("https://fake-service.com"))
	s.Require().Nil(err)

	resp := GroupSigRequestsByServiceURL(ctx, srs)
	s.Require().NotEmpty(resp)
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
	svc, err := GetServiceURLs(ctx, cctx)
	s.Require().Nil(err)
	s.Require().NotEmpty(svc)
}

func TestValidatorServiceAuthRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceAuthRoutesTestSuite))
}
