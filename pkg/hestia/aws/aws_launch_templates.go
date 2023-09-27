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

var (
	SlugToInstanceID = map[string]string{
		"i3.4xlarge":  "lt-0e97987981123738a",
		"i3.8xlarge":  "lt-0a2a4a58163737e16",
		"i4i.4xlarge": "lt-0872c2f0aff238f9a",
	}

	SlugToInstanceTemplateName = map[string]string{
		"i3.4xlarge":  "eks-pv-raid-launch-template-i3.4xlarge",
		"i3.8xlarge":  "eks-pv-raid-launch-template-i3.8xlarge",
		"i4i.4xlarge": "eks-pv-raid-launch-template-i4i.4xlarge",
	}
)

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

func GetLaunchTemplate(id, instanceName string) *eksTypes.LaunchTemplateSpecification {
	lt := &eksTypes.LaunchTemplateSpecification{
		Id: aws.String(id),
		//Name:    aws.String(instanceName),
		Version: nil,
	}
	return lt
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

func (a *AwsEc2) UpdateInstanceTemplate(templateID string) (*ec2.ModifyLaunchTemplateOutput, error) {
	mti := &ec2.ModifyLaunchTemplateInput{
		ClientToken:        nil,
		DefaultVersion:     nil,
		DryRun:             nil,
		LaunchTemplateId:   aws.String(templateID),
		LaunchTemplateName: nil,
	}

	launchTemplateOutput, err := a.ModifyLaunchTemplate(context.Background(), mti)
	if err != nil {
		log.Err(err).Interface("lto", launchTemplateOutput).Msg("failed to create launch template")
		return launchTemplateOutput, err
	}
	return launchTemplateOutput, err
}

// CreateNvmeLaunchTemplate to read how retarded this is refer to: https://docs.aws.amazon.com/eks/latest/userguide/launch-templates.html
func CreateNvmeLaunchTemplate(slug string) *ec2.CreateLaunchTemplateInput {
	// Create EC2 Launch Template with User Data
	userData := `MIME-Version: 1.0
Content-Type: multipart/mixed; boundary=3bfcfdaa6a583f7487ad3c90c4853b64cdf474e9b43f0b14b706e557e015

--3bfcfdaa6a583f7487ad3c90c4853b64cdf474e9b43f0b14b706e557e015
Content-Type: text/x-shellscript
Content-Type: charset="us-ascii"

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

--3bfcfdaa6a583f7487ad3c90c4853b64cdf474e9b43f0b14b706e557e015--`

	ebs := &types.LaunchTemplateEbsBlockDeviceRequest{
		Encrypted:           aws.Bool(false), // Modify this as per your requirement
		DeleteOnTermination: aws.Bool(true),
		Iops:                aws.Int32(3000),
		VolumeSize:          aws.Int32(20),
		VolumeType:          types.VolumeTypeGp3,
		Throughput:          aws.Int32(125),
	}
	blockDeviceMapping := types.LaunchTemplateBlockDeviceMappingRequest{
		DeviceName: aws.String("/dev/xvda"),
		Ebs:        ebs,
	}
	encodedUserData := base64.StdEncoding.EncodeToString([]byte(userData))
	lt := &ec2.CreateLaunchTemplateInput{
		LaunchTemplateName: aws.String(fmt.Sprintf("eks-pv-raid-launch-template-%s", slug)),
		VersionDescription: aws.String("eks nvme bootstrap"),
		LaunchTemplateData: &types.RequestLaunchTemplateData{
			BlockDeviceMappings: []types.LaunchTemplateBlockDeviceMappingRequest{blockDeviceMapping},
			InstanceType:        types.InstanceType(slug),
			SecurityGroupIds:    []string{AwsUsWestSecurityGroupID},
			UserData:            aws.String(encodedUserData),
		},
	}
	return lt
}
