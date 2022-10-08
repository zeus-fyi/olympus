package readers

import (
	"io/ioutil"
	"path"
)

func ReadFile(subDir, fn string) []byte {
	file := path.Join(subDir, fn)
	byteArray, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return byteArray
}
