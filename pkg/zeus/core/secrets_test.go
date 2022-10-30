package zeus_core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SecretsTestSuite struct {
	K8TestSuite
}

func (s *SecretsTestSuite) TestGetSecrets() {
	ctx := context.Background()
	var kns = KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "dev-sfo3-zeus", Namespace: "eth-indexer"}

	secret, err := s.K.GetSecretWithKns(ctx, kns, "postgres-auth", nil)
	s.Require().Nil(err)
	s.Require().NotEmpty(secret)
}

func (s *SecretsTestSuite) TestCreateSecrets() {
	ctx := context.Background()
	var kns = KubeCtxNs{CloudProvider: "do", Region: "sfo3", CtxType: "dev-sfo3-zeus", Namespace: "eth-indexer"}

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
func TestSecretsTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsTestSuite))
}
