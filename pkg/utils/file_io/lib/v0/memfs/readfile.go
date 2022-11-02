package memfs

import (
	"errors"
	"io/fs"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (m *MemFS) ReadFileFromPath(p *structs.Path) ([]byte, error) {
	var b []byte
	if p == nil {
		return b, errors.New("need to include a path")
	}
	b, err := fs.ReadFile(m, p.FileInPath())
	if err != nil {
		return b, err
	}
	return b, nil
}
