package printer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strings"
	"time"
)

var logLevelFilter = "error"

func InterfacePrinter(path, fn, env string, v interface{}) (interface{}, error) {
	jsonParams, e := json.MarshalIndent(&v, "", " ")
	if e != nil {
		return v, e
	}
	Printer(path, fmt.Sprintf("%s.json", fn), env, jsonParams)

	return v, nil
}

func Printer(subDir, filename, env string, data []byte) {
	ts := time.Now().Format(time.Stamp)

	if env == "development" {
		env = "dev"
	}

	fn := fmt.Sprintf("%s.%s.%s", ts, env, filename)
	folder := path.Join(env, subDir)
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
