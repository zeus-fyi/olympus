package v0

import (
	"io/fs"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

func (l *Lib) BuildPathsFromDirInPath(path structs.Path, ext string) structs.Paths {
	dirOut := path.DirOut
	paths := structs.Paths{}
	_ = filepath.WalkDir(path.DirIn, func(walkDir string, d fs.DirEntry, err error) error {
		if l.Log.ErrHandler(err) != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			dirIn := filepath.Dir(walkDir)
			paths.AddPathToSlice(l.NewPathInOut(dirIn, filepath.Join(dirOut, dirIn), d.Name()))
		}
		return l.Log.ErrHandler(err)
	})
	return paths
}
