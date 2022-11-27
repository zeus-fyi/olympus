package configs

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
	"github.com/zeus-fyi/olympus/pkg/aegis/auth_startup/auth_keys_config"
	temporal_client "github.com/zeus-fyi/olympus/pkg/iris/temporal/auth"
)

var testCont TestContainer

type TestURLs struct {
	ProdZeusApiURL  string
	LocalZeusApiURL string
}

type TestContainer struct {
	Env string

	LocalDbPgconn     string
	StagingDbPgconn   string
	ProdDbPgconn      string
	ProdLocalDbPgconn string

	LocalBeaconConn  string
	LocalRedisConn   string
	StagingRedisConn string
	ProdRedisConn    string

	LocalAgePubkey      string
	LocalAgePkey        string
	LocalS3SpacesKey    string
	LocalS3SpacesSecret string

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

	ProdLocalTemporalAuthArtemis temporal_client.TemporalAuth

	TestURLs

	ArtemisHexKeys

	LocalEcsdaTestPkey  string
	LocalEcsdaTestPkey2 string

	GoerliNodeUrl  string
	MainnetNodeUrl string
}

type ArtemisHexKeys struct {
	ArtemisGoerliEcdsaKey  string
	ArtemisMainnetEcdsaKey string
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
	testCont.MainnetNodeUrl = viper.GetString("MAINNET_NODE_URL")
	testCont.GoerliNodeUrl = viper.GetString("GOERLI_NODE_URL")

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
	testCont.LocalBeaconConn = viper.GetString("LOCAL_BEACON_CONN_STR")

	testCont.LocalBearerToken = viper.GetString("LOCAL_BEARER_TOKEN")
	testCont.ProductionLocalBearerToken = viper.GetString("PROD_LOCAL_BEARER_TOKEN")

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
