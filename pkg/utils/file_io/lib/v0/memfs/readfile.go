package memfs

import (
	"io/fs"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) ReadFile(p structs.Path) error {

	_, err := fs.ReadFile(m, p.FileInPath())
	if err != nil {
		return err
	}
	return nil
}
