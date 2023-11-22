package auth_startup

import (
	"context"
	"encoding/json"

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

	// gmailAuthJson
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
	sw.OpenAIToken = sw.MustReadSecret(ctx, inMemSecrets, heraOpenAIAuth)

	hera_openai.InitHeraOpenAI(sw.OpenAIToken)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)
	sw.GmailApiKey = sw.MustReadSecret(ctx, inMemSecrets, gmailApiKey)
	sw.SecretsManagerAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)
	InitAtlassianKeys(ctx, inMemSecrets, &sw)

	sw.TwitterConsumerPublicAPIKey = sw.MustReadSecret(ctx, inMemSecrets, twitterConsumerPublicAPIKey)
	sw.TwitterConsumerSecretAPIKey = sw.MustReadSecret(ctx, inMemSecrets, twitterConsumerSecretAPIKey)
	sw.TwitterAccessToken = sw.MustReadSecret(ctx, inMemSecrets, twitterAccessToken)
	sw.TwitterAccessTokenSecret = sw.MustReadSecret(ctx, inMemSecrets, twitterAccessTokenSecret)

	ra := RedditAuthConfig{}
	sb := sw.ReadSecretBytes(ctx, inMemSecrets, redditSecretsJson)
	err := json.Unmarshal(sb, &ra)
	if err != nil {
		panic(err)
	}
	return inMemSecrets, sw
}
