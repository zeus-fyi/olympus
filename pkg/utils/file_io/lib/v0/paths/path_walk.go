package paths

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
	"github.com/zeus-fyi/olympus/pkg/utils/string_utils"
)

func WalkAndApplyFuncToFileType(p structs.Path, ext string, f func(p string) error) error {
	fileSystem := os.DirFS(p.DirIn)
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ext {
			if string_utils.FilterStringWithOpts(path, &p.FilterFiles) {
				filePath := pathJoin(p.DirIn, path)
				return f(filePath)
			}
		}
		return nil
	})
	return err
}

func pathJoin(root, file string) string {
	return path.Join(root, file)
}
