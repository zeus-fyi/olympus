package printer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/zeus-fyi/olympus/pkg/utils/env"
)

var logLevelFilter = "error"

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

	CreateFile(p, fn, folder, nil)

	file, err := OpenFile(p)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	linesToWrite := strings.Split(string(data), "time")
	for _, line := range linesToWrite {
		if !strings.Contains(line, "vendor") && loglevel(line, logLevelFilter) {
			_, berr := writer.WriteString("time" + string(line))
			if berr != nil {
				log.Fatalf("Got error while writing to a file. Err: %s", berr.Error())
			}
		}
		_ = writer.Flush()
		return
	}
}

func loglevel(line, level string) bool {
	if level == "error" {
		return strings.Contains(line, level) || strings.Contains(line, "warn")
	} else {
		return strings.Contains(line, level)
	}
}
