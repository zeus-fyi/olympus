package datastores

import (
	"os"
	"path"
	"runtime"
)

type Datastores struct {
}

func (d *Datastores) ForceDirToCallerDatastoreDirRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}
