package hestia_eks_aws

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/rs/zerolog/log"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type AwsEc2 struct {
	*ec2.Client
}

func InitAwsEc2(ctx context.Context, accessCred aegis_aws_auth.AuthAWS) (AwsEc2, error) {
	creds := credentials.NewStaticCredentialsProvider(accessCred.AccessKey, accessCred.SecretKey, "")
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(creds),
		config.WithRegion(accessCred.Region),
	)
	if err != nil {
		return AwsEc2{}, err
	}
	return AwsEc2{ec2.NewFromConfig(cfg)}, nil
}

func GetLaunchTemplate(slug string) *eksTypes.LaunchTemplateSpecification {
	// todo
	lt := &eksTypes.LaunchTemplateSpecification{
		Id:      nil,
		Name:    nil,
		Version: nil,
	}
	return lt
}

func GetLaunchTemplateName(slug string) string {
	return fmt.Sprintf("eks-pv-raid-launch-template-%s", slug)
}

func (a *AwsEc2) RegisterInstanceTemplate(slug string) (*ec2.CreateLaunchTemplateOutput, error) {
	lti := CreateNvmeLaunchTemplate(slug)
	launchTemplateOutput, err := a.CreateLaunchTemplate(context.Background(), lti)
	if err != nil {
		log.Err(err).Interface("lto", launchTemplateOutput).Msg("failed to create launch template")
		return launchTemplateOutput, err
	}
	return launchTemplateOutput, err
}

func CreateNvmeLaunchTemplate(slug string) *ec2.CreateLaunchTemplateInput {
	// Create EC2 Launch Template with User Data
	userData := `#!/bin/bash
		# Install NVMe CLI
        yum install nvme-cli -y
        
        # Get list of NVMe Drives
        nvme_drives=$(nvme list | grep "Amazon EC2 NVMe Instance Storage" | cut -d " " -f 1 || true)
        readarray -t nvme_drives <<< "$nvme_drives"
        num_drives=${#nvme_drives[@]}
        
        # Install software RAID utility
        yum install mdadm -y
        
        # Create RAID-0 array across the instance store NVMe SSDs
        mdadm --create /dev/md0 --level=0 --name=md0 --raid-devices=$num_drives "${nvme_drives[@]}"

        # Format drive with Ext4
        mkfs.ext4 /dev/md0

        # Get RAID array's UUID
        uuid=$(blkid -o value -s UUID /dev/md0)
        
        # Create a filesystem path to mount the disk
        mount_location="/mnt/fast-disks/${uuid}"
        mkdir -p $mount_location
        
        # Mount RAID device
        mount /dev/md0 $mount_location
        
        # Have disk be mounted on reboot
        mdadm --detail --scan >> /etc/mdadm.conf 
        echo /dev/md0 $mount_location ext4 defaults,noatime 0 2 >> /etc/fstab
	`
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	lt := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String(fmt.Sprintf("eks-pv-raid-launch-template-%s", slug)),
		VersionDescription: aws.String("eks nvme bootstrap"),
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			UserData:     aws.String(encodedUserData),
			InstanceType: types.InstanceType(slug),
		},
	}
	return lt
}
