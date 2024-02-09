package zeus_core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/aegis/aws_secrets"
	artemis_hydra_orchestrations_aws_auth "github.com/zeus-fyi/olympus/pkg/artemis/ethereum/orchestrations/validator_signature_requests/aws_auth"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretsTestSuite struct {
	K8TestSuite
}

func (s *SecretsTestSuite) TestGetSecrets() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ethereum"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	//err = s.K.DeleteSecretWithKns(ctx, kns, "postgres-auth", nil)
	//s.Require().Nil(err)
	//
	//m := make(map[string]string)
	//
	//sec := v1.Secret{
	//	TypeMeta: metav1.TypeMeta{
	//		Kind:       "Secret",
	//		APIVersion: "v1",
	//	},
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:      "postgres-auth",
	//		Namespace: kns.Namespace,
	//	},
	//	StringData: m,
	//	Type:       "Opaque",
	//}
	//newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	//s.Require().Nil(err)
	//s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestGetSecretsUpdate() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "zeusfyi", Namespace: "5cf3a2c0-1d65-48cb-8b85-dc777ad956a0"}

	routes := []string{
		"http://anvil.191aada9-055d-4dba-a906-7dfbc4e632c6.svc.cluster.local:8545",
		"http://anvil.427c5536-4fc0-4257-90b5-1789d290058c.svc.cluster.local:8545",
		"http://anvil.5cf3a2c0-1d65-48cb-8b85-dc777ad956a0.svc.cluster.local:8545",
		"http://anvil.78ab2d4c-82eb-4bbc-b0fb-b702639e78c0.svc.cluster.local:8545",
		"http://anvil.a49ca82d-ff96-4c4f-8653-001d56cab5e5.svc.cluster.local:8545",
		"http://anvil.be58f278-1fbe-4bc8-8db5-03d8901cc060.svc.cluster.local:8545",
		"http://anvil.e56def19-190f-4b45-9fdb-8468ddbe0eb5.svc.cluster.local:8545",
		"http://anvil.eeb335ad-78da-458f-9cfb-9928514d65d0.svc.cluster.local:8545",
	}
	for i, r := range routes {
		r = strings.TrimPrefix(r, "http://anvil.")
		r = strings.TrimSuffix(r, ".svc.cluster.local:8545")
		kns.Namespace = r

		node := s.Tc.QuikNodeURLS.Routes[i]

		secret, err := s.K.GetSecretWithKns(ctx, kns, "hardhat", nil)
		s.Require().Nil(err)
		s.Require().NotEmpty(secret)
		sb := []byte(node)
		secret.Data["rpc"] = sb

		err = s.K.DeleteSecretWithKns(ctx, kns, "hardhat", nil)
		s.Require().Nil(err)

		sec := v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "hardhat",
				Namespace: kns.Namespace,
			},
			Data: secret.Data,
			Type: "Opaque",
		}
		newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
		s.Require().Nil(err)
		s.Require().NotEmpty(newSecret)

		secret, err = s.K.GetSecretWithKns(ctx, kns, "hardhat", nil)
		s.Require().Nil(err)
		s.Require().NotEmpty(secret)
		s.Require().Equal(secret.Data["rpc"], sb)
	}
}
func (s *SecretsTestSuite) TestCreateChoreographySecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "p2p-crawler"}
	m := make(map[string]string)

	m["bearer"] = "bearer"
	m["cloud-provider"] = "cloud"
	m["ctx"] = "ctx"
	m["ns"] = "ns"
	m["region"] = "region"

	sec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "choreography",
			Namespace: kns.Namespace,
		},
		StringData: m,
		Type:       "Opaque",
	}

	//_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	//s.Require().Nil(err)
	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateAwsDynamoDBSecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "cec2d631-0330-4792-8139-fa18752a93d8"}
	m := make(map[string]string)
	m["postgres-conn-str"] = s.Tc.ProdLocalDbPgconn
	m["dynamodb-access-key"] = s.Tc.AwsAccessKeyDynamoDB
	m["dynamodb-secret-key"] = s.Tc.AwsSecretKeyDynamoDB
	sec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dynamodb-auth",
			Namespace: kns.Namespace,
		},
		StringData: m,
		Type:       "Opaque",
	}

	err := s.K.DeleteSecretWithKns(ctx, kns, "dynamodb-auth", nil)
	s.Require().Nil(err)

	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateAwsSecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "zeus"}
	m := make(map[string]string)
	m["aws-access-key"] = s.Tc.AwsAccessKeyEks
	m["aws-secret-key"] = s.Tc.AwsSecretKeyEks
	m["aws-default-region"] = "us-west-1"
	sec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: kns.Namespace,
		},
		StringData: m,
		Type:       "Opaque",
	}

	_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateSecrets() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	kns.Namespace = "demo"
	secret.Namespace = kns.Namespace
	_, err = s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

	secret.ResourceVersion = ""
	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, secret, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCopySecrets() {
	//var knsFrom = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "zeus"}
	var knsTo = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "hestia"}
	fromKns := zeus_common_types.CloudCtxNs{
		CloudProvider: "do",
		Region:        "sfo3",
		Context:       "do-sfo3-dev-do-sfo3-zeus",
		Namespace:     "zeus",
		Alias:         "zeus",
		Env:           "",
	}

	secList := []string{"age-auth", "spaces-auth", "spaces-key", "zeus-fyi-ext"}
	for _, se := range secList {
		_, err := s.K.CopySecretToAnotherKns(ctx, fromKns, knsTo, se, nil)
		s.Require().Nil(err)
	}
	//_, err := s.K.CopySecretToAnotherKns(ctx, fromKns, knsTo, "aws-auth", nil)
	//s.Require().Nil(err)
}

func (s *SecretsTestSuite) TestCopySecretToAnotherNs() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "zeus"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	// change key & value here
	kns.Namespace = "zeus"
	secret.Namespace = kns.Namespace
	//_, err = s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	//s.Require().Nil(err)

	secret.ResourceVersion = ""
	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, secret, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateWeb3SignerSecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral-staking"}
	m := make(map[string]string)

	m["postgres-db"] = s.Tc.DevWeb3SignerPgconn
	m["postgres-username"] = "doadmin"
	m["postgres-auth"] = s.Tc.DevWeb3SignerPgconnAuth
	m["network"] = "minimal"

	sec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "web3signer",
			Namespace: kns.Namespace,
		},
		StringData: m,
		Type:       "Opaque",
	}

	_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateS3AwsSecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "nyc1", Context: "do-nyc1-do-nyc1-zeus-demo", Namespace: "sui-testnet-do-1daf4b8e"}
	m := make(map[string]string)
	m["AWS_ACCESS_KEY_ID"] = s.Tc.AwsS3AccessKey
	m["AWS_SECRET_ACCESS_KEY"] = s.Tc.AwsS3SecretKey
	m["AWS_REGION"] = "us-west-1"
	sec := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-credentials",
			Namespace: kns.Namespace,
		},
		StringData: m,
		Type:       "Opaque",
	}

	_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestDockerSecret() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "zeus"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "zeus-fyi-ext", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	auth := aegis_aws_auth.AuthAWS{
		Region:    "us-west-1",
		AccessKey: s.Tc.AwsAccessKeySecretManager,
		SecretKey: s.Tc.AwsSecretKeySecretManager,
	}
	artemis_hydra_orchestrations_aws_auth.InitHydraSecretManagerAuthAWS(ctx, auth)
	ps, err := aws_secrets.GetDockerSecret(ctx, s.Ou, "zeus-do-docker")
	s.Require().Nil(err)
	s.Require().NotEmpty(ps)

	s.Require().Equal(ps.DockerAuthJson, string(secret.Data[".dockerconfigjson"]))
}

func TestSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsTestSuite))
}
