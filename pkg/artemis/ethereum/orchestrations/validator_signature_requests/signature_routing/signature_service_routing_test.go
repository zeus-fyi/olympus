package artemis_validator_signature_service_routing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	hestia_test "github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/test"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aws_secrets "github.com/zeus-fyi/zeus/pkg/aegis/aws"
	aegis_inmemdbs "github.com/zeus-fyi/zeus/pkg/aegis/inmemdbs"
)

type ValidatorServiceAuthRoutesTestSuite struct {
	hestia_test.BaseHestiaTestSuite
}

var ctx = context.Background()

func (s *ValidatorServiceAuthRoutesTestSuite) SetupTest() {
	s.InitLocalConfigs()
	err := InitRouteMapInMemFS(ctx)
	s.Require().Nil(err)
	auth := aws_secrets.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
}

func (s *ValidatorServiceAuthRoutesTestSuite) TestServiceGroupingHelper() {
	srs := aegis_inmemdbs.EthereumBLSKeySignatureRequests{Map: make(map[string]aegis_inmemdbs.EthereumBLSKeySignatureRequest)}
	keyOne, keyTwo, keyThree := "0x1", "0x2", "0x3"
	keyOneGroupName, keyTwoGroupName, keyThreeGroupName := "groupOne", "groupOne", "groupThree"
	srs.Map[keyOne] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "one"}
	srs.Map[keyTwo] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "two"}
	srs.Map[keyThree] = aegis_inmemdbs.EthereumBLSKeySignatureRequest{Message: "three"}
	err := SetPubkeyToGroupService(ctx, keyOne, keyOneGroupName)
	s.Require().Nil(err)
	err = SetPubkeyToGroupService(ctx, keyTwo, keyTwoGroupName)
	s.Require().Nil(err)
	err = SetPubkeyToGroupService(ctx, keyThree, keyThreeGroupName)
	s.Require().Nil(err)

	resp := GroupSigRequestsByGroupName(ctx, srs)
	s.Require().NotEmpty(resp)
	s.Require().Equal(2, len(resp))
	for k, v := range resp {
		if k == keyOneGroupName {
			s.Require().Equal(2, len(v.Map))
			s.Require().Equal("one", v.Map[keyOne].Message)
			s.Require().Equal("two", v.Map[keyTwo].Message)
		}
		if k == keyThreeGroupName {
			s.Require().Equal(1, len(v.Map))
			s.Require().Equal("three", v.Map[keyThree].Message)
		}
	}
}

func TestValidatorServiceAuthRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorServiceAuthRoutesTestSuite))
}
