package hestia_eks_aws

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

type AwsPricingClientTestSuite struct {
	test_suites_base.TestSuite
	pc AwsPricing
}

func (s *AwsPricingClientTestSuite) SetupTest() {
	s.InitLocalConfigs()
	eksCreds := EksCredentials{
		Region:       UsWest1,
		AccessKey:    s.Tc.AwsAccessKeyEks,
		AccessSecret: s.Tc.AwsSecretKeyEks,
	}
	p, err := InitPricingClient(ctx, eksCreds)
	s.Require().NoError(err)
	s.pc = p
	s.Require().NotNil(s.pc.Client)
}

func (s *AwsPricingClientTestSuite) TestGetPrice() {

}

func TestAwsPricingClientTestSuite(t *testing.T) {
	suite.Run(t, new(AwsPricingClientTestSuite))
}
