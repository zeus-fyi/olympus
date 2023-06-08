package auth_startup

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

const (
	secretsManagerAccessKey = "secrets/aws.secrets.manager.access.key.txt"
	secretsManagerSecretKey = "secrets/aws.secrets.manager.secret.key.txt"

	sesAccessKey = "secrets/aws.ses.access.key.txt"
	sesSecretKey = "secrets/aws.ses.secret.key.txt"

	sendGridAPIKey = "secrets/sendgrid.api.key.txt"

	stripePublishableKey = "secrets/stripe.api.access.key.txt"
	stripeSecretKey      = "secrets/stripe.api.secret.key.txt"
)

func RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, pgSecret)
	sw.BearerToken = sw.MustReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	sw.SecretsManagerAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)

	sw.SESAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, sesAccessKey)
	sw.SESAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, sesSecretKey)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)
	sw.StripePubKey = sw.MustReadSecret(ctx, inMemSecrets, stripePublishableKey)
	sw.StripeSecretKey = sw.MustReadSecret(ctx, inMemSecrets, stripeSecretKey)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}
