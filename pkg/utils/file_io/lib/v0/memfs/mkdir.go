package memfs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
)

func (m *MemFS) MkPathDirAll(p *filepaths.Path) error {
	if err := m.MkdirAll(p.DirOut, 0700); err != nil {
		return err
	}
	if err := m.MkdirAll(p.DirIn, 0700); err != nil {
		return err
	}
	return nil
}
