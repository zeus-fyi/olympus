package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeus-fyi/olympus/pkg/zeus/core/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretsTestSuite struct {
	K8TestSuite
}

func (s *SecretsTestSuite) TestGetSecrets() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)
}

func (s *SecretsTestSuite) TestCreateChoreographySecret() {
	ctx := context.Background()
	var kns = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "do-sfo3-dev-do-sfo3-zeus", Namespace: "ephemeral-staking"}
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

	_, err := s.K.CreateNamespaceIfDoesNotExist(ctx, kns)
	s.Require().Nil(err)

	newSecret, err := s.K.CreateSecretWithKns(ctx, kns, &sec, nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(newSecret)
}

func (s *SecretsTestSuite) TestCreateSecrets() {
	ctx := context.Background()
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
	ctx := context.Background()
	var knsFrom = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "zeus"}
	var knsTo = zeus_common_types.CloudCtxNs{CloudProvider: "do", Region: "sfo3", Context: "dev-sfo3-zeus", Namespace: "beacon"}

	_, err := s.K.CopySecretToAnotherKns(ctx, knsFrom, knsTo, "spaces-auth", nil)
	s.Require().Nil(err)
}

func (s *SecretsTestSuite) TestCopySecretToAnotherNs() {
	ctx := context.Background()
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

func TestSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsTestSuite))
}
