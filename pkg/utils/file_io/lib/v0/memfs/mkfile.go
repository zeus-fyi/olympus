package memfs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) MakeFile(p *structs.Path, content []byte) error {
	merr := m.MkPathDirAll(p)
	if merr != nil {
		return merr
	}
	if err := m.WriteFile(p.FileDirOutFnInPath(), content, 0644); err != nil {
		return err
	}
	return nil
}
