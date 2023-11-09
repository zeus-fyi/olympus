package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func RunZeusDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Zeus: RunZeusDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunZeusDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.GcpAuthJsonBytes = sw.ReadSecretBytes(ctx, inMemSecrets, gcpAuthJson)
	sw.DoctlToken = sw.MustReadSecret(ctx, inMemSecrets, doctlSecret)
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.StripeSecretKey = sw.MustReadSecret(ctx, inMemSecrets, stripeSecretKey)
	sw.EksAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, eksAccessKey)
	sw.EksAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, eksSecretKey)
	// TODO allow for multiple regions
	sw.OvhAppKey = sw.MustReadSecret(ctx, inMemSecrets, ovhAppKey)
	sw.OvhSecretKey = sw.MustReadSecret(ctx, inMemSecrets, ovhSecretKey)
	sw.OvhConsumerKey = sw.MustReadSecret(ctx, inMemSecrets, ovhConsumerKey)
	sw.AwsS3AccessKey = sw.MustReadSecret(ctx, inMemSecrets, awsS3ReaderAccessKey)
	sw.AwsS3SecretKey = sw.MustReadSecret(ctx, inMemSecrets, awsS3ReaderSecretKey)
	sw.EksAuthAWS.Region = "us-west-1"
	hera_openai.InitHeraOpenAI(sw.OpenAIToken)

	InitAtlassianKeys(ctx, inMemSecrets, &sw)
	return inMemSecrets, sw
}
