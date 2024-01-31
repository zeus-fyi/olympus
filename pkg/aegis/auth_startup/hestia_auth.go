package auth_startup

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	hera_openai "github.com/zeus-fyi/olympus/pkg/hera/openai"
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

	quicknodePassword = "secrets/quicknode.http.basic.password.txt"
	quicknodeBearer   = "secrets/quicknode.http.bearer.txt"
	quicknodeJWT      = "secrets/quicknode.jwt.txt"

	gmailApiKey = "secrets/gmail.api.key.txt"
)

func RunHestiaDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sk, err := sw.ReadSecret(ctx, inMemSecrets, hestiaSessionKey)
	if err != nil {
		log.Err(err).Msg("error reading session key")
		err = nil
	} else {
		sw.HestiaSessionKey = sk
	}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.BearerToken = sw.MustReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	sw.SecretsManagerAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerAccessKey)
	sw.SecretsManagerAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, secretsManagerSecretKey)

	sw.PagerDutyApiKey = sw.MustReadSecret(ctx, inMemSecrets, pagerDutySecret)
	sw.PagerDutyRoutingKey = sw.MustReadSecret(ctx, inMemSecrets, pagerDutyRoutingKey)

	sw.SESAuthAWS.AccessKey = sw.MustReadSecret(ctx, inMemSecrets, sesAccessKey)
	sw.SESAuthAWS.SecretKey = sw.MustReadSecret(ctx, inMemSecrets, sesSecretKey)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)
	sw.StripePubKey = sw.MustReadSecret(ctx, inMemSecrets, stripePublishableKey)
	sw.StripeSecretKey = sw.MustReadSecret(ctx, inMemSecrets, stripeSecretKey)
	sw.QuickNodeBearer = sw.MustReadSecret(ctx, inMemSecrets, quicknodeBearer)
	sw.QuickNodeJWT = sw.MustReadSecret(ctx, inMemSecrets, quicknodeJWT)
	sw.QuickNodePassword = sw.MustReadSecret(ctx, inMemSecrets, quicknodePassword)

	sw.GoogClientID = sw.MustReadSecret(ctx, inMemSecrets, googClientID)
	sw.GoogClientSecret = sw.MustReadSecret(ctx, inMemSecrets, googClientSecret)
	sw.GoogGtagSecret = sw.MustReadSecret(ctx, inMemSecrets, googGtagSecret)

	InitAtlassianKeys(ctx, inMemSecrets, &sw)
	sw.GmailApiKey = sw.MustReadSecret(ctx, inMemSecrets, gmailApiKey)
	sw.OpenAIToken = sw.MustReadSecret(ctx, inMemSecrets, heraOpenAIAuth)
	sw.GmailAuthJsonBytes = sw.ReadSecretBytes(ctx, inMemSecrets, gmailAuthJson)

	cid, err := sw.ReadSecret(ctx, inMemSecrets, twitterClientID)
	if err != nil {
		log.Err(err).Msg("error reading twitter client id")
		err = nil
	} else {
		sw.TwitterMbClientID = cid
	}
	cs, err := sw.ReadSecret(ctx, inMemSecrets, twitterClientSecret)
	if err != nil {
		log.Err(err).Msg("error reading twitter client id")
		err = nil
	} else {
		sw.TwitterMbClientSecret = cs
	}
	hera_openai.InitHeraOpenAI(sw.OpenAIToken)
	log.Info().Msg("Hestia: RunDigitalOceanS3BucketObjSecretsProcedure succeeded")

	da := DiscordAuthConfig{}
	sb := sw.ReadSecretBytes(ctx, inMemSecrets, discordSecretsJson)
	err = json.Unmarshal(sb, &da)
	if err != nil {
		panic(err)
	}
	sw.DiscordAuthConfig = da
	return inMemSecrets, sw
}
