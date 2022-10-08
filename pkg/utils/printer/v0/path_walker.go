package v0

import (
	"io/fs"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

// BuildPathsFromDirInPath assumes the package name is the parent dir name
func (l *Lib) BuildPathsFromDirInPath(root structs.Path, ext string) structs.Paths {
	dirOut := root.DirOut
	paths := structs.Paths{}
	_ = filepath.WalkDir(root.DirIn, func(walkDir string, d fs.DirEntry, err error) error {
		if l.Log.ErrHandler(err) != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			dirIn := filepath.Dir(walkDir)
			pkgName := filepath.Base(dirIn)
			paths.AddPathToSlice(l.NewPkgPathInOut(pkgName, dirIn, filepath.Join(dirOut, dirIn), d.Name()))
		}
		return l.Log.ErrHandler(err)
	})
	return paths
}
