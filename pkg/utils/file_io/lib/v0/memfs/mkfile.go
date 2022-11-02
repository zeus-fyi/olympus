package memfs

import (
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) MakeFile(p structs.Path, content []byte) error {

	if err := m.WriteFile(p.FileOutPath(), content, 0o600); err != nil {
		return err
	}

	return nil
}
