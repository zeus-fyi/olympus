package hestia_eks_aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

var ctx = context.Background()

type AwsEKSTestSuite struct {
	test_suites_base.TestSuite
	ek AwsEKS
}

func (s *AwsEKSTestSuite) SetupTest() {
	s.InitLocalConfigs()
	eksCreds := aegis_aws_auth.AuthAWS{
		Region:    UsWest1,
		AccessKey: s.Tc.AwsAccessKeyEks,
		SecretKey: s.Tc.AwsSecretKeyEks,
	}
	eka, err := InitAwsEKS(ctx, eksCreds)
	s.Require().NoError(err)
	s.ek = eka
	s.Require().NotNil(s.ek.Client)
}

func (s *AwsEKSTestSuite) TestCreateNodeGroup() {
	machineType := "t2.micro"
	nodeGroupName := "test-node-group"

	//awsTaint := types.Taint{
	//	Effect: "NoSchedule",
	//	Key:    aws.String("org"),
	//	Value:  aws.String("org"),
	//}

	params := &eks.CreateNodegroupInput{
		ClusterName:        aws.String(AwsUsWest1Context),
		NodeRole:           aws.String(AwsEksRole),
		NodegroupName:      aws.String(nodeGroupName),
		AmiType:            types.AMITypesAl2X8664,
		Subnets:            UsWestSubnetIDs,
		CapacityType:       "",
		ClientRequestToken: nil,
		InstanceTypes:      []string{machineType},
		Labels:             nil,
		ScalingConfig: &types.NodegroupScalingConfig{
			DesiredSize: aws.Int32(1),
			MaxSize:     aws.Int32(1),
			MinSize:     aws.Int32(1),
		},
		Taints: nil,
	}
	ngs, err := s.ek.AddNodeGroup(ctx, params)
	s.Require().Nil(err)
	s.Require().NotNil(ngs)
}

func (s *AwsEKSTestSuite) TestRemoveNodeGroup() {
	nodeGroupName := "test-node-group"
	params := &eks.DeleteNodegroupInput{
		ClusterName:   aws.String(AwsUsWest1Context),
		NodegroupName: aws.String(nodeGroupName),
	}
	ngs, err := s.ek.RemoveNodeGroup(ctx, params)
	s.Require().Nil(err)
	s.Require().NotNil(ngs)
}

func (s *AwsEKSTestSuite) TestListNodeGroups() {
	results := int32(10)
	params := &eks.ListNodegroupsInput{
		ClusterName: aws.String(AwsUsWest1Context),
		MaxResults:  &results,
	}
	ngs, err := s.ek.ListNodegroups(ctx, params)
	s.Require().Nil(err)
	s.Require().NotNil(ngs)
}

func TestAwsEKSTestSuite(t *testing.T) {
	suite.Run(t, new(AwsEKSTestSuite))
}
