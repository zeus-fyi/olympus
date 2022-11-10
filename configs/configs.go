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

	LocalBearerToken           string
	ProductionLocalBearerToken string

	ProductionLocalTemporalBearerToken string

	DevTemporalHostPort string
	DevTemporalNs       string

	DevTemporalAuth temporal_client.TemporalAuth
	DevAuthKeysCfg  auth_keys_config.AuthKeysCfg

	ProdLocalTemporalAuth temporal_client.TemporalAuth
	ProdLocalAuthKeysCfg  auth_keys_config.AuthKeysCfg
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

func InitLocalTestConfigs() TestContainer {
	InitEnvFromConfig(forceDirToCallerLocation())
	testCont.Env = viper.GetString("ENV")

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

	testCont.ProductionLocalTemporalBearerToken = viper.GetString("PROD_LOCAL_TEMPORAL_BEARER_TOKEN")

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
