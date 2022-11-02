package memfs

import (
	"errors"
	"io/fs"
	"os"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/readers"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) MakeFile(p *structs.Path, content []byte) error {
	merr := m.MkPathDirAll(p)
	if merr != nil {
		return merr
	}
	if err := m.WriteFile(p.FileOutPath(), content, 0644); err != nil {
		return err
	}
	return nil
}

func (m *MemFS) MakeFilesFromWalk(p *structs.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r := readers.ReaderLib{}
	fileSystem := os.DirFS(p.DirIn)

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			bytesToTransfer := r.ReadFilePathPtr(p)
			newPath := structs.Path{
				DirOut: path,
			}
			terr := m.MakeFile(&newPath, bytesToTransfer)
			if terr != nil {
				return terr
			}
		}
		return nil
	})
	return err
}
