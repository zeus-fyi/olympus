package hestia_eks_aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type KubeConfig struct {
	APIVersion     string         `json:"apiVersion"`
	Kind           string         `json:"kind"`
	Clusters       []ClusterEntry `json:"clusters"`
	Contexts       []ContextEntry `json:"contexts"`
	CurrentContext string         `json:"current-context"`
	Users          []UserEntry    `json:"users"`

	EksKubeInfo *EksKubeInfo `json:"eksKubeInfo,omitempty"`
}

func (k *KubeConfig) GetEksSubnets() ([]string, error) {
	if k.EksKubeInfo == nil || k.EksKubeInfo.ResourcesVpcConfig == nil {
		return nil, fmt.Errorf("GetEksSubnets: EksKubeInfo or ResourcesVpcConfig is nil")
	}
	return k.EksKubeInfo.ResourcesVpcConfig.SubnetIds, nil
}
func (k *KubeConfig) GetEksRoleArn() (*string, error) {
	if k.EksKubeInfo == nil || k.EksKubeInfo.RoleArn == nil {
		return nil, fmt.Errorf("GetEksRoleArn: EksKubeInfo or RoleArn is nil")
	}
	return k.EksKubeInfo.RoleArn, nil
}

type EksKubeInfo struct {
	Arn                *string                  `json:"arn,omitempty"`
	RoleArn            *string                  `json:"roleArn,omitempty"`
	ResourcesVpcConfig *types.VpcConfigResponse `json:"resourcesVpcConfig,omitempty"`
}

type ClusterEntry struct {
	Name    string      `json:"name"`
	Cluster ClusterInfo `json:"cluster"`
}

type ClusterInfo struct {
	Server                   string `json:"server"`
	CertificateAuthorityData string `json:"certificate-authority-data"`
}

type ContextEntry struct {
	Name    string      `json:"name"`
	Context ContextInfo `json:"context"`
}

type ContextInfo struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

type UserEntry struct {
	Name string   `json:"name"`
	User UserInfo `json:"user"`
}

type UserInfo struct {
	Exec ExecConfig `json:"exec"`
}

type ExecConfig struct {
	APIVersion string   `json:"apiVersion"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
}
