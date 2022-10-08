package v0

import (
	"encoding/json"
	"io/ioutil"

	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

func (l *Lib) ReadFile(p structs.Path) []byte {
	byteArray, err := ioutil.ReadFile(p.FileInPath())
	if err != nil {
		panic(err)
	}
	return byteArray
}

func (l *Lib) ReadJsonObject(p structs.Path, obj interface{}) interface{} {
	jsonByteArray := l.ReadFile(p)
	err := json.Unmarshal(jsonByteArray, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}
