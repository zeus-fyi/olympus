package hestia_eks_aws

import "fmt"

func (s *AwsEKSTestSuite) TestCreateInstanceTemplate() {
	instanceTypes := []string{
		"i3.4xlarge",
		//"i3.8xlarge",
		//"i4i.4xlarge",
	}

	for _, instanceType := range instanceTypes {
		template := CreateNvmeLaunchTemplate(instanceType)
		s.Require().NotNil(template)
		lto, err := s.ecc.RegisterInstanceTemplate(instanceType)
		s.Require().NoError(err)
		s.Require().NotNil(lto)
		fmt.Println(lto.LaunchTemplate.LaunchTemplateId)
		fmt.Println(lto.LaunchTemplate.LaunchTemplateName)
	}
}
