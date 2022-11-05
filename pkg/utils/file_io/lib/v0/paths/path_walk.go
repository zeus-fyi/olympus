package paths

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func WalkAndApplyFuncToFileType(p structs.Path, ext string, f func(p string) error) error {
	fileSystem := os.DirFS(p.DirIn)
	root := path.Dir(p.DirIn)
	err := fs.WalkDir(fileSystem, root, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ext {
			filePath := pathJoin(p.DirIn, path)
			return f(filePath)
		}
		return nil
	})
	return err
}

func pathJoin(root, file string) string {
	return path.Join(root, file)
}
