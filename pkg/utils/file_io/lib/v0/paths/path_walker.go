package paths

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

// BuildPathsFromDirInPath assumes the package name is the parent dir name
func (l *PathLib) BuildPathsFromDirInPath(root filepaths.Path, ext string) map[int]filepaths.Paths {
	m := make(map[int]filepaths.Paths)
	depth := 0
	dirOut := root.DirOut
	depthStart := len(strings.Split(root.DirIn, "/")) - 1
	_ = filepath.WalkDir(root.DirIn, func(walkDir string, d fs.DirEntry, err error) error {
		if l.Log.ErrHandler(err) != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			dirIn := filepath.Dir(walkDir)
			pkgName := filepath.Base(dirIn)
			depth = len(strings.Split(walkDir, "/")) - len(strings.Split(dirOut, "/")) - depthStart
			if _, ok := m[depth]; !ok {
				m[depth] = filepaths.Paths{}
			}
			depthPaths := m[depth]

			nextPath := dirOut
			if depth != 0 {
				nextPath = path.Join(dirOut, pkgName)
			}
			depthPaths.AddPathToSlice(l.NewPkgPathInOut(pkgName, dirIn, nextPath, d.Name()))
			m[depth] = depthPaths
		}
		return l.Log.ErrHandler(err)
	})
	return m
}
