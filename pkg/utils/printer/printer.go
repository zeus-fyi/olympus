package printer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/env"
)

type PrintPath string

func (p PrintPath) Local() string {
	return "artifacts/local/"
}

func (p PrintPath) Dev() string {
	return "artifacts/dev/"
}

func (p PrintPath) Staging() string {
	return "artifacts/staging/"
}

func (p PrintPath) Production() string {
	return "artifacts/production/"
}

func InterfacePrinter(path, fn string, v interface{}) (interface{}, error) {
	jsonParams, e := json.MarshalIndent(&v, "", " ")
	if e != nil {
		return v, e
	}
	Printer(path, fmt.Sprintf("%s.json", fn), jsonParams)

	return v, nil
}

func Printer(subDir, filename string, data []byte) {
	var pp PrintPath
	ts := time.Now().Format(time.Stamp)

	envPrinter := env.Str
	if env.Str == "development" {
		envPrinter = "dev"
	}
	fn := fmt.Sprintf("%s.%s.%s", ts, envPrinter, filename)
	folder := path.Join(env.SetEnvParam(pp), subDir)
	p := path.Join(folder, fn)

	// make path if it doesn't exist
	if _, err := os.Stat(p); os.IsNotExist(err) {
		_ = os.MkdirAll(folder, 0700) // Create your dir
	}

	err := ioutil.WriteFile(p, data, 0644)
	if err != nil {
		log.Fatalf("error writing %s: %s", fn, err)
	}
}
