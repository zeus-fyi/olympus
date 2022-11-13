package mocks

import (
	"os"
	"path"
	"runtime"
)

func ChangeToMockDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
