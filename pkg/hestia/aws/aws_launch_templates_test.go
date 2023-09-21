package hestia_eks_aws

func (s *AwsEKSTestSuite) TestCreateInstanceTemplate() {
	instanceTypes := []string{
		"i3.4xlarge",
		"i3.8xlarge",
		"i4i.4xlarge",
	}

	for _, instanceType := range instanceTypes {
		template := CreateNvmeLaunchTemplate(instanceType)
		s.Require().NotNil(template)
	}
}
