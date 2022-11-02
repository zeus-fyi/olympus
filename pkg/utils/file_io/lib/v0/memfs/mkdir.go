package memfs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) MkDir(p structs.Path) error {

	if err := m.MkdirAll(p.DirOut, 0o700); err != nil {
		return err
	}

	return nil
}
