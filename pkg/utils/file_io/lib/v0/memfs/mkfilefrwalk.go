package memfs

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) WalkAndApplyFuncToFileType(p *structs.Path, ext string, f func(p string, fs *MemFS) error) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	fileSystem, err := m.Sub(p.DirIn)
	if err != nil {
		return err
	}
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(d.Name()) == ext {
			p.FnIn = path
			return f(p.FileDirOutFnInPath(), m)
		}
		return nil
	})
	return err
}
