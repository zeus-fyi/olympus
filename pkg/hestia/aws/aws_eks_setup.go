package hestia_eks_aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/rs/zerolog/log"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

type EksCredentials struct {
	Creds       aegis_aws_auth.AuthAWS `json:"creds"`
	ClusterName string                 `json:"clusterName"`
}

func GetEksKubeConfig(ctx context.Context, eksCreds EksCredentials) (*AwsEKS, *KubeConfig, error) {
	eka, err := InitAwsEKS(ctx, eksCreds.Creds)
	if err != nil {
		log.Err(err).Msg("GetKubeConfig: failed to init EKS client")
		return nil, nil, err
	}
	clusterInput := &eks.DescribeClusterInput{
		Name: aws.String(eksCreds.ClusterName),
	}
	clusterOutput, err := eka.DescribeCluster(ctx, clusterInput)
	if err != nil {
		log.Err(err).Msg("GetKubeConfig: failed to describe cluster")
		return nil, nil, err
	}
	if clusterOutput == nil || clusterOutput.Cluster == nil || clusterOutput.Cluster.Endpoint == nil || clusterOutput.Cluster.CertificateAuthority == nil || clusterOutput.Cluster.CertificateAuthority.Data == nil {
		err = fmt.Errorf("GetKubeConfig: clusterOutput is nil")
		log.Err(err).Msg("GetKubeConfig: clusterOutput is nil")
		return nil, nil, err
	}

	return &eka, populateEksKubeConfig(eksCreds.ClusterName, clusterOutput, eksCreds.Creds.Region), nil
}

func populateEksKubeConfig(clusterName string, clusterOutput *eks.DescribeClusterOutput, region string) *KubeConfig {
	kubeConfig := KubeConfig{
		APIVersion: "v1",
		Kind:       "Config",
		EksKubeInfo: &EksKubeInfo{
			Arn:                clusterOutput.Cluster.Arn,
			RoleArn:            clusterOutput.Cluster.RoleArn,
			ResourcesVpcConfig: clusterOutput.Cluster.ResourcesVpcConfig,
		},
		Clusters: []ClusterEntry{
			{
				Name: clusterName,
				Cluster: ClusterInfo{
					Server:                   *clusterOutput.Cluster.Endpoint,
					CertificateAuthorityData: *clusterOutput.Cluster.CertificateAuthority.Data,
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
						Args:       []string{"eks", "get-token", "--cluster-name", clusterName, "--region", region},
					},
				},
			},
		},
	}

	return &kubeConfig
}
