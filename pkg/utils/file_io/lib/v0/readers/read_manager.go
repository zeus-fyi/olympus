package readers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (l *ReaderLib) ReadFile(p structs.Path) []byte {
	byteArray, err := ioutil.ReadFile(p.FileInPath())
	if err != nil {
		panic(err)
	}
	return byteArray
}

func (l *ReaderLib) ReadJsonObject(p structs.Path, obj interface{}) interface{} {
	jsonByteArray := l.ReadFile(p)
	err := json.Unmarshal(jsonByteArray, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}
