package cookbook

import (
	"os"
	"path"
	"runtime"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"golang.org/x/exp/maps"
)

func UseCookbookDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "")
	err := os.Chdir(dir)
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func (a *Cookbook) GetTopologicallySortedPaths(path filepaths.Path) filepaths.Paths {
	templatePathMaps := fileIO.BuildPathsFromDirInPath(path, ".go")
	depth := len(maps.Keys(templatePathMaps))
	tmp := filepaths.Paths{}
	for i := 0; i <= depth; i++ {
		recipePaths := templatePathMaps[i]
		for _, r := range recipePaths.Slice {
			tmp.AddPathToSlice(r)
		}
		i += i
	}
	return tmp
}
