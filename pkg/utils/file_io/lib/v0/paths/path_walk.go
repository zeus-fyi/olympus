package paths

import (
	"io/fs"
	"os"
	"path/filepath"
)

func WalkAndApplyFuncToFileType(dir, ext string, f func(p string) error) error {
	fileSystem := os.DirFS(dir)

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
