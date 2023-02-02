package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/datastores/postgres/apps/hestia/models/read/auth"
	s3base "github.com/zeus-fyi/olympus/datastores/s3"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/encryption"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

type AuthConfig struct {
	Path         filepaths.Path
	a            encryption.Age
	s3BaseClient s3base.S3Client
	S3KeyValue   *s3.GetObjectInput
}

func FetchTemporalAuthBearer(ctx context.Context) string {
	key, err := auth.FetchTemporalAuthToken(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("SetTemporalAuthBearer: failed to find auth token, shutting down the server")
		misc.DelayedPanic(err)
	}
	return key.PublicKey
}

func NewDefaultAuthClient(ctx context.Context, keysCfg auth_keys_config.AuthKeysCfg) AuthConfig {
	if len(keysCfg.AgePrivKey) <= 0 {
		log.Warn().Msg("no age priv key provided, auth will fail")
	}
	a := encryption.NewAge(keysCfg.AgePrivKey, keysCfg.AgePubKey)
	s3BaseClient := NewDigitalOceanS3AuthClient(ctx, keysCfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String("zeus-fyi"),
		Key:    aws.String("kube.tar.gz.age"),
	}
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "./",
		DirOut:      "./",
		FnIn:        "kube.tar.gz.age",
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

	err := s3SecretsReader.MemFS.MakeFileIn(&authCfg.Path, buf.Bytes())
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
