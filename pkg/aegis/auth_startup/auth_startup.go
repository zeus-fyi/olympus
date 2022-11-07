package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AuthConfig struct {
	Path         structs.Path
	a            encryption.Age
	s3BaseClient s3base.S3Client
	S3KeyValue   *s3.GetObjectInput
}

func NewDefaultAuthClient(ctx context.Context, keysCfg auth_keys_config.AuthKeysCfg) AuthConfig {
	a := encryption.NewAge(keysCfg.AgePrivKey, keysCfg.AgePubKey)
	s3BaseClient, err := s3base.NewConnS3ClientWithStaticCreds(ctx, keysCfg.SpacesKey, keysCfg.SpacesPrivKey)
	if err != nil {
		log.Fatal().Msg("NewDefaultAuthClient: NewConnS3ClientWithStaticCreds failed, shutting down the server")
		misc.DelayedPanic(err)
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
	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjAuthProcedure starting")

	s3Reader := s3reader.NewS3ClientReader(authCfg.s3BaseClient)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	buf := s3SecretsReader.ReadBytes(ctx, &authCfg.Path, authCfg.S3KeyValue)

	tmpPath := structs.Path{}
	tmpPath.DirOut = "./"
	tmpPath.FnOut = "kube.tar.gz.age"
	err := s3SecretsReader.MemFS.MakeFile(&authCfg.Path, buf.Bytes())
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjAuthProcedure: MakeFile failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	unzipDir := "./.kube"
	err = s3SecretsReader.DecryptAndUnGzipToInMemFs(&authCfg.Path, unzipDir)
	if err != nil {
		log.Fatal().Msg("RunDigitalOceanS3BucketObjAuthProcedure: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}

	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjAuthProcedure finished")
	return s3SecretsReader.MemFS
}
