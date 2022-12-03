package compression

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/pierrec/lz4"
	"github.com/rs/zerolog/log"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/filepaths"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
)

func (c *Compression) Lz4Decompress(p *filepaths.Path) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := os.Open(p.FileInPath())
	if err != nil {
		log.Err(err).Msg("Compression: Lz4Decompress, os.Open(p.FileInPath())")
		return err
	}
	defer r.Close()
	lz4Reader := lz4.NewReader(r)
	if err != nil {
		log.Err(err).Msg("Compression: Lz4Decompress, lz4.NewReader(r)")
		return err
	}
	return tarReader(p, lz4Reader)
}

func (c *Compression) Lz4DecompressInMemFsFile(p *filepaths.Path, inMemFs memfs.MemFS) (memfs.MemFS, error) {
	if p == nil {
		return inMemFs, errors.New("need to include a path")
	}

	b, err := inMemFs.ReadFileInPath(p)
	if err != nil {
		log.Err(err).Msg("Compression: Lz4CompressInMemFsFile")
		return inMemFs, err
	}
	in := bytes.Buffer{}
	r := bufio.NewReader(&in)

	buf := bytes.Buffer{}
	// make a write buffer
	w := bufio.NewWriter(&buf)

	lz4Reader := lz4.NewReader(r)

	// make a buffer to keep chunks that are read
	buffer := make([]byte, 1024)
	for {
		// read a chunk
		n, rerr := lz4Reader.Read(buffer)
		if rerr != nil && rerr != io.EOF {
			panic(rerr)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, cerr := w.Write(buffer[:n]); cerr != nil {
			panic(cerr)
		}
	}
	err = inMemFs.MakeFileOut(p, b)
	p.DirIn = p.DirOut
	p.FnIn = p.FnOut
	if err != nil {
		log.Err(err).Msg("Compression: Lz4DecompressInMemFsFile")
		return inMemFs, err
	}
	return inMemFs, err
}
