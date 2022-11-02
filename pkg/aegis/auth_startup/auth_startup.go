package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

type AuthConfig struct {
	Path       structs.Path
	a          encryption.Age
	S3KeyValue *s3.GetObjectInput
}

func RunDigitalOceanS3BucketObjAuthProcedure(ctx context.Context, authCfg AuthConfig) (s3secrets.S3Secrets, error) {
	s3ClientBase := s3base.NewS3ClientBase()
	err := s3ClientBase.ConnectS3SpacesDO(ctx)
	if err != nil {
		return s3secrets.S3Secrets{}, err
	}
	s3Reader := s3reader.NewS3ClientReader(s3ClientBase)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	err = s3SecretsReader.Read(ctx, &authCfg.Path, authCfg.S3KeyValue)
	if err != nil {
		return s3SecretsReader, err
	}
	return s3SecretsReader, err
}
