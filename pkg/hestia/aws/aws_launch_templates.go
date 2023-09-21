package hestia_eks_aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func CreateNvmeLaunchTemplate(instanceType string) *ec2.CreateLaunchTemplateInput {
	// Create EC2 Launch Template with User Data
	userData := `#!/bin/bash
	# Your pre-bootstrap commands here
	# ...
	`
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	lt := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String("eks-pv-raid-launch-template"),
		VersionDescription: aws.String("eks nvme bootstrap"),
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			UserData:     aws.String(encodedUserData),
			InstanceType: types.InstanceType(instanceType),
		},
	}
	return lt
}
