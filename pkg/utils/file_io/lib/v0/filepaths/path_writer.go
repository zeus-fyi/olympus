package filepaths

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/chronos"
)

func (p *Path) WriteToFileOutPath(data []byte) error {
	// make path if it doesn't exist
	if _, err := os.Stat(p.FileOutPath()); os.IsNotExist(err) {
		_ = os.MkdirAll(p.DirOut, 0700) // Create your dir
	}
	err := os.WriteFile(p.FileOutPath(), data, 0644)
	return err
}

func (p *Path) Print(data []byte) error {
	ts := chronos.Chronos{}
	p.FnOut += fmt.Sprintf("_%d.log", ts.UnixTimeStampNow())
	err := p.WriteToFileOutPath(data)
	if err != nil {
		return err
	}
	file, err := p.OpenFileOutPath()
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	linesToWrite := strings.Split(string(data), "time")
	for _, line := range linesToWrite {
		_, berr := writer.WriteString("time" + string(line))
		if berr != nil {
			return berr
		}

		_ = writer.Flush()
		return nil
	}
	return nil
}
