package auth_startup

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	s3reader "github.com/zeus-fyi/olympus/datastores/s3/read"
	"github.com/zeus-fyi/olympus/pkg/aegis/s3secrets"
	temporal_auth "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/misc"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
	aegis_aws_auth "github.com/zeus-fyi/zeus/pkg/aegis/aws/auth"
)

const (
	temporalBearerSecret = "secrets/temporal.bearer.txt"
	PgSecret             = "secrets/postgres-auth.txt"
	doctlSecret          = "secrets/doctl.txt"
	rcloneSecret         = "secrets/rclone.conf"
	encryptedSecret      = "secrets.tar.gz.age"
	secretBucketName     = "zeus-fyi"
	pagerDutySecret      = "secrets/pagerduty.txt"
	pagerDutyRoutingKey  = "secrets/pagerduty.routing.key.txt"
	gcpAuthJson          = "secrets/zeusfyi-23264580e41d.json"

	gmailAuthJson = "secrets/zgmail.json"

	eksAccessKey = "secrets/aws.eks.access.key.txt"
	eksSecretKey = "secrets/aws.eks.secret.key.txt"

	ovhAppKey      = "secrets/ovh.app.key.txt"
	ovhSecretKey   = "secrets/ovh.secret.key.txt"
	ovhConsumerKey = "secrets/ovh.consumer.key.txt"

	zeroXApiKey = "secrets/zero.x.api.key.txt"

	googClientID     = "secrets/google.client.id.txt"
	googClientSecret = "secrets/google.client.secret.txt"
	googGtagSecret   = "secrets/google.gtag.secret.txt"

	awsS3ReaderAccessKey = "secrets/aws.s3.reader.access.key.txt"
	awsS3ReaderSecretKey = "secrets/aws.s3.reader.secret.key.txt"

	twitterConsumerPublicAPIKey = "secrets/twitter.consumer.public.api.key.txt"
	twitterConsumerSecretAPIKey = "secrets/twitter.consumer.secret.api.key.txt"
	twitterAccessToken          = "secrets/twitter.access.token.txt"
	twitterAccessTokenSecret    = "secrets/twitter.access.secret.token.txt"

	redditSecretsJson  = "secrets/reddit.api.keys.json"
	discordSecretsJson = "secrets/discord.auth.json"
)

type SecretsWrapper struct {
	TwitterConsumerPublicAPIKey string
	TwitterConsumerSecretAPIKey string
	TwitterAccessToken          string
	TwitterAccessTokenSecret    string

	OvhAppKey              string
	OvhSecretKey           string
	OvhConsumerKey         string
	PostgresAuth           string
	AegisPostgresAuth      string
	DoctlToken             string
	MainnetBeaconURL       string
	BearerToken            string
	OpenAIToken            string
	AccessKeyHydraDynamoDB string
	SecretKeyHydraDynamoDB string
	PagerDutyApiKey        string
	PagerDutyRoutingKey    string
	SendGridAPIKey         string
	GcpAuthJsonBytes       []byte
	GmailAuthJsonBytes     []byte

	GoogClientID     string
	GoogClientSecret string
	GoogGtagSecret   string

	QuickNodePassword string
	QuickNodeBearer   string
	QuickNodeJWT      string
	StripePubKey      string
	StripeSecretKey   string
	ZeroXApiKey       string

	SecretsManagerAuthAWS aegis_aws_auth.AuthAWS
	SESAuthAWS            aegis_aws_auth.AuthAWS
	EksAuthAWS            aegis_aws_auth.AuthAWS

	TemporalAuth temporal_auth.TemporalAuth

	AwsS3AccessKey  string
	AwsS3SecretKey  string
	AtlassianOrgId  string
	AtlassianApiKey string
	GmailApiKey     string

	RedditAuthConfig  RedditAuthConfig
	DiscordAuthConfig DiscordAuthConfig
}
type RedditAuthConfig struct {
	RedditUsername     string `json:"redditUsername"`
	RedditPassword     string `json:"redditPassword"`
	RedditSecretOAuth2 string `json:"redditSecretOAuth2"`
	RedditPublicOAuth2 string `json:"redditPublicOAuth2"`
}

type DiscordAuthConfig struct {
	DiscordClientID     string `json:"discordClientID"`
	DiscordClientSecret string `json:"discordClientSecret"`
}

var secretsBucket = &s3.GetObjectInput{
	Bucket: aws.String(secretBucketName),
	Key:    aws.String(encryptedSecret),
}

func (s *SecretsWrapper) MustReadSecret(ctx context.Context, inMemSecrets memfs.MemFS, fileName string) string {
	secret, err := inMemSecrets.ReadFile(fileName)
	if err != nil {
		log.Fatal().Msgf("SecretsWrapper: MustReadSecret failed, shutting down the server: %s", fileName)
		misc.DelayedPanic(err)
	}
	return string(secret)
}

func (s *SecretsWrapper) ReadSecret(ctx context.Context, inMemSecrets memfs.MemFS, fileName string) (string, error) {
	secret, err := inMemSecrets.ReadFile(fileName)
	if err != nil {
		log.Err(err).Msgf("SecretsWrapper: ReadSecret failed, shutting down the server: %s", fileName)
		return "", err
	}
	return string(secret), err
}

func (s *SecretsWrapper) ReadSecretBytes(ctx context.Context, inMemSecrets memfs.MemFS, fileName string) []byte {
	secret, err := inMemSecrets.ReadFile(fileName)
	if err != nil {
		log.Fatal().Msgf("SecretsWrapper: MustReadSecret failed, shutting down the server: %s", fileName)
		misc.DelayedPanic(err)
	}
	return secret
}

var (
	Sp = filepaths.Path{
		PackageName: "",
		DirIn:       "/secrets",
		DirOut:      "/secrets",
		FnIn:        "secrets.tar.gz.age",
		Env:         "",
		FilterFiles: string_utils.FilterOpts{},
	}
)

func ReadEncSecretsFromInMemDir() []byte {
	b, err := Sp.ReadFileInPath()
	if err != nil {
		log.Fatal().Err(err).Msg("ReadEncSecretsFromInMemDir: failed to read file")
		misc.DelayedPanic(err)
	}
	return b
}

func ReadEncryptedSecretsData(ctx context.Context, authCfg AuthConfig) memfs.MemFS {
	authCfg.S3KeyValue = secretsBucket
	s3Reader := s3reader.NewS3ClientReader(authCfg.s3BaseClient)
	s3SecretsReader := s3secrets.NewS3Secrets(authCfg.a, s3Reader)
	//buf := s3SecretsReader.ReadBytes(ctx, &authCfg.Path, authCfg.S3KeyValue)
	buf := ReadEncSecretsFromInMemDir()
	tmpPath := filepaths.Path{}
	tmpPath.DirOut = "./"
	tmpPath.FnOut = encryptedSecret
	err := s3SecretsReader.MemFS.MakeFileIn(&authCfg.Path, buf)
	if err != nil {
		log.Fatal().Msg("ReadEncryptedSecretsData: MakeFile failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	unzipDir := "./secrets"
	err = s3SecretsReader.DecryptAndUnGzipToInMemFs(&authCfg.Path, unzipDir)
	if err != nil {
		log.Fatal().Msg("ReadEncryptedSecretsData: DecryptAndUnGzipToInMemFs failed, shutting down the server")
		misc.DelayedPanic(err)
	}
	return s3SecretsReader.MemFS
}

func RunDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Zeus: RunDigitalOceanS3BucketObjSecretsProcedure starting")

	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.DoctlToken = sw.MustReadSecret(ctx, inMemSecrets, doctlSecret)
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.StripeSecretKey = sw.MustReadSecret(ctx, inMemSecrets, stripeSecretKey)
	sw.OpenAIToken = sw.MustReadSecret(ctx, inMemSecrets, heraOpenAIAuth)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)
	return inMemSecrets, sw
}

func RunArtemisDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Artemis: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Artemis: RunArtemisDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.BearerToken = sw.MustReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	sw.AccessKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, HydraAccessKeyDynamoDB)
	sw.SecretKeyHydraDynamoDB = sw.MustReadSecret(ctx, inMemSecrets, HydraSecretKeyDynamoDB)
	sw.ZeroXApiKey = sw.MustReadSecret(ctx, inMemSecrets, zeroXApiKey)
	sw.SendGridAPIKey = sw.MustReadSecret(ctx, inMemSecrets, sendGridAPIKey)
	log.Info().Msg("Artemis: RunArtemisDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func RunPoseidonDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Poseidon: RunDigitalOceanS3BucketObjSecretsProcedure starting")
	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("RunPoseidonDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)
	sw.BearerToken = sw.MustReadSecret(ctx, inMemSecrets, temporalBearerSecret)
	log.Info().Msg("RunPoseidonDigitalOceanS3BucketObjSecretsProcedure succeeded")
	return inMemSecrets, sw
}

func RunAthenaDigitalOceanS3BucketObjSecretsProcedure(ctx context.Context, authCfg AuthConfig) (memfs.MemFS, SecretsWrapper) {
	log.Info().Msg("Athena: RunDigitalOceanS3BucketObjSecretsProcedure starting")

	inMemSecrets := ReadEncryptedSecretsData(ctx, authCfg)
	log.Info().Msg("Athena: RunDigitalOceanS3BucketObjSecretsProcedure finished")
	sw := SecretsWrapper{}
	sw.PostgresAuth = sw.MustReadSecret(ctx, inMemSecrets, PgSecret)

	p := filepaths.Path{
		PackageName: "",
		DirIn:       "",
		DirOut:      "/root/.config/rclone",
		FnOut:       "rclone.conf",
		Env:         "",
		Metadata:    nil,
		FilterFiles: string_utils.FilterOpts{},
	}
	rcloneConf, err := inMemSecrets.ReadFile(rcloneSecret)
	if err != nil {
		log.Err(err).Msg("Athena:  RunAthenaDigitalOceanS3BucketObjSecretsProcedure failed to set rclone conf")
		misc.DelayedPanic(err)
	}
	err = p.WriteToFileOutPath(rcloneConf)
	if err != nil {
		log.Err(err).Msg("Athena:  RunAthenaDigitalOceanS3BucketObjSecretsProcedure failed to set rclone conf")
		misc.DelayedPanic(err)
	}
	return inMemSecrets, sw
}
