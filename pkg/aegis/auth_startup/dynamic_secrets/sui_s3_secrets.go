package dynamic_secrets

import (
	"context"

	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup"
	"github.com/zeus-fyi/zeus/zeus/z_client/zeus_common_types"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	awsS3ReaderAccessKey = "secrets/aws.s3.reader.access.key.txt"
	awsS3ReaderSecretKey = "secrets/aws.s3.reader.secret.key.txt"
)

func GetS3SecretSui(ctx context.Context, kns zeus_common_types.CloudCtxNs) v1.Secret {
	sw := auth_startup.SecretsWrapper{}
	m := make(map[string]string)
	m["AWS_ACCESS_KEY_ID"] = sw.MustReadSecret(ctx, AegisInMemSecrets, awsS3ReaderAccessKey)
	m["AWS_SECRET_ACCESS_KEY"] = sw.MustReadSecret(ctx, AegisInMemSecrets, awsS3ReaderSecretKey)
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
	return sec
}

func GetS3FakeSecretSui(ctx context.Context, kns zeus_common_types.CloudCtxNs) v1.Secret {
	//sw := auth_startup.SecretsWrapper{}
	m := make(map[string]string)
	m["AWS_ACCESS_KEY_ID"] = "fake"
	m["AWS_SECRET_ACCESS_KEY"] = "fake"
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
	return sec
}
