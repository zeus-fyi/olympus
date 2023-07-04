package zeus_core

import (
	"testing"

	"github.com/stretchr/testify/suite"
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

	secret, err := s.K.GetSecretWithKns(ctx, kns, "hardhat", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	sb := []byte("")
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
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "ovh", Region: "us-west-or-1", Context: "kubernetes-admin@zeusfyi", Namespace: "3cfa3022-5f71-478c-bebd-3b2e20f4caea"}
	m := make(map[string]string)
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

	_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
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
	var knsFrom = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "zeus"}
	var knsTo = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "beacon"}

	_, err := s.K.CopySecretToAnotherKns(ctx, knsFrom, knsTo, "spaces-auth", nil)
	s.Require().Nil(err)
}

func (s *SecretsTestSuite) TestCopySecretToAnotherNs() {
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)

	// change key & value here
	kns.Namespace = "zeus"
	secret.Namespace = kns.Namespace
	_, err = s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

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

func TestSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsTestSuite))
}
