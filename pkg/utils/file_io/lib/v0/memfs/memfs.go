package memfs

import (
	"github.com/zeus-fyi/memoryfs"
)

type MemFS struct {
	*memoryfs.FS
}

func NewMemFs() MemFS {
	memfs := memoryfs.New()
	m := MemFS{memfs}
	return m
}
