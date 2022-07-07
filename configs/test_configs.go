package configs

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

var testCont TestContainer

type TestContainer struct {
	PG_CONN            string
	BEACON_NODE_INFURA string
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

	testCont.PG_CONN = viper.GetString("PG_CONN")
	testCont.BEACON_NODE_INFURA = viper.GetString("BEACON_NODE_INFURA")
	return testCont
}
