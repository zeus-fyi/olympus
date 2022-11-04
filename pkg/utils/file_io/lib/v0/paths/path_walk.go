package paths

import (
	"io/fs"
	"path/filepath"
)

func WalkAndApplyFuncToFileType(fileSystem fs.FS, dir, ext string, f func(p string) error) error {
	err := fs.WalkDir(fileSystem, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ext {
			return f(path)
		}
		return nil
	})
	return err
}
