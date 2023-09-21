package hestia_eks_aws

func (s *AwsEKSTestSuite) TestCreateInstanceTemplate() {
	template := CreateNvmeLaunchTemplate("t2.micro")
	s.Require().NotNil(template)
}
