package hestia_eks_aws

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/ghodss/yaml"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var ctx = context.Background()

type AwsEKSTestSuite struct {
	test_suites_base.TestSuite
	ek  AwsEKS
	ecc AwsEc2
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

	ecc, err := InitAwsEc2(ctx, eksCreds)
	s.Require().NoError(err)
	s.ecc = ecc
}

func (s *AwsEKSTestSuite) TestGetKubeConfig() {
	eksCreds := aegis_aws_auth.AuthAWS{
		Region:    "us-east-2",
		AccessKey: s.Tc.AwsZeusEksServiceAccessKey,
		SecretKey: s.Tc.AwsZeusEksServiceSecretKey,
	}
	eka, err := InitAwsEKS(ctx, eksCreds)
	s.Require().NoError(err)
	clusterName := "zeus-eks-us-east-2"
	// Retrieve cluster details
	clusterInput := &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}
	clusterOutput, err := eka.DescribeCluster(ctx, clusterInput)
	s.Require().Nil(err)
	s.Require().NotNil(clusterOutput)

	clusterEndpoint := *clusterOutput.Cluster.Endpoint
	clusterCA := *clusterOutput.Cluster.CertificateAuthority.Data
	kubeConfig := KubeConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []ClusterEntry{
			{
				Name: clusterName,
				Cluster: ClusterInfo{
					Server:                   clusterEndpoint,
					CertificateAuthorityData: clusterCA,
				},
			},
		},
		Contexts: []ContextEntry{
			{
				Name: clusterName,
				Context: ContextInfo{
					Cluster: clusterName,
					User:    clusterName,
				},
			},
		},
		CurrentContext: clusterName,
		Users: []UserEntry{
			{
				Name: clusterName,
				User: UserInfo{
					Exec: ExecConfig{
						APIVersion: "client.authentication.k8s.io/v1beta1",
						Command:    "aws",
						Args:       []string{"eks", "get-token", "--cluster-name", clusterName},
					},
				},
			},
		},
	}

	kubeConfigYAML, err := yaml.Marshal(&kubeConfig)
	s.Require().Nil(err)

	p := filepaths.Path{
		PackageName: "",
		DirIn:       "/.kube",
		FnIn:        "config",
	}

	inMemFilestore := memfs.NewMemFs()
	err = inMemFilestore.MakeFileIn(&p, kubeConfigYAML)
	s.Require().Nil(err)

	k := zeus_core.K8Util{}
	k.ConnectToK8sFromInMemFsCfgPath(inMemFilestore)

	ctxes, err := k.GetContexts()
	s.Require().Nil(err)
	s.Require().NotNil(ctxes)
	for name, _ := range ctxes {
		fmt.Println(name)

		kctx := zeus_common_types.CloudCtxNs{
			CloudProvider: "aws",
			Region:        "us-east-2",
			Context:       name,
		}
		nses, nerr := k.GetNamespaces(ctx, kctx)
		s.Require().Nil(nerr)
		s.Require().NotNil(nses)

		for _, ns := range nses.Items {
			fmt.Println(ns.Name)
		}
	}

	// Write the kubeconfig to a file

	//pathPrefix := "/Users/alex/go/Olympus/olympus/pkg/hestia/aws/"
	//err = os.WriteFile(pathPrefix+"kubeconfig.yaml", kubeConfigYAML, 0600)
	//s.Require().Nil(err)
	//
	//fmt.Println("Kubeconfig successfully written to kubeconfig.yaml")
}

func (s *AwsEKSTestSuite) TestCreateNodeGroup() {
	machineType := "t2.micro"
	nodeGroupName := "test-node-group"

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

func (s *AwsEKSTestSuite) TestCreateNvmeNodeGroup() {
	nodeGroupName := "test-" + uuid.New().String()
	ltID := SlugToLaunchTemplateID["i3.8xlarge"]
	labels := make(map[string]string)
	labels = AddAwsEksNvmeLabels(labels)

	oid := 7138983863666903883
	orgTaint := types.Taint{
		Effect: "NO_SCHEDULE",
		Key:    aws.String(fmt.Sprintf("org-%d", oid)),
		Value:  aws.String(fmt.Sprintf("org-%d", oid)),
	}

	params := &eks.CreateNodegroupInput{
		ClusterName:   aws.String(AwsUsWest1Context),
		NodeRole:      aws.String(AwsEksRole),
		NodegroupName: aws.String(nodeGroupName),
		AmiType:       types.AMITypesAl2X8664,
		Subnets:       UsWestSubnetIDs,
		LaunchTemplate: &types.LaunchTemplateSpecification{
			Id: aws.String(ltID),
		},
		ClientRequestToken: aws.String(nodeGroupName),
		ScalingConfig: &types.NodegroupScalingConfig{
			DesiredSize: aws.Int32(1),
			MaxSize:     aws.Int32(1),
			MinSize:     aws.Int32(1),
		},
		Labels: labels,
		Taints: []types.Taint{orgTaint},
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
