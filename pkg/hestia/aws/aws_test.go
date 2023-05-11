package hestia_eks_aws

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
)

var ctx = context.Background()

type AwsEKSTestSuite struct {
	test_suites_base.TestSuite
	ek AwsEKS
}

func (s *AwsEKSTestSuite) SetupTest() {
	s.InitLocalConfigs()
	eksCreds := EksCredentials{
		Region:       UsWest1,
		AccessKey:    s.Tc.AwsAccessKeyEks,
		AccessSecret: s.Tc.AwsSecretKeyEks,
	}
	eka, err := InitAwsEKS(ctx, eksCreds)
	s.Require().NoError(err)
	s.ek = eka
	s.Require().NotNil(s.ek.Client)
}

func (s *AwsEKSTestSuite) TestListSizes() {

}
