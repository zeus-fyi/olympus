package hestia_cluster_configs

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
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

}
func (s *ExtClusterCfgsTestSuite) TestGetPlatformServiceAccountsToExtClusterCfgs() {
	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)

	aiUserOrgID := 1699642242976434000
	ou := org_users.NewOrgUserWithID(aiUserOrgID, aiUserOrgID)
	ps, perr := aws_secrets.GetServiceAccountSecrets(ctx, ou)
	s.Require().Nil(perr)
	s.Require().NotNil(ps)

	var extClusterConfigs []authorized_clusters.K8sClusterConfig
	for clusterName, creds := range ps.AwsEksServiceMap {
		eksCredsAuth := hestia_eks_aws.EksCredentials{
			Creds:       creds,
			ClusterName: clusterName,
		}
		_, kubeConfig, err := hestia_eks_aws.GetEksKubeConfig(ctx, eksCredsAuth)
		s.Require().NoError(err)

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
				Region:        "us-east-2",
				Namespace:     "kube-system",
				Context:       name,
			}
			nses, nerr := k.GetNamespaces(ctx, kctx)
			s.Require().Nil(nerr)
			s.Require().NotNil(nses)

			for _, ns := range nses.Items {
				fmt.Println(ns.Name)
			}

			cms, cerr := k.GetConfigMapWithKns(ctx, kctx, "aws-auth", nil)
			s.Require().Nil(cerr)
			s.Require().NotNil(cms)
			var awsAuthMapRoles AwsAuthConfigMap
			err = yaml.Unmarshal([]byte(cms.Data["mapRoles"]), &awsAuthMapRoles.MapRoles)
			s.Require().Nil(err)

			eka, eerr := hestia_eks_aws.InitAwsEKS(ctx, eksCredsAuth.Creds)
			s.Require().Nil(eerr)
			s.Require().NotNil(eka)

			awsAuthMapRoles.MapUsers = []UserEntry{
				{
					UserARN:  aws.StringValue(eka.Arn),
					Username: eka.Username,
					Groups:   []string{"system:masters"},
				},
			}

			b, berr := yaml.Marshal(awsAuthMapRoles.MapUsers)
			s.Require().Nil(berr)

			cms.Data["mapUsers"] = string(b)

			cms2, cerr := k.UpdateConfigMapWithKns(ctx, kctx, cms, nil)
			s.Require().Nil(cerr)
			s.Require().Equal(cms.Data, cms2.Data)
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
