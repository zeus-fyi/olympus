package configs

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_client "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var testCont TestContainer

type TestURLs struct {
	ProdZeusApiURL  string
	LocalZeusApiURL string
}

type QuickNodeURLS struct {
	TestRoute string
	Routes    []string
}
type QuickNodeMarketplace struct {
	Password  string
	JWTToken  string
	AuthToken string
}

type TestContainer struct {
	Env string

	AwsS3AccessKey string
	AwsS3SecretKey string

	QuickNodeMarketplace        QuickNodeMarketplace
	ZeroXApiKey                 string
	OvhAppKey                   string
	OvhSecretKey                string
	OvhConsumerKey              string
	QuikNodeStreamWsNode        string
	QuikNodeLiveNode            string
	HardhatNode                 string
	AwsAccessKeyEks             string
	AwsSecretKeyEks             string
	InfraCostAPIKey             string
	GcpAuthJson                 []byte
	TwitterAccessToken          string
	TwitterAccessTokenSecret    string
	TwitterBearerToken          string
	TwitterConsumerPublicAPIKey string
	TwitterConsumerSecretAPIKey string
	RedditUsername              string
	RedditPassword              string
	RedditSecretOAuth2          string
	RedditPublicOAuth2          string
	StripeTestPublicAPIKey      string
	StripeTestSecretAPIKey      string
	StripeProdPublicAPIKey      string
	StripeProdSecretAPIKey      string

	DigitalOceanAPIKey      string
	LocalDbPgconn           string
	StagingDbPgconn         string
	ProdDbPgconn            string
	ProdLocalDbPgconn       string
	ProdLocalApolloDbPgconn string

	LocalBeaconConn  string
	LocalRedisConn   string
	StagingRedisConn string
	ProdRedisConn    string

	LocalAgePubkey      string
	LocalAgePkey        string
	LocalS3SpacesKey    string
	LocalS3SpacesSecret string

	ProductionLocalTemporalOrgID  int
	ProductionLocalTemporalUserID int

	LocalBearerToken                   string
	ProductionLocalBearerToken         string
	ProductionLocalTemporalBearerToken string
	DemoUserBearerToken                string

	DevTemporalHostPort string
	DevTemporalNs       string

	DevTemporalAuth temporal_client.TemporalAuth

	DevAuthKeysCfg auth_keys_config.AuthKeysCfg

	ProdLocalTemporalAuth temporal_client.TemporalAuth
	ProdLocalAuthKeysCfg  auth_keys_config.AuthKeysCfg

	ProdLocalTemporalAuthArtemis  temporal_client.TemporalAuth
	ProdLocalTemporalAuthPoseidon temporal_client.TemporalAuth
	ProdLocalTemporalAuthHestia   temporal_client.TemporalAuth

	TestURLs

	ArtemisHexKeys

	LocalEcsdaTestPkey  string
	LocalEcsdaTestPkey2 string

	EphemeralNodeUrl string
	GoerliNodeUrl    string
	MainnetNodeUrl   string

	OpenAIAuth string

	DevWeb3SignerPgconn     string
	DevWeb3SignerPgconnAuth string

	AwsAccessKey string
	AwsSecretKey string

	AwsAccessKeySecretManager string
	AwsSecretKeySecretManager string

	AwsAccessKeyDynamoDB string
	AwsSecretKeyDynamoDB string

	AwsAccessKeySES string
	AwsSecretKeySES string

	SendGridAPIKey        string
	AwsAccessKeyLambdaExt string
	AwsSecretKeyLambdaExt string
	AwsLamdbaTestURL      string

	TestEthKeyOneBLS string
	TestEthKeyTwoBLS string

	PagerDutyApiKey     string
	PagerDutyRoutingKey string
	AdminLoginPassword  string

	EtherScanAPIKey string

	QuikNodeURLS QuickNodeURLS

	GoogClientID     string
	GoogClientSecret string
	GoogTagSecret    string

	AtlassianKeys
}

type AtlassianKeys struct {
	OrgId  string
	ApiKey string
}

type ArtemisHexKeys struct {
	ArtemisEphemeralEcdsaKey string
	ArtemisGoerliEcdsaKey    string
	ArtemisMainnetEcdsaKey   string
}

func SetBaseURLs() TestURLs {
	tu := TestURLs{}
	tu.ProdZeusApiURL = viper.GetString("PROD_ZEUS_URL")
	tu.LocalZeusApiURL = viper.GetString("LOCAL_ZEUS_URL")
	return tu
}

func forceDirToCallerLocation() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func InitEnvFromConfig(dir string) {
	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml") // for a YAML file
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func InitArtemisLocalAccounts() {
	testCont.ArtemisMainnetEcdsaKey = viper.GetString("PROD_MAINNET_ARTEMIS_ECDSA_PKEY")
	testCont.ArtemisGoerliEcdsaKey = viper.GetString("PROD_GOERLI_ARTEMIS_ECDSA_PKEY")
}

func InitLocalTestConfigs() TestContainer {
	InitEnvFromConfig(forceDirToCallerLocation())
	p := filepaths.Path{
		PackageName: "",
		DirIn:       "secrets",
		FnIn:        "zeusfyi-23264580e41d.json",
	}
	b, err := p.ReadFileInPath()
	if err != nil {
		log.Info().Err(err).Msg("error reading gcp auth json file")
	}
	qn := QuickNodeURLS{
		TestRoute: viper.GetString("QUIKNODE_TEST"),
		Routes:    []string{},
	}
	for i := 1; i < 9; i++ {
		qn.Routes = append(qn.Routes, viper.GetString(fmt.Sprintf("QUIKNODE_%d", i)))
	}

	testCont.AtlassianKeys.ApiKey = viper.GetString("ATLASSIAN_API_KEY")
	testCont.AtlassianKeys.OrgId = viper.GetString("ATLASSIAN_ORG_ID")

	testCont.AwsS3AccessKey = viper.GetString("AWS_S3_ACCESS_KEY")
	testCont.AwsS3SecretKey = viper.GetString("AWS_S3_SECRET_KEY")

	testCont.GoogTagSecret = viper.GetString("GOOGLE_GTAG_SECRET")

	testCont.GoogClientID = viper.GetString("GOOGLE_CLIENT_ID")
	testCont.GoogClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	testCont.QuickNodeMarketplace.Password = viper.GetString("QUICKNODE_PASSWORD")
	testCont.QuickNodeMarketplace.JWTToken = viper.GetString("QUICKNODE_JWT")
	testCont.QuickNodeMarketplace.AuthToken = viper.GetString("QUICKNODE_AUTH_TOKEN")

	testCont.QuikNodeURLS = qn
	testCont.GcpAuthJson = b
	testCont.EtherScanAPIKey = viper.GetString("ETHERSCAN_API_KEY")
	testCont.HardhatNode = viper.GetString("HARDHAT_NODE_URL")
	testCont.TwitterAccessToken = viper.GetString("TWITTER_ACCESS_TOKEN_KEY")
	testCont.TwitterAccessTokenSecret = viper.GetString("TWITTER_ACCESS_TOKEN_SECRET_KEY")

	testCont.ZeroXApiKey = viper.GetString("ZERO_X_API_KEY")
	testCont.QuikNodeLiveNode = viper.GetString("QUIKNODE_LIVE_NODE_URL")
	testCont.OvhAppKey = viper.GetString("OVH_APP_KEY")
	testCont.OvhSecretKey = viper.GetString("OVH_SECRET_KEY")
	testCont.OvhConsumerKey = viper.GetString("OVH_CONSUMER_KEY")
	testCont.QuikNodeStreamWsNode = viper.GetString("QUIKNODE_STREAM_WS_URL")
	testCont.InfraCostAPIKey = viper.GetString("INFRA_COST_API_KEY")
	testCont.TwitterBearerToken = viper.GetString("TWITTER_BEARER_TOKEN")
	testCont.TwitterConsumerPublicAPIKey = viper.GetString("TWITTER_PUBLIC_API_KEY")
	testCont.TwitterConsumerSecretAPIKey = viper.GetString("TWITTER_SECRET_API_KEY")
	testCont.RedditUsername = viper.GetString("REDDIT_USERNAME")
	testCont.RedditPassword = viper.GetString("REDDIT_PASSWORD")
	testCont.RedditPublicOAuth2 = viper.GetString("REDDIT_PUBLIC_OAUTH2")
	testCont.RedditSecretOAuth2 = viper.GetString("REDDIT_SECRET_OAUTH2")
	testCont.StripeTestPublicAPIKey = viper.GetString("STRIPE_TEST_API_PUBLIC_KEY")
	testCont.StripeTestSecretAPIKey = viper.GetString("STRIPE_TEST_API_SECRET_KEY")

	testCont.StripeProdPublicAPIKey = viper.GetString("STRIPE_PROD_API_PUBLIC_KEY")
	testCont.StripeProdSecretAPIKey = viper.GetString("STRIPE_PROD_API_SECRET_KEY")

	testCont.DigitalOceanAPIKey = viper.GetString("DO_API_KEY")
	testCont.SendGridAPIKey = viper.GetString("SENDGRID_API_KEY")
	testCont.AwsAccessKeySES = viper.GetString("AWS_ACCESS_KEY_SES")
	testCont.AwsSecretKeySES = viper.GetString("AWS_SECRET_KEY_SES")

	testCont.AdminLoginPassword = viper.GetString("ADMIN_LOGIN_PW")
	testCont.AwsAccessKeyDynamoDB = viper.GetString("AWS_ACCESS_KEY_DYNAMODB")
	testCont.AwsSecretKeyDynamoDB = viper.GetString("AWS_SECRET_KEY_DYNAMODB")

	testCont.TestEthKeyOneBLS = viper.GetString("BLS_ETH_TEST_SK_ONE")
	testCont.TestEthKeyTwoBLS = viper.GetString("BLS_ETH_TEST_SK_TWO")

	testCont.AwsAccessKey = viper.GetString("AWS_ACCESS_KEY")
	testCont.AwsSecretKey = viper.GetString("AWS_SECRET_KEY")

	testCont.AwsAccessKeyEks = viper.GetString("AWS_ACCESS_KEY_EKS")
	testCont.AwsSecretKeyEks = viper.GetString("AWS_SECRET_KEY_EKS")
	testCont.AwsAccessKeySecretManager = viper.GetString("AWS_ACCESS_KEY_SECRET_MANAGER")
	testCont.AwsSecretKeySecretManager = viper.GetString("AWS_SECRET_KEY_SECRET_MANAGER")

	testCont.AwsLamdbaTestURL = viper.GetString("BLS_SERVERLESS_LAMBA_FUNC_ADDR")

	testCont.AwsAccessKeyLambdaExt = viper.GetString("AWS_LAMBDA_INVOKE_ACCESS_KEY")
	testCont.AwsSecretKeyLambdaExt = viper.GetString("AWS_LAMBDA_INVOKE_SECRET_KEY")

	testCont.DevWeb3SignerPgconn = viper.GetString("WEB3SIGNER_PG_DB")
	testCont.DevWeb3SignerPgconnAuth = viper.GetString("WEB3SIGNER_PG_AUTH_DEV")

	testCont.OpenAIAuth = viper.GetString("OPEN_AI_AUTH")

	testCont.MainnetNodeUrl = viper.GetString("MAINNET_NODE_URL")
	testCont.GoerliNodeUrl = viper.GetString("GOERLI_NODE_URL")
	testCont.EphemeralNodeUrl = viper.GetString("EPHEMERAL_NODE_URL")

	testCont.ProductionLocalTemporalOrgID = viper.GetInt("PROD_LOCAL_TEMPORAL_ORG_ID")
	testCont.ProductionLocalTemporalUserID = viper.GetInt("PROD_LOCAL_TEMPORAL_USER_ID")

	InitArtemisLocalAccounts()
	// local test keys
	testCont.LocalEcsdaTestPkey = viper.GetString("LOCAL_TESTING_ECDSA_PKEY")
	testCont.LocalEcsdaTestPkey2 = viper.GetString("LOCAL_TESTING_ECDSA_PKEY_2")

	// urls & env
	testCont.TestURLs = SetBaseURLs()
	testCont.Env = viper.GetString("ENV")

	// demo user for testing
	testCont.DemoUserBearerToken = viper.GetString("DEMO_USER_BEARER_TOKEN")

	// temporal auth
	testCont.ProductionLocalTemporalBearerToken = viper.GetString("PROD_LOCAL_TEMPORAL_BEARER_TOKEN")

	// temporal zeus
	testCont.DevTemporalNs = viper.GetString("DEV_TEMPORAL_NS")
	testCont.DevTemporalHostPort = viper.GetString("DEV_TEMPORAL_HOST_PORT")
	certPath := "./zeus.fyi/ca.pem"
	pemPath := "./zeus.fyi/ca.key"
	namespace := testCont.DevTemporalNs
	hostPort := testCont.DevTemporalHostPort
	testCont.DevTemporalAuth = temporal_client.TemporalAuth{
		ClientCertPath:   certPath,
		ClientPEMKeyPath: pemPath,
		Namespace:        namespace,
		HostPort:         hostPort,
	}

	testCont.ProdLocalTemporalAuth = testCont.DevTemporalAuth
	testCont.ProdLocalTemporalAuth.Namespace = viper.GetString("PROD_LOCAL_TEMPORAL_NS")
	testCont.ProdLocalTemporalAuth.HostPort = viper.GetString("PROD_LOCAL_TEMPORAL_HOST_PORT")

	// temporal artemis
	testCont.ProdLocalTemporalAuthArtemis.Namespace = viper.GetString("PROD_LOCAL_ARTEMIS_TEMPORAL_NS")
	testCont.ProdLocalTemporalAuthArtemis.HostPort = viper.GetString("PROD_LOCAL_ARTEMIS_TEMPORAL_HOST_PORT")
	testCont.ProdLocalTemporalAuthArtemis.ClientPEMKeyPath = "./zeus.fyi/ca.key"
	testCont.ProdLocalTemporalAuthArtemis.ClientCertPath = "./zeus.fyi/ca.pem"
	// temporal hestia
	testCont.ProdLocalTemporalAuthHestia.Namespace = viper.GetString("PROD_LOCAL_HESTIA_TEMPORAL_NS")
	testCont.ProdLocalTemporalAuthHestia.HostPort = viper.GetString("PROD_LOCAL_HESTIA_TEMPORAL_HOST_PORT")
	testCont.ProdLocalTemporalAuthHestia.ClientPEMKeyPath = "./zeus.fyi/ca.key"
	testCont.ProdLocalTemporalAuthHestia.ClientCertPath = "./zeus.fyi/ca.pem"

	// temporal poseidon
	testCont.ProdLocalTemporalAuthPoseidon = testCont.DevTemporalAuth
	testCont.ProdLocalTemporalAuthPoseidon.Namespace = viper.GetString("PROD_LOCAL_POSEIDON_TEMPORAL_NS")
	testCont.ProdLocalTemporalAuthPoseidon.HostPort = viper.GetString("PROD_LOCAL_POSEIDON_TEMPORAL_HOST_PORT")

	// age keys
	testCont.LocalAgePubkey = viper.GetString("LOCAL_AGE_PUBKEY")
	testCont.LocalAgePkey = viper.GetString("LOCAL_AGE_PKEY")

	testCont.LocalS3SpacesKey = viper.GetString("LOCAL_S3_SPACES_KEY")
	testCont.LocalS3SpacesSecret = viper.GetString("LOCAL_S3_SPACES_SECRET")

	testCont.LocalRedisConn = viper.GetString("LOCAL_REDIS_CONN")
	testCont.StagingRedisConn = viper.GetString("STAGING_REDIS_CONN")
	testCont.ProdRedisConn = viper.GetString("PROD_REDIS_CONN")

	testCont.LocalDbPgconn = viper.GetString("LOCAL_DB_PGCONN")
	testCont.StagingDbPgconn = viper.GetString("STAGING_DB_PGCONN")
	testCont.ProdDbPgconn = viper.GetString("PROD_DB_PGCONN")
	testCont.ProdLocalDbPgconn = viper.GetString("PROD_LOCAL_DB_PGCONN")
	testCont.ProdLocalApolloDbPgconn = viper.GetString("PROD_LOCAL_APOLLO_PGCONN")
	testCont.LocalBeaconConn = viper.GetString("LOCAL_BEACON_CONN_STR")

	testCont.LocalBearerToken = viper.GetString("LOCAL_BEARER_TOKEN")
	testCont.ProductionLocalBearerToken = viper.GetString("PROD_LOCAL_BEARER_TOKEN")

	testCont.PagerDutyApiKey = viper.GetString("PAGERDUTY_API_KEY")
	testCont.PagerDutyRoutingKey = viper.GetString("PAGERDUTY_ROUTING_KEY")
	testCont.DevAuthKeysCfg = getDevAuthKeysCfg()
	testCont.ProdLocalAuthKeysCfg = testCont.DevAuthKeysCfg
	return testCont
}

func getDevAuthKeysCfg() auth_keys_config.AuthKeysCfg {
	var DevAuthKeysCfg auth_keys_config.AuthKeysCfg
	DevAuthKeysCfg.AgePubKey = testCont.LocalAgePubkey
	DevAuthKeysCfg.AgePrivKey = testCont.LocalAgePkey
	DevAuthKeysCfg.SpacesKey = testCont.LocalS3SpacesKey
	DevAuthKeysCfg.SpacesPrivKey = testCont.LocalS3SpacesSecret
	return DevAuthKeysCfg
}
func InitProductionConfigs() TestContainer {
	InitEnvFromConfig(forceDirToCallerLocation())
	testCont.Env = "production"
	testCont.LocalRedisConn = viper.GetString("LOCAL_REDIS_CONN")
	testCont.StagingRedisConn = viper.GetString("STAGING_REDIS_CONN")
	testCont.LocalDbPgconn = viper.GetString("LOCAL_DB_PGCONN")
	testCont.StagingDbPgconn = viper.GetString("STAGING_DB_PGCONN")
	testCont.ProdDbPgconn = viper.GetString("PROD_DB_PGCONN")
	testCont.LocalBeaconConn = viper.GetString("LOCAL_BEACON_CONN_STR")
	testCont.ProdLocalApolloDbPgconn = viper.GetString("PROD_LOCAL_APOLLO_PGCONN")
	return testCont
}

func InitStagingConfigs() TestContainer {
	InitEnvFromConfig(forceDirToCallerLocation())
	testCont.Env = "staging"
	testCont.LocalRedisConn = viper.GetString("LOCAL_REDIS_CONN")
	testCont.StagingRedisConn = viper.GetString("STAGING_REDIS_CONN")
	testCont.LocalDbPgconn = viper.GetString("LOCAL_DB_PGCONN")
	testCont.StagingDbPgconn = viper.GetString("STAGING_DB_PGCONN")
	testCont.ProdDbPgconn = viper.GetString("PROD_DB_PGCONN")
	testCont.LocalBeaconConn = viper.GetString("LOCAL_BEACON_CONN_STR")
	return testCont
}
