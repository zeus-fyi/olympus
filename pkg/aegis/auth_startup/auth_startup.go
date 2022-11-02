package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AuthConfig struct {
	Path         structs.Path
	a            encryption.Age
	s3BaseClient s3base.S3Client
	S3KeyValue   *s3.GetObjectInput
}

type AuthKeysCfg struct {
	AgePrivKey    string
	AgePubKey     string
	SpacesKey     string
	SpacesPrivKey string
}

func NewDefaultAuthClient(ctx context.Context, keysCfg AuthKeysCfg) AuthConfig {
	a := encryption.NewAge(keysCfg.AgePrivKey, keysCfg.AgePubKey)
	s3BaseClient, err := s3base.NewConnS3ClientWithStaticCreds(ctx, keysCfg.SpacesKey, keysCfg.SpacesPrivKey)
	if err != nil {
		panic(err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("kube.tar.gz.age"),
	}
	p := structs.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		Fn:          "kube.tar.gz.age",
		FnOut:       "kube.tar.gz",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
	authCfg := AuthConfig{
		Path:         p,
		a:            a,
		s3BaseClient: s3BaseClient,
		S3KeyValue:   input,
	}
	return authCfg
}

func RunDigitalOceanS3BucketObjAuthProcedure(ctx context.Context, authCfg AuthConfig) memfs.MemFS {
	s3Reader := s3reader.NewS3ClientReader(authCfg.s3BaseClient)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	err := s3SecretsReader.Read(ctx, &authCfg.Path, authCfg.S3KeyValue)
	if err != nil {
		panic(err)
	}
	return s3SecretsReader.MemFS
}

//func RunDigitalOceanS3BucketObjAuthProcedure(ctx context.Context, authCfg AuthConfig) (s3secrets.S3Secrets, error) {
//	s3ClientBase := s3base.NewS3ClientBase()
//	err := s3ClientBase.ConnectS3SpacesDO(ctx)
//	if err != nil {
//		return s3secrets.S3Secrets{}, err
//	}
//	s3Reader := s3reader.NewS3ClientReader(s3ClientBase)
//	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
//	err = s3SecretsReader.Read(ctx, &authCfg.Path, authCfg.S3KeyValue)
//	if err != nil {
//		return s3SecretsReader, err
//	}
//	return s3SecretsReader, err
//}
