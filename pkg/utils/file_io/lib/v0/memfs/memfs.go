package memfs

import (
	"io/fs"

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

func tmp() {
	memfs := memoryfs.New()

	if err := memfs.MkdirAll("my/dir", 0o700); err != nil {
		panic(err)
	}

	if err := memfs.WriteFile("my/dir/file.txt", []byte("hello world"), 0o600); err != nil {
		panic(err)
	}

	_, err := fs.ReadFile(memfs, "my/dir/file.txt")
	if err != nil {
		panic(err)
	}
}
