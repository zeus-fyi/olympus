package v0

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

// BuildPathsFromDirInPath assumes the package name is the parent dir name
func (l *Lib) BuildPathsFromDirInPath(root structs.Path, ext string) map[int]structs.Paths {
	m := make(map[int]structs.Paths)
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
				m[depth] = structs.Paths{}
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
