package compression

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"path"

	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/memfs"
	"github.com/zeus-fyi/olympus/pkg/utils/file_io/lib/v0/structs"
)

func (c *Compression) UnGzipFromInMemFsOutToInMemFS(p *structs.Path, fs memfs.MemFS) error {
	if p == nil {
		return errors.New("need to include a path")
	}
	r, err := fs.ReadFileFromPath(p)
	if err != nil {
		return err
	}
	in := &bytes.Buffer{}
	_, err = in.Write(r)
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, herr := tr.Next()

		switch {

		// if no more files are found return
		case herr == io.EOF:
			return nil

		// return any other error
		case herr != nil:
			return herr

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// check the file type
		switch header.Typeflag {
		// if it's a dir do nothing
		case tar.TypeDir:
		// if it's a file create it
		case tar.TypeReg:
			p.Fn = header.Name

			fo := p.FileOutPath()
			p.DirIn = path.Dir(fo)

			out := &bytes.Buffer{}
			if _, cerr := io.Copy(out, tr); cerr != nil {
				return cerr
			}
			ferr := fs.MakeFile(p, out.Bytes())
			if ferr != nil {
				return ferr
			}
		}
	}
}
