package readers

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

func ReadJsonFile(subDir, fn string) []byte {
	jsonFile := path.Join(subDir, fn)
	jsonByteArray, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	return jsonByteArray
}

func ReadJsonObject(subDir, fn string, obj interface{}) interface{} {
	jsonByteArray := ReadJsonFile(subDir, fn)
	err := json.Unmarshal(jsonByteArray, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}
