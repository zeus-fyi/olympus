package writers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

var logLevelFilter = "error"

func (l *WriterLib) InterfacePrinter(path filepaths.Path, v interface{}) (interface{}, error) {
	jsonParams, e := json.MarshalIndent(&v, "", " ")
	if e != nil {
		return v, e
	}
	err := l.Print(path, jsonParams)
	return v, err
}

func (l *WriterLib) Print(p filepaths.Path, data []byte) error {
	ts := l.Log.UnixTimeStampNow()

	if p.Env == "development" {
		p.Env = "dev"
	}
	fn := fmt.Sprintf("%d.%s.%s", ts, p.Env, p.FnIn)
	p.FnIn = fn
	err := l.CreateFile(p, nil)
	if l.Log.ErrHandler(err) != nil {
		return err
	}

	file, err := l.OpenFile(p)
	if l.Log.ErrHandler(err) != nil {
		return err
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	linesToWrite := strings.Split(string(data), "time")
	for _, line := range linesToWrite {
		if !strings.Contains(line, "vendor") && l.loglevel(line, logLevelFilter) {
			_, berr := writer.WriteString("time" + string(line))
			if l.Log.ErrHandler(berr) != nil {
				return berr
			}
		}
		_ = writer.Flush()
		return nil
	}
	return nil
}

func (l *WriterLib) loglevel(line, level string) bool {
	if level == "error" {
		return strings.Contains(line, level) || strings.Contains(line, "warn")
	} else {
		return strings.Contains(line, level)
	}
}
