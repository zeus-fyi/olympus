package hestia_cluster_configs

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/authorized_clusters"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/bases/org_users"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	hestia_eks_aws "github.com/zeus-fyi/olympus/pkg/hestia/aws"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/compression"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/test_utils/test_suites/test_suites_base"
	zeus_core "github.com/zeus-fyi/olympus/pkg/zeus/core"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
)

var ctx = context.Background()

type ExtClusterCfgsTestSuite struct {
	test_suites_base.TestSuite
}

func (s *ExtClusterCfgsTestSuite) SetupTest() {
	s.InitLocalConfigs()
	home, exists := os.LookupEnv("HOME")
	if exists {
		aws_secrets.CredBasePath = path.Join(home, aws_secrets.CredBasePath)
		aws_secrets.ConfigPath = path.Join(home, aws_secrets.ConfigPath)
	}
}

func (s *ExtClusterCfgsTestSuite) TestGetPlatformServiceAccountsExplore() {
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)

	aiUserOrgID := 1685378241971196000
	ou := org_users.NewOrgUserWithID(aiUserOrgID, aiUserOrgID)
	ps, perr := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	s.Require().Nil(perr)
	s.Require().NotNil(ps)

	for clusterName, creds := range ps.AwsEksServiceMap {

		if clusterName == "zeus-eks-us-east-2" {
			continue
		}
		eksCredsAuth := hestia_eks_aws.EksCredentials{
			Creds:       creds,
			ClusterName: clusterName,
			Ou:          ou,
		}

		ek, _, err := hestia_eks_aws.GetEksKubeConfig(ctx, eksCredsAuth)
		s.Require().NoError(err)
		s.Require().NotNil(ek)
	}
}

func (s *ExtClusterCfgsTestSuite) TestGetPlatformServiceAccountsToExtClusterCfgs() {
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)

	aiUserOrgID := 1685378241971196000
	ou := org_users.NewOrgUserWithID(aiUserOrgID, aiUserOrgID)
	ps, perr := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	s.Require().Nil(perr)
	s.Require().NotNil(ps)

	var extClusterConfigs []authorized_clusters.K8sClusterConfig
	for clusterName, creds := range ps.AwsEksServiceMap {

		if clusterName == "zeus-eks-us-east-2" {
			continue
		}
		eksCredsAuth := hestia_eks_aws.EksCredentials{
			Creds:       creds,
			ClusterName: clusterName,
			Ou:          ou,
		}
		ek, kubeConfig, err := hestia_eks_aws.GetEksKubeConfig(ctx, eksCredsAuth)
		s.Require().NoError(err)
		s.Require().NotNil(ek)

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
		inCmp, err := compression.GzipDirectoryToMemoryFS(p, inMemFilestore)
		s.Require().Nil(err)
		s.Require().NotNil(inCmp)

		k := zeus_core.K8Util{}
		k.ConnectToK8sFromInMemFsCfgPath(inMemFilestore)

		ctxes, err := k.GetContexts()
		s.Require().Nil(err)
		s.Require().NotNil(ctxes)
		for name, _ := range ctxes {
			fmt.Println(name)

			kctx := zeus_common_types.CloudCtxNs{
				CloudProvider: "aws",
				Region:        "us-east-1",
				Namespace:     "kube-system",
				Context:       name,
			}
			nses, nerr := k.GetNamespaces(ctx, kctx)
			s.Require().Nil(nerr)
			s.Require().NotNil(nses)

			for _, ns := range nses.Items {
				fmt.Println(ns.Name)
			}

			ec := authorized_clusters.K8sClusterConfig{
				CloudCtxNs:   kctx,
				ContextAlias: clusterName,
				IsActive:     false,
			}
			extClusterConfigs = append(extClusterConfigs, ec)
		}
	}
	for _, ec := range extClusterConfigs {
		fmt.Println(ec)
	}
}

func TestExtClusterCfgsTestSuite(t *testing.T) {
	suite.Run(t, new(ExtClusterCfgsTestSuite))
}

//cms, cerr := k.GetConfigMapWithKns(ctx, kctx, "aws-auth", nil)
//s.Require().Nil(cerr)
//s.Require().NotNil(cms)
//var awsAuthMapRoles AwsAuthConfigMap
//err = yaml.Unmarshal([]byte(cms.Data["mapRoles"]), &awsAuthMapRoles.MapRoles)
//s.Require().Nil(err)
//
//eka, eerr := hestia_eks_aws.InitAwsEKS(ctx, eksCredsAuth.Creds)
//s.Require().Nil(eerr)
//s.Require().NotNil(eka)
//
//awsAuthMapRoles.MapUsers = []UserEntry{
//	{
//		UserARN:  aws.StringValue(eka.Arn),
//		Username: eka.Username,
//		Groups:   []string{"system:masters"},
//	},
//}
//
//b, berr := yaml.Marshal(awsAuthMapRoles.MapUsers)
//s.Require().Nil(berr)
//
//cms.Data["mapUsers"] = string(b)
//sc := &v1.StorageClass{
//	ObjectMeta: metav1.ObjectMeta{
//		Name: "aws-ebs-gp3-max-performance", // Provide a meaningful name for the StorageClass
//	},
//	Provisioner: "ebs.csi.aws.com", // AWS EBS CSI driver
//	Parameters: map[string]string{
//		"type":       "gp3",   // Specify gp3 type for the EBS volume
//		"iops":       "16000", // Maximum IOPS for gp3
//		"throughput": "1000",  // Maximum throughput in MB/s for gp3
//		//"encrypted":  "true",  // Optionally, ensure encryption is enabled
//		// "fsType":      "ext4",            // Specify filesystem type if needed, e.g., ext4 or xfs
//	},
//	ReclaimPolicy:        nil,                   // You can specify a ReclaimPolicy if needed
//	AllowVolumeExpansion: pointer.BoolPtr(true), // Optionally allow volume expansion
//}
//_, kerr := k.CreateStorageClass(ctx, kctx, sc)
//s.Require().Nil(kerr)
//
//cms2, cerr := k.UpdateConfigMapWithKns(ctx, kctx, cms, nil)
//s.Require().Nil(cerr)
//s.Require().Equal(cms.Data, cms2.Data)
